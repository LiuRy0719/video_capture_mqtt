package model

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"time"

	"github.com/gogf/gf/v2/frame/g"
)

// DeviceModel 设备表结构
type DeviceModel struct {
	Id         string    `json:"id" dc:"设备ID"`
	Name       string    `json:"name" dc:"设备名称"`
	Status     string    `json:"status" dc:"设备状态 online/offline"`
	LastActive time.Time `json:"lastActive" dc:"最后活跃时间"`
	CreatedAt  time.Time `json:"createdAt" dc:"创建时间"`
	UpdatedAt  time.Time `json:"updatedAt" dc:"更新时间"`
}

// 请求结构体
type DeviceListReq struct {
	g.Meta `path:"/devices" method:"get" tags:"设备管理" summary:"获取设备列表"`
}

type DeviceGetReq struct {
	g.Meta   `path:"/devices/{deviceId}" method:"get" tags:"设备管理" summary:"获取设备详情"`
	DeviceId string `json:"deviceId" v:"required" dc:"设备ID"`
}

type DeviceAddReq struct {
	g.Meta `path:"/devices" method:"post" tags:"设备管理" summary:"添加设备"`
	Id     string `json:"id" v:"required" dc:"设备ID"`
	Name   string `json:"name" v:"required" dc:"设备名称"`
}

type DeviceUpdateReq struct {
	g.Meta   `path:"/devices/{deviceId}" method:"put" tags:"设备管理" summary:"更新设备"`
	DeviceId string `json:"deviceId" v:"required" dc:"设备ID"`
	Name     string `json:"name" v:"required" dc:"设备名称"`
	Status   string `json:"status" dc:"设备状态"`
}

type DeviceDeleteReq struct {
	g.Meta   `path:"/devices/{deviceId}" method:"delete" tags:"设备管理" summary:"删除设备"`
	DeviceId string `json:"deviceId" v:"required" dc:"设备ID"`
}

type DeviceStatusReq struct {
	g.Meta   `path:"/devices/{deviceId}/status" method:"get" tags:"设备管理" summary:"获取设备状态"`
	DeviceId string `json:"deviceId" v:"required" dc:"设备ID"`
}

// 响应结构体
type DeviceListRes struct {
	Code    int           `json:"code"`
	Message string        `json:"message"`
	Data    []DeviceModel `json:"data"`
}

type DeviceGetRes struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    DeviceModel `json:"data"`
}

type DeviceAddRes DeviceModel

type DeviceUpdateRes DeviceModel

type DeviceDeleteRes struct {
	Success bool `json:"success"`
}

type DeviceStatusRes struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		Status     string    `json:"status"`
		LastActive time.Time `json:"lastActive"`
	} `json:"data"`
}

type DeviceHistoryImageReq struct {
	g.Meta    `path:"/devices/{deviceId}/images" method:"get" tags:"设备管理" summary:"获取设备历史图像"`
	DeviceId  string `json:"deviceId" v:"required" dc:"设备ID"`
	StartTime string `json:"startTime" v:"required" dc:"开始时间"`
	EndTime   string `json:"endTime" v:"required" dc:"结束时间"`
}

type DeviceHistoryImageRes struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    []struct {
		Timestamp string `json:"timestamp"`
		ImageData string `json:"imageData"`
	} `json:"data"`
}

type DeviceRealtimeImageReq struct {
	g.Meta   `path:"/devices/{deviceId}/realtime" method:"get" tags:"设备管理" summary:"获取设备实时图像"`
	DeviceId string `json:"deviceId" v:"required" dc:"设备ID"`
}

type DeviceRealtimeImageRes struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		ImageData string `json:"imageData"`
	} `json:"data"`
}

// 设备数据访问对象
type DeviceDao struct{}

var Device = new(DeviceDao)

// 获取所有设备
func (dao *DeviceDao) List(ctx g.Ctx) (devices []DeviceModel, err error) {
	log.Printf("正在从数据库获取设备列表...")
	
	// 检查数据库连接
	if g.DB() == nil {
		log.Printf("错误：数据库连接未初始化")
		return nil, fmt.Errorf("数据库连接未初始化")
	}
	
	// 初始化空数组
	devices = make([]DeviceModel, 0)
	
	// 执行查询
	err = g.DB().Model("device").Scan(&devices)
	if err != nil {
		log.Printf("数据库查询失败: %v", err)
		return nil, fmt.Errorf("数据库查询失败: %v", err)
	}
	
	// 检查结果
	if len(devices) == 0 {
		log.Printf("数据库中没有找到任何设备")
	} else {
		log.Printf("从数据库中成功获取到 %d 个设备", len(devices))
		for i, device := range devices {
			log.Printf("设备记录 %d: ID=%s, Name=%s, Status=%s, LastActive=%v",
				i+1, device.Id, device.Name, device.Status, device.LastActive)
		}
	}
	
	return devices, nil
}

