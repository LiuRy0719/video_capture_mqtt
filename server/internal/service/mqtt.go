package service

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
	"video-platform/internal/model"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type MQTTService struct {
	client     mqtt.Client
	imageDir   string        // 图像存储目录
	deviceData sync.Map      // 存储设备数据
}

var (
	mqttService *MQTTService
	once        sync.Once
)

// 获取MQTT服务实例
func GetMQTTService() *MQTTService {
	once.Do(func() {
		mqttService = &MQTTService{}
		mqttService.init()
	})
	return mqttService
}

// 初始化MQTT客户端
func (s *MQTTService) init() {
	// 创建图像存储目录
	s.imageDir = "images"
	if err := os.MkdirAll(s.imageDir, 0755); err != nil {
		log.Printf("创建图像存储目录失败: %v", err)
	}
	
	opts := mqtt.NewClientOptions()
	opts.AddBroker("tcp://112.6.224.25:20042")
	opts.SetClientID(fmt.Sprintf("server_%d", time.Now().Unix()))
	opts.SetDefaultPublishHandler(s.messageHandler)
	opts.SetOnConnectHandler(s.onConnect)
	opts.SetConnectionLostHandler(s.onConnectionLost)
	
	// 添加更多连接选项
	opts.SetAutoReconnect(true)
	opts.SetMaxReconnectInterval(time.Second * 10)
	opts.SetKeepAlive(time.Second * 60)
	opts.SetCleanSession(true)
	opts.SetOrderMatters(false) // 不要求消息顺序
	
	// 设置详细的调试日志
	mqtt.DEBUG = log.New(log.Writer(), "MQTT DEBUG: ", log.Ltime|log.Lshortfile)
	mqtt.ERROR = log.New(log.Writer(), "MQTT ERROR: ", log.Ltime|log.Lshortfile)

	s.client = mqtt.NewClient(opts)
	if token := s.client.Connect(); token.Wait() && token.Error() != nil {
		log.Printf("MQTT连接失败: %v", token.Error())
		return
	}
}

// 连接成功回调
func (s *MQTTService) onConnect(client mqtt.Client) {
	log.Println("MQTT已连接")
	// 订阅所有设备的图像主题
	if token := client.Subscribe("device/+/image", 1, s.messageHandler); token.Wait() && token.Error() != nil {
		log.Printf("订阅主题失败: %v", token.Error())
	} else {
		log.Printf("成功订阅主题: device/+/image")
	}
}

// 连接断开回调
func (s *MQTTService) onConnectionLost(client mqtt.Client, err error) {
	log.Printf("MQTT连接断开: %v", err)
}

// 消息处理
func (s *MQTTService) messageHandler(client mqtt.Client, msg mqtt.Message) {
	// 详细的消息日志
	log.Printf("收到MQTT消息:")
	log.Printf("- 主题: '%s'", msg.Topic())
	log.Printf("- 主题字节: %v", []byte(msg.Topic()))
	log.Printf("- QoS: %d", msg.Qos())
	log.Printf("- 重复消息: %v", msg.Duplicate())
	log.Printf("- 消息ID: %d", msg.MessageID())
	log.Printf("- 数据长度: %d字节", len(msg.Payload()))
	
	// 从主题中提取设备ID
	topic := strings.TrimSpace(msg.Topic())
	parts := strings.Split(topic, "/")
	log.Printf("- 主题分割: %v", parts)
	
	if len(parts) != 3 || parts[0] != "device" || parts[2] != "image" {
		log.Printf("无效的主题格式: '%s'", msg.Topic())
		return
	}
	deviceId := parts[1]
	log.Printf("- 设备ID: %s", deviceId)
	
	// 检查设备是否存在，如果不存在则添加
	device, err := model.Device.Get(context.Background(), deviceId)
	if err != nil || device == nil {
		// 设备不存在，创建新设备
		newDevice := &model.DeviceModel{
			Id:         deviceId,
			Name:       fmt.Sprintf("Device-%s", deviceId),
			Status:     "online",
			LastActive: time.Now(),
		}
		if err := model.Device.Add(context.Background(), newDevice); err != nil {
			log.Printf("添加新设备失败: %v", err)
		} else {
			log.Printf("成功添加新设备: %s", deviceId)
		}
	} else {
		// 更新设备状态为在线
		if err := model.Device.UpdateStatus(context.Background(), deviceId, "online"); err != nil {
			log.Printf("更新设备状态失败: %v", err)
		} else {
			log.Printf("设备 %s 状态已更新为在线", deviceId)
		}
	}
	
	// 创建设备专属的图像存储目录
	deviceDir := filepath.Join(s.imageDir, deviceId)
	if err := os.MkdirAll(deviceDir, 0755); err != nil {
		log.Printf("创建设备图像目录失败: %v", err)
		return
	}
	
	// 生成图像文件名（使用时间戳）
	timestamp := time.Now().Format("20060102_150405")
	filename := filepath.Join(deviceDir, fmt.Sprintf("%s.jpg", timestamp))
	
	// 保存图像到文件
	if err := os.WriteFile(filename, msg.Payload(), 0644); err != nil {
		log.Printf("保存图像文件失败: %v", err)
		return
	}
	
	// 更新设备最新图像的文件路径
	s.deviceData.Store(deviceId, filename)
	log.Printf("设备 %s 的图像已保存到文件: %s", deviceId, filename)
}

// 获取设备最新图像
func (s *MQTTService) GetDeviceImage(deviceId string) []byte {
	// 从内存中获取最新图像的文件路径
	if filename, ok := s.deviceData.Load(deviceId); ok {
		// 读取图像文件
		if data, err := os.ReadFile(filename.(string)); err == nil {
			return data
		} else {
			log.Printf("读取图像文件失败: %v", err)
		}
	}
	return nil
}

// 测试发布消息
func (s *MQTTService) TestPublish(deviceId string, data []byte) error {
	topic := fmt.Sprintf("device/%s/image", deviceId)
	token := s.client.Publish(topic, 1, false, data)
	if token.Wait() && token.Error() != nil {
		return fmt.Errorf("发布消息失败: %v", token.Error())
	}
	log.Printf("测试消息已发布到主题: %s", topic)
	return nil
} 