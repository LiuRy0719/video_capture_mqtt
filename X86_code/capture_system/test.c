#include <stdio.h>
#include <sys/types.h>
#include <sys/stat.h>
#include <fcntl.h>
#include <stdlib.h>
#include <unistd.h>
#include <sys/ioctl.h>//iocontrol设备控制操作
#include <linux/videodev2.h>//用于视频设备（如摄像头）的接口
#include <string.h>
#include <sys/mman.h>
#include <jpeglib.h>//JPEG 图像处理库
#include <linux/fb.h>//帧缓冲区相关
#include <stdio.h>
#include <mosquitto.h>
#include <signal.h>
 
#define MQTT_BROKER "112.6.224.25"//定义要连接的代理的ip地址，localhost指连接本地的ip
#define MQTT_PORT 20042//指代理监听的端口号，这也是连接的入口
#define MQTT_TOPIC "screen/stream"//订阅的主题
#define MQTT_QOS 1//消息服务等级
 
volatile sig_atomic_t stop = 0;
//定义了一个sig_atomic类型的全局变量stop，volatile 表示该变量可能会在程序外部被修改（如通过信号处理函数）
//sig_atomic类型的作用？通常是int类型，确保这些操作是原子性的（即不会被信号中断）
//信号处理函数是在内存的哪里执行的？在栈中吗？信号属于异步事件处理机制，会中断正常的程序执行，并在当前线程的上下文中运行，上下文（context）指的是程序执行时的状态
void handle_sigint(int sig) {
    stop = 1; 
}//定义了一个信号处理函数，当接收到 SIGINT（通常是通过按下 Ctrl+C 触发）时，将全局变量 stop 设置为 1，以通知程序退出
 
// 连接回调函数，客户端收到代理的消息后执行，这个执行是在主线程之外的
//回调函数是在栈中执行吗？是否要多开一个进程？这个网络信号的回调处理函数是在后台监听线程中执行的，是同步机制。
void on_connect(struct mosquitto *mosq, void *obj, int reason_code) {
    if (reason_code != 0) {
        fprintf(stderr, "Failed to connect to MQTT broker with reason code: %d\n", reason_code);
        exit(EXIT_FAILURE);
    }
    printf("Connected to MQTT broker.\n");
    mosquitto_subscribe(mosq, NULL, "screen/stream", 1);//订阅主题
}
 