// 获取单个设备
func (dao *DeviceDao) Get(ctx g.Ctx, id string) (device *DeviceModel, err error) {
	err = g.DB().Model("device").Where("id", id).Scan(&device)
	if err != nil {
		log.Printf("获取设备信息失败: %v", err)
		return nil, err
	}
	
	// 检查设备状态
	if device != nil {
		// 计算时间差
		now := time.Now()
		timeDiff := now.Sub(device.LastActive)
		log.Printf("设备 %s 状态检查:", id)
		log.Printf("- 最后活跃时间: %v", device.LastActive)
		log.Printf("- 当前时间: %v", now)
		log.Printf("- 时间差: %.2f秒", timeDiff.Seconds())
		log.Printf("- 当前状态: %s", device.Status)
		
		// 如果最后活跃时间超过15秒，将状态设置为离线
		if timeDiff > 15*time.Second {
			device.Status = "offline"
			log.Printf("设备 %s 已超过15秒无活动，标记为离线", id)
			// 更新数据库中的状态
			_ = dao.Update(ctx, id, g.Map{
				"status": "offline",
			})
		} else if device.Status == "offline" && timeDiff <= 15*time.Second {
			// 如果设备状态为离线，但最后活跃时间在15秒内，则更新为在线
			device.Status = "online"
			log.Printf("设备 %s 在15秒内有活动，标记为在线", id)
			_ = dao.Update(ctx, id, g.Map{
				"status": "online",
			})
		}
		
		log.Printf("设备 %s 最终状态: %s", id, device.Status)
	}
	
	return device, nil
}

// 添加设备
func (dao *DeviceDao) Add(ctx g.Ctx, device *DeviceModel) error {
	device.CreatedAt = time.Now()
	device.UpdatedAt = time.Now()
	_, err := g.DB().Model("device").Data(device).Insert()
	return err
}

// 更新设备
func (dao *DeviceDao) Update(ctx g.Ctx, id string, data g.Map) error {
	data["updated_at"] = time.Now()
	_, err := g.DB().Model("device").Where("id", id).Data(data).Update()
	return err
}

// 删除设备
func (dao *DeviceDao) Delete(ctx g.Ctx, id string) error {
	_, err := g.DB().Model("device").Where("id", id).Delete()
	return err
}

// 更新设备状态
func (dao *DeviceDao) UpdateStatus(ctx g.Ctx, id string, status string) error {
	now := time.Now()
	log.Printf("更新设备 %s 状态为 %s，时间：%v", id, status, now)
	
	_, err := g.DB().Model("device").
		Where("id", id).
		Data(g.Map{
			"status":      status,
			"last_active": now,
			"updated_at":  now,
		}).
		Update()
		
	if err != nil {
		log.Printf("更新设备状态失败: %v", err)
	}
	return err
}

// 初始化数据库表
func (dao *DeviceDao) InitTable(ctx g.Ctx) error {
	sql := `
	CREATE TABLE IF NOT EXISTS device (
		id VARCHAR(64) PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		status VARCHAR(20) NOT NULL DEFAULT 'offline',
		last_active DATETIME,
		created_at DATETIME,
		updated_at DATETIME,
		INDEX idx_status (status),
		INDEX idx_last_active (last_active)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
	`
	_, err := g.DB().Exec(ctx, sql)
	return err
}

