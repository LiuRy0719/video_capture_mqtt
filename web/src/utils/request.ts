import axios from "axios";
import type { AxiosRequestConfig } from "axios";

const request = axios.create({
  baseURL: "/api/v1",
  timeout: 5000,
  headers: {
    "Content-Type": "application/json",
    Accept: "application/json",
  },
});

// 请求拦截器
request.interceptors.request.use(
  (config) => {
    console.log("发送请求:", {
      method: config.method?.toUpperCase(),
      url: config.url,
      headers: config.headers,
      data: config.data,
    });
    return config;
  },
  (error) => {
    console.error("请求错误:", error);
    return Promise.reject(error);
  }
);

// 响应拦截器
request.interceptors.response.use(
  (response) => {
    console.log("收到响应:", {
      url: response.config.url,
      status: response.status,
      headers: response.headers,
      data: response.data,
    });
    return response;
  },
  (error) => {
    console.error("响应错误:", {
      url: error.config?.url,
      status: error.response?.status,
      message: error.message,
      response: error.response?.data,
    });
    return Promise.reject(error);
  }
);

export default request;
