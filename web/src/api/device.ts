import request from "@/utils/request";
import type { ApiResponse } from "./types";

export interface Device {
  id: string;
  name: string;
  status: string;
  lastActive: string;
}

export interface DeviceForm {
  deviceId: string;
  name: string;
}

// 获取设备列表
export function getDevices() {
  return request<ApiResponse<Device[]>>({
    url: "/devices",
    method: "get",
  });
}

// 添加设备
export function addDevice(deviceId: string, name: string) {
  return request<ApiResponse<null>>({
    url: "/devices",
    method: "post",
    data: { deviceId, name },
  });
}

// 删除设备
export function deleteDevice(deviceId: string) {
  return request<ApiResponse<null>>({
    url: `/devices/${deviceId}`,
    method: "delete",
  });
}

// 获取设备状态
export function getDeviceStatus(deviceId: string) {
  return request<ApiResponse<{ status: string }>>({
    url: `/devices/${deviceId}/status`,
    method: "get",
  });
}

// 获取设备实时图像
export function getRealtimeImage(deviceId: string) {
  return request<ApiResponse<{ imageData: string }>>({
    url: `/devices/${deviceId}/realtime`,
    method: "get",
  });
}

// 获取设备历史图像
export function getHistoryImages(
  deviceId: string,
  startTime?: string,
  endTime?: string
) {
  return request<
    ApiResponse<{
      images: { id: string; timestamp: string; imageData: string }[];
    }>
  >({
    url: `/devices/${deviceId}/images`,
    method: "get",
    params: { startTime, endTime },
  });
}

// 获取单个设备信息
export function getDevice(deviceId: string) {
  return request<ApiResponse<Device>>({
    url: `/devices/${deviceId}`,
    method: "get",
  });
}
