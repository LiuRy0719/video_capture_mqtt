package controller

import (
	"context"
	"fmt"
	"log"
	"time"
	"video-platform/internal/model"

	"github.com/gogf/gf/v2/frame/g"
)

var DeviceController = new(deviceController)

type deviceController struct{}

// 获取设备列表
func (c *deviceController) List(ctx context.Context, req *model.DeviceListReq) (res *model.DeviceListRes, err error) {
	log.Printf("正在获取设备列表...")
	devices, err := model.Device.List(ctx)
	if err != nil {
		log.Printf("获取设备列表失败: %v", err)
		r := g.RequestFromCtx(ctx)
		r.Response.WriteJson(g.Map{
			"code": 1,
			"message": fmt.Sprintf("获取设备列表失败: %v", err),
			"data": []model.DeviceModel{},
		})
		return nil, nil
	}
	
	log.Printf("成功获取设备列表，共 %d 个设备", len(devices))
	for i, device := range devices {
		log.Printf("设备 %d: ID=%s, Name=%s, Status=%s", i+1, device.Id, device.Name, device.Status)
	}
	
	// 确保返回一个空数组而不是 nil
	if devices == nil {
		devices = make([]model.DeviceModel, 0)
	}
	
	r := g.RequestFromCtx(ctx)
	r.Response.WriteJson(g.Map{
		"code": 0,
		"message": "success",
		"data": devices,
	})
	return nil, nil
}

// 获取单个设备信息
func (c *deviceController) Get(ctx context.Context, req *model.DeviceGetReq) (res *model.DeviceGetRes, err error) {
	device, err := model.Device.Get(ctx, req.DeviceId)
	if err != nil {
		return nil, fmt.Errorf("获取设备信息失败: %v", err)
	}
	
	return &model.DeviceGetRes{
		Code: 0,
		Message: "success",
		Data: *device,
	}, nil
}

// 添加设备
func (c *deviceController) Add(ctx context.Context, req *model.DeviceAddReq) (res *model.DeviceAddRes, err error) {
	// 先检查设备是否已存在
	if existingDevice, _ := model.Device.Get(ctx, req.Id); existingDevice != nil {
		return nil, fmt.Errorf("设备ID '%s' 已存在", req.Id)
	}

	device := &model.DeviceModel{
		Id:         req.Id,
		Name:       req.Name,
		Status:     "offline",
		LastActive: time.Now(),
	}
	
	if err = model.Device.Add(ctx, device); err != nil {
		return nil, fmt.Errorf("添加设备失败: %v", err)
	}
	
	result := model.DeviceAddRes(*device)
	return &result, nil
}

// 更新设备信息
func (c *deviceController) Update(ctx context.Context, req *model.DeviceUpdateReq) (res *model.DeviceUpdateRes, err error) {
	data := g.Map{
		"name": req.Name,
	}
	if req.Status != "" {
		data["status"] = req.Status
	}
	
	if err = model.Device.Update(ctx, req.DeviceId, data); err != nil {
		return nil, err
	}
	
	device, err := model.Device.Get(ctx, req.DeviceId)
	if err != nil {
		return nil, err
	}
	*res = model.DeviceUpdateRes(*device)
	return
}

// 删除设备
func (c *deviceController) Delete(ctx context.Context, req *model.DeviceDeleteReq) (res *model.DeviceDeleteRes, err error) {
	log.Printf("删除设备: %s", req.DeviceId)
	if err = model.Device.Delete(ctx, req.DeviceId); err != nil {
		return nil, err
	}
	res = &model.DeviceDeleteRes{Success: true}
	return
}

// 获取设备状态
func (c *deviceController) GetStatus(ctx context.Context, req *model.DeviceStatusReq) (res *model.DeviceStatusRes, err error) {
	log.Printf("获取设备状态: %s", req.DeviceId)
	device, err := model.Device.Get(ctx, req.DeviceId)
	if err != nil {
		r := g.RequestFromCtx(ctx)
		r.Response.WriteJson(g.Map{
			"code":    1,
			"message": fmt.Sprintf("获取设备状态失败: %v", err),
			"data": g.Map{
				"status":     "unknown",
				"lastActive": time.Now().Format(time.RFC3339),
			},
		})
		return nil, nil
	}
	
	if device == nil {
		r := g.RequestFromCtx(ctx)
		r.Response.WriteJson(g.Map{
			"code":    1,
			"message": "设备不存在",
			"data": g.Map{
				"status":     "offline",
				"lastActive": time.Now().Format(time.RFC3339),
			},
		})
		return nil, nil
	}
	
	// 确保时区信息正确
	lastActive := device.LastActive.In(time.Local)
	log.Printf("设备 %s 最后活跃时间: %v", req.DeviceId, lastActive.Format(time.RFC3339))
	
	r := g.RequestFromCtx(ctx)
	r.Response.WriteJson(g.Map{
		"code":    0,
		"message": "success",
		"data": g.Map{
			"status":     device.Status,
			"lastActive": lastActive.Format(time.RFC3339),
		},
	})
	return nil, nil
}

// 获取设备实时图像
func (c *deviceController) GetRealtimeImage(ctx context.Context, req *model.DeviceRealtimeImageReq) (res *model.DeviceRealtimeImageRes, err error) {
	log.Printf("获取设备实时图像: %s", req.DeviceId)
	
	// 获取设备最新图像
	imageData, err := model.Device.GetLatestImage(ctx, req.DeviceId)
	if err != nil {
		log.Printf("获取实时图像失败: %v", err)
		r := g.RequestFromCtx(ctx)
		r.Response.WriteJson(g.Map{
			"code":    1,
			"message": fmt.Sprintf("获取实时图像失败: %v", err),
			"data": g.Map{
				"imageData": "",
			},
		})
		return nil, nil
	}
	
	log.Printf("成功获取实时图像，数据长度: %d", len(imageData))
	r := g.RequestFromCtx(ctx)
	r.Response.WriteJson(g.Map{
		"code":    0,
		"message": "success",
		"data": g.Map{
			"imageData": imageData,
		},
	})
	return nil, nil
}

// 获取设备历史图像
func (c *deviceController) GetHistoryImages(ctx context.Context, req *model.DeviceHistoryImageReq) (res *model.DeviceHistoryImageRes, err error) {
	log.Printf("获取设备历史图像: %s, 时间范围: %s - %s", req.DeviceId, req.StartTime, req.EndTime)
	
	// 获取历史图像列表
	images, err := model.Device.GetHistoryImages(ctx, req.DeviceId, req.StartTime, req.EndTime)
	if err != nil {
		log.Printf("获取历史图像失败: %v", err)
		r := g.RequestFromCtx(ctx)
		r.Response.WriteJson(g.Map{
			"code":    1,
			"message": fmt.Sprintf("获取历史图像失败: %v", err),
			"data":    []interface{}{},
		})
		return nil, nil
	}
	
	log.Printf("成功获取历史图像，共 %d 张", len(images))
	r := g.RequestFromCtx(ctx)
	r.Response.WriteJson(g.Map{
		"code":    0,
		"message": "success",
		"data":    images,
	})
	return nil, nil
} 