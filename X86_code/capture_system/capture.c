#include <stdio.h>
#include <stdlib.h>
#include <sys/types.h>
#include <sys/stat.h>
#include <fcntl.h>
#include <unistd.h>
#include <sys/ioctl.h>
#include <string.h>
#include <sys/mman.h>
#include <errno.h>
#include <linux/videodev2.h>
#include <linux/fb.h>
int main(void)
{
    // 定义一个设备描述符
    int fd;
    struct v4l2_capability cap; // 定义结构体类型
    fd = open("/dev/video0", O_RDWR);
    if (fd < 0)
    {
        perror("video设备打开失败\n");
        return -1;
    }
    else
    {
        printf("video设备打开成功\n");
    }
    ioctl(fd, VIDIOC_QUERYCAP, &cap);
    if (!(V4L2_CAP_VIDEO_CAPTURE & cap.capabilities))
    {
        perror("Error: No capture video device!\n");
        return -1;
    }
    printf("驱动名 : %s\n", cap.driver);
    printf("设备名字 : %s\n", cap.card);
    printf("总线信息 : %s\n", cap.bus_info);
    printf("驱动版本号 : %d\n", cap.version);
    struct v4l2_frmsizeenum frmsize;
    frmsize.index = 0;
    frmsize.type = V4L2_BUF_TYPE_VIDEO_CAPTURE;

    printf("MJEPG格式支持所有分辨率如下:\n");
    frmsize.pixel_format = V4L2_PIX_FMT_MJPEG;
    while (ioctl(fd, VIDIOC_ENUM_FRAMESIZES, &frmsize) == 0)
    {
        printf("frame_size<%d*%d>\n", frmsize.discrete.width, frmsize.discrete.height);
        frmsize.index++;
    }
    struct v4l2_format vfmt; // 头文件自带v4l2_format结构体，引用即可
    vfmt.type = V4L2_BUF_TYPE_VIDEO_CAPTURE;
    vfmt.fmt.pix.width = 1920;
    vfmt.fmt.pix.height = 1080;
    vfmt.fmt.pix.pixelformat = V4L2_PIX_FMT_MJPEG;
    if (ioctl(fd, VIDIOC_S_FMT, &vfmt) < 0)
    {
        perror("设置格式失败\n");
        return -1;
    }
    // 检查设置参数是否生效
    if (ioctl(fd, VIDIOC_G_FMT, &vfmt) < 0)
    {
        perror("获取设置格式失败\n");
        return -1;
    }
    else if (vfmt.fmt.pix.width == 1920 && vfmt.fmt.pix.height == 1080 && vfmt.fmt.pix.pixelformat == V4L2_PIX_FMT_MJPEG)
    {
        printf("设置格式生效,实际分辨率大小<%d * %d>,图像格式:Motion-JPEG\n", vfmt.fmt.pix.width, vfmt.fmt.pix.height);
    }
    else
    {
        printf("设置格式未生效\n");
    }
    /* 获取 streamparm */
    struct v4l2_streamparm streamparm = {0};
    streamparm.type = V4L2_BUF_TYPE_VIDEO_CAPTURE;
    ioctl(fd, VIDIOC_G_PARM, &streamparm);
    if (V4L2_CAP_TIMEPERFRAME & streamparm.parm.capture.capability)
    {
        streamparm.parm.capture.timeperframe.numerator = 1;
        streamparm.parm.capture.timeperframe.denominator = 60; // 60fps
        if (0 > ioctl(fd, VIDIOC_S_PARM, &streamparm))
        {
            fprintf(stderr, "ioctl error: VIDIOC_S_PARM: %s\n", strerror(errno));
            return -1;
        }
    }
    else
        fprintf(stderr, "不支持帧率设置");
    struct v4l2_requestbuffers reqbuf;
    reqbuf.type = V4L2_BUF_TYPE_VIDEO_CAPTURE;
    reqbuf.count = 3; // 3个帧缓冲
    reqbuf.memory = V4L2_MEMORY_MMAP;
    if (ioctl(fd, VIDIOC_REQBUFS, &reqbuf) < 0)
    {
        perror("申请缓冲区失败\n");
        return -1;
    }
    void *frm_base[3]; // 映射后到空间的首地址
    unsigned int frm_size[3];

    struct v4l2_buffer buf;
    buf.type = V4L2_BUF_TYPE_VIDEO_CAPTURE;
    buf.memory = V4L2_MEMORY_MMAP;
    for (buf.index = 0; buf.index < 3; buf.index++)
    {
        ioctl(fd, VIDIOC_QUERYBUF, &buf);
        frm_base[buf.index] = mmap(NULL, buf.length, PROT_READ | PROT_WRITE, MAP_SHARED, fd, buf.m.offset); //
        frm_size[buf.index] = buf.length;

        if (frm_base[buf.index] == MAP_FAILED)
        {
            perror("mmap failed\n");
            return -1;
        }
        // 入队操作
        if (ioctl(fd, VIDIOC_QBUF, &buf) < 0)
        {
            perror("入队失败\n");
            return -1;
        }
    }
    enum v4l2_buf_type type = V4L2_BUF_TYPE_VIDEO_CAPTURE;
    if (ioctl(fd, VIDIOC_STREAMON, &type) < 0)
    {
        perror("开始采集失败\n");
        return -1;
    }
    while (1)
    {
        struct v4l2_buffer readbuffer;
        readbuffer.type = V4L2_BUF_TYPE_VIDEO_CAPTURE;
        readbuffer.memory = V4L2_MEMORY_MMAP;
        if (ioctl(fd, VIDIOC_DQBUF, &readbuffer) < 0)
        {
            perror("读取帧失败\n");
        }
        // 保存这一帧，格式为jpg
        FILE *file = fopen("qqq.jpg", "w+");
        fwrite(frm_base[readbuffer.index], buf.length, 1, file);
        fclose(file);
        if (ioctl(fd, VIDIOC_QBUF, &readbuffer) < 0)
        {
            perror("入队失败\n");
        }
        sleep(1);
    }
    // 停止采集
    if (ioctl(fd, VIDIOC_STREAMOFF, 1) < 0)
    {
        perror("停止采集失败\n");
        return -1;
    }

    // 释放映射
    for (int i = 0; i < 3; i++)
    {
        munmap(frm_base[i], frm_size[i]);
    }

    close(fd);
}
