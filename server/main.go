package main

import (
	"context"
	"log"
	"video-platform/internal/controller"
	"video-platform/internal/middleware"
	"video-platform/internal/model"
	"video-platform/internal/service"

	_ "github.com/gogf/gf/contrib/drivers/mysql/v2"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

// 初始化数据库
func initDatabase() error {
	// 创建数据库
	sql := "CREATE DATABASE IF NOT EXISTS video_platform DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;"
	_, err := g.DB().Exec(context.Background(), sql)
	if err != nil {
		return err
	}
	
	// 切换到新创建的数据库
	_, err = g.DB().Exec(context.Background(), "USE video_platform;")
	return err
}

func main() {
	ctx := context.Background()
	
	// 初始化数据库
	if err := initDatabase(); err != nil {
		log.Fatalf("创建数据库失败: %v", err)
	}
	
	// 初始化数据库表
	if err := model.Device.InitTable(ctx); err != nil {
		log.Fatalf("初始化数据库表失败: %v", err)
	}

	// 初始化MQTT服务
	service.GetMQTTService()

	s := g.Server()

	s.Group("/api", func(group *ghttp.RouterGroup) {
		group.Middleware(middleware.CORS)
		group.Group("/v1", func(group *ghttp.RouterGroup) {
			// 设备管理路由
			group.GET("/devices", controller.DeviceController.List)
			group.POST("/devices", controller.DeviceController.Add)
			group.GET("/devices/:deviceId", controller.DeviceController.Get)
			group.PUT("/devices/:deviceId", controller.DeviceController.Update)
			group.DELETE("/devices/:deviceId", controller.DeviceController.Delete)
			group.GET("/devices/:deviceId/status", controller.DeviceController.GetStatus)
			
			// 设备图像路由
			group.GET("/devices/:deviceId/realtime", controller.DeviceController.GetRealtimeImage)
			group.GET("/devices/:deviceId/images", controller.DeviceController.GetHistoryImages)

			// 测试MQTT消息发布
			group.GET("/test/mqtt/:deviceId", func(r *ghttp.Request) {
				deviceId := r.Get("deviceId").String()
				testData := []byte("test image data")
				if err := service.GetMQTTService().TestPublish(deviceId, testData); err != nil {
					r.Response.WriteStatus(500)
					r.Response.WriteJson(g.Map{"error": err.Error()})
					return
				}
				r.Response.WriteJson(g.Map{"message": "测试消息已发送"})
			})
		})
	})

	s.Run()
} 