int main()
{
    signal(SIGINT, handle_sigint);//把SIGINT信号和信号处理函数连接起来
    int fd = open("/dev/video0",O_RDWR); //打开摄像头设备
    if (fd < 0)return -1;
	
    struct v4l2_format vfmt;//这个结构体用于设置摄像头的格式，然后作为实参输入v4l2的方法中
 
    vfmt.type = V4L2_BUF_TYPE_VIDEO_CAPTURE; //表示视频捕获
    vfmt.fmt.pix.width = 640; //设置摄像头采集参数，分辨率
    vfmt.fmt.pix.height = 480;
    vfmt.fmt.pix.pixelformat = V4L2_PIX_FMT_MJPEG; //设置像素格式为 MJPEG（即摄像头输出 JPEG 格式的视频流）
 
    int ret = ioctl(fd,VIDIOC_S_FMT,&vfmt);//使用 ioctl 系统调用将这些设置应用到摄像头设备
//ioctl是如何传参的？
    if (ret < 0)perror("设置格式失败1");
 
    //申请内核空间
    struct v4l2_requestbuffers reqbuffer;
    reqbuffer.type = V4L2_BUF_TYPE_VIDEO_CAPTURE;//视频捕获
    reqbuffer.count = 4; //申请四个缓冲区
    reqbuffer.memory = V4L2_MEMORY_MMAP;  //使用内存映射（mmap）方式访问缓冲区
 
    ret = ioctl(fd,VIDIOC_REQBUFS,&reqbuffer);//请求内核分配缓冲区
    if (ret < 0)perror("申请空间失败");
    
 
    //映射
    unsigned char *mptr[4];//保存映射后用户空间的首地址
    unsigned int size[4];//缓冲区大小
    struct v4l2_buffer mapbuffer;
    //初始化type和index
    mapbuffer.type = V4L2_BUF_TYPE_VIDEO_CAPTURE;
 
    for(int i = 0; i <4;i++) {
        mapbuffer.index = i;
        ret = ioctl(fd,VIDIOC_QUERYBUF,&mapbuffer); //从内核空间中查询一个空间作映射
        if (ret < 0)perror("查询内核空间失败");
        //映射到用户空间
        mptr[i] = (unsigned char *)mmap(NULL,mapbuffer.length,PROT_READ|PROT_WRITE,MAP_SHARED,fd,mapbuffer.m.offset);
        size[i] = mapbuffer.length; //保存映射长度用于后期释放
        //查询后通知内核已经放回
        ret = ioctl(fd,VIDIOC_QBUF,&mapbuffer);
        if (ret < 0)perror("放回失败");    
    }
 
    //开始采集
    int type = V4L2_BUF_TYPE_VIDEO_CAPTURE;
    ret = ioctl(fd,VIDIOC_STREAMON,&type);
    if (ret < 0)perror("开启失败");
 
//----------------------------------------------------
    struct mosquitto *mosq = NULL;
 
    // 初始化Mosquitto客户端
    mosquitto_lib_init();
    mosq = mosquitto_new(NULL, true, NULL);
    if (!mosq) {
        fprintf(stderr, "Failed to create Mosquitto client.\n");
        return EXIT_FAILURE;
    }
 
    // 设置回调函数
    mosquitto_connect_callback_set(mosq, on_connect);
 
    // 连接到MQTT代理
    ret = mosquitto_connect(mosq, MQTT_BROKER, MQTT_PORT, 60);
    if (ret != MOSQ_ERR_SUCCESS) {
        fprintf(stderr, "Failed to connect to MQTT broker: %d\n", ret);
        mosquitto_destroy(mosq);
        return EXIT_FAILURE;
    }
 
    // 开始网络循环，处理消息，网络循环这个机制是怎么实现的？
    //启动一个后台线程（thread）来处理 MQTT 客户端的网络通信和消息分发。这个线程会独立于主线程运行，负责处理与 MQTT 代理（broker）之间的网络交互、消息接收和发送等任务。
    // int a = mosquitto_loop_start(mosq);
    // if(a!=MOSQ_ERR_SUCCESS)
    // {
    //     fprintf(stderr,"Failed to start to loop_start\n");
    //     exit(EXIT_FAILURE);
    // }
    
    while(!stop)
    {
	 //从队列中提取一帧数据
        struct v4l2_buffer readbuffer;
        readbuffer.type = V4L2_BUF_TYPE_VIDEO_CAPTURE; //每个结构体都需要设置type
        ret = ioctl(fd,VIDIOC_DQBUF,&readbuffer);//用于获取内核缓冲区的信息
        if (ret < 0)perror("读取数据失败");
 
        // 将数据封装为MQTT消息
        unsigned char *jpeg_data = mptr[readbuffer.index];//通过 mptr[readbuffer.index] 访问缓冲区的用户空间地址，并处理该缓冲区中的 MJPEG 数据
	//readbuffer.index 表示当前帧数据所在的缓冲区索引。mptr[readbuffer.index] 是之前通过 mmap 映射的用户空间地址，指向当前帧数据的起始位置。
        printf("readbuffer.index=%d",readbuffer.index);
	int jpeg_size = readbuffer.bytesused;
//readbuffer.bytesused 表示当前帧数据的实际大小（以字节为单位）。
	if (readbuffer.bytesused == 0) {
        fprintf(stderr, "No data in buffer.\n");
        continue; // 跳过当前循环
        }
	
        // 发布消息到MQTT主题
        ret = mosquitto_publish(mosq, NULL, MQTT_TOPIC, jpeg_size, jpeg_data, MQTT_QOS, false);
        if (ret != MOSQ_ERR_SUCCESS)perror("MQTT发布失败");
        
	//通知内核使用完毕，内核标记这个缓冲区为可用状态，重新入队，等待摄像头写入MJEPG数据块
        ret = ioctl(fd, VIDIOC_QBUF, &readbuffer);
        if(ret < 0)perror("放回队列失败");
     }
       
 
//清理资源--------------------------------------------------------
    printf("successful out\n");
    mosquitto_destroy(mosq);
    mosquitto_lib_cleanup();
 
    ret = ioctl(fd,VIDIOC_STREAMOFF,&type);
 
    for(int i=0; i<4; i++)munmap(mptr[i], size[i]);
 
    close(fd); //关闭文件
    return 0;
    
 
}