// 获取最新图像
func (dao *DeviceDao) GetLatestImage(ctx g.Ctx, deviceId string) (string, error) {
	log.Printf("获取设备实时图像: %s", deviceId)
	
	// 获取设备状态
	device, err := dao.Get(ctx, deviceId)
	if err != nil {
		log.Printf("获取设备信息失败: %v", err)
		return "", fmt.Errorf("获取设备信息失败: %v", err)
	}

	if device == nil {
		log.Printf("设备不存在: %s", deviceId)
		return "", fmt.Errorf("设备不存在")
	}

	// 获取当前工作目录
	workDir, err := filepath.Abs(".")
	if err != nil {
		log.Printf("获取工作目录失败: %v", err)
		return "", fmt.Errorf("获取工作目录失败: %v", err)
	}

	// 使用绝对路径查找图像文件
	pattern := filepath.Join(workDir, "images", deviceId, "*.jpg")
	log.Printf("查找图像文件: %s", pattern)
	
	files, err := filepath.Glob(pattern)
	if err != nil {
		log.Printf("查找图像文件失败: %v", err)
		return "", fmt.Errorf("查找图像文件失败: %v", err)
	}

	if len(files) == 0 {
		log.Printf("设备 %s 未找到图像文件", deviceId)
		return "", fmt.Errorf("未找到图像文件")
	}

	// 获取最新的图像文件
	latestFile := files[len(files)-1]
	log.Printf("找到最新图像文件: %s", latestFile)

	// 读取图像文件
	imageData, err := ioutil.ReadFile(latestFile)
	if err != nil {
		log.Printf("读取图像文件失败: %v", err)
		return "", fmt.Errorf("读取图像文件失败: %v", err)
	}

	// 转换为base64
	base64Data := base64.StdEncoding.EncodeToString(imageData)
	log.Printf("成功读取图像文件 %s，大小: %d 字节", latestFile, len(imageData))

	return base64Data, nil
}

// 获取历史图像
func (dao *DeviceDao) GetHistoryImages(ctx g.Ctx, deviceId string, startTime string, endTime string) ([]struct {
	Timestamp string `json:"timestamp"`
	ImageData string `json:"imageData"`
}, error) {
	log.Printf("开始获取历史图像，设备ID: %s, 时间范围: %s 至 %s", deviceId, startTime, endTime)
	
	// 解析时间范围
	loc := time.Local // 使用本地时区
	start, err := time.ParseInLocation("2006-01-02 15:04:05", startTime, loc)
	if err != nil {
		return nil, fmt.Errorf("解析开始时间失败: %v", err)
	}
	
	end, err := time.ParseInLocation("2006-01-02 15:04:05", endTime, loc)
	if err != nil {
		return nil, fmt.Errorf("解析结束时间失败: %v", err)
	}
	
	log.Printf("解析后的时间范围: %v 至 %v", start.Format("2006-01-02 15:04:05"), end.Format("2006-01-02 15:04:05"))
	
	// 获取当前工作目录
	workDir, err := filepath.Abs(".")
	if err != nil {
		return nil, fmt.Errorf("获取工作目录失败: %v", err)
	}
	log.Printf("当前工作目录: %s", workDir)
	
	// 从文件系统获取指定时间范围内的图像
	pattern := filepath.Join(workDir, "images", deviceId, "*.jpg")
	log.Printf("查找图像文件: %s", pattern)
	files, err := filepath.Glob(pattern)
	if err != nil {
		return nil, fmt.Errorf("查找图像文件失败: %v", err)
	}
	
	log.Printf("找到 %d 个图像文件", len(files))
	
	var images []struct {
		Timestamp string `json:"timestamp"`
		ImageData string `json:"imageData"`
	}
	
	for _, file := range files {
		// 从文件名解析时间戳
		timestamp := filepath.Base(file)
		timestamp = timestamp[:len(timestamp)-4] // 移除 .jpg 后缀
		fileTime, err := time.ParseInLocation("20060102_150405", timestamp, loc)
		if err != nil {
			log.Printf("解析文件时间戳失败: %s, %v", timestamp, err)
			continue
		}
		
		// 检查时间范围
		if fileTime.Before(start) || fileTime.After(end) {
			log.Printf("文件 %s 不在时间范围内: %v", file, fileTime.Format("2006-01-02 15:04:05"))
			continue
		}
		
		log.Printf("处理文件: %s, 时间: %v", file, fileTime.Format("2006-01-02 15:04:05"))
		
		// 读取图像文件
		imageData, err := ioutil.ReadFile(file)
		if err != nil {
			log.Printf("读取图像文件失败: %v", err)
			continue
		}
		
		// 添加到结果列表
		images = append(images, struct {
			Timestamp string `json:"timestamp"`
			ImageData string `json:"imageData"`
		}{
			Timestamp: fileTime.Format("2006-01-02 15:04:05"),
			ImageData: base64.StdEncoding.EncodeToString(imageData),
		})
		log.Printf("成功添加图像: %s", fileTime.Format("2006-01-02 15:04:05"))
	}
	
	log.Printf("总共找到 %d 张符合条件的图像", len(images))
	return images, nil
} 