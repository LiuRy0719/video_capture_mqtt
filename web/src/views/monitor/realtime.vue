<template>
  <div class="realtime-monitor">
    <el-card>
      <template #header>
        <div class="card-header">
          <div class="header-left">
            <span>实时监控 - 设备 {{ deviceId }}</span>
            <div class="layout-controls">
              <el-radio-group v-model="screenLayout" size="small">
                <el-radio-button label="1">全屏</el-radio-button>
                <el-radio-button label="2">两分屏</el-radio-button>
                <el-radio-button label="4">四分屏</el-radio-button>
                <el-radio-button label="8">八分屏</el-radio-button>
                <el-radio-button label="16">十六分屏</el-radio-button>
              </el-radio-group>
            </div>
          </div>
          <div class="header-controls">
            <el-tag
              :type="deviceStatus === 'online' ? 'success' : 'danger'"
              class="device-status"
            >
              {{ deviceStatus === "online" ? "设备在线" : "设备离线" }}
            </el-tag>
            <el-button type="primary" @click="toggleMonitor">
              {{ isMonitoring ? "停止监控" : "开始监控" }}
            </el-button>
          </div>
        </div>
      </template>
      <div class="monitor-content">
        <div class="image-grid" :class="gridLayoutClass">
          <div v-for="n in screenLayout" :key="n" class="grid-item">
            <div class="image-container">
              <img
                v-if="screenImages[n]"
                :src="screenImages[n]"
                :alt="'实时图像 ' + n"
              />
              <div v-else class="no-image">
                {{ isMonitoring ? "等待图像..." : "未开始监控" }}
              </div>
            </div>
            <div class="grid-item-title">摄像头 {{ n }}</div>
          </div>
        </div>
        <div class="info-panel">
          <el-descriptions title="设备信息" :column="1" border>
            <el-descriptions-item label="设备ID">{{
              deviceId
            }}</el-descriptions-item>
            <el-descriptions-item label="状态">
              <el-tag :type="isMonitoring ? 'success' : 'info'">
                {{ isMonitoring ? "监控中" : "未监控" }}
              </el-tag>
            </el-descriptions-item>
            <el-descriptions-item label="最后更新">{{
              lastUpdate
            }}</el-descriptions-item>
            <el-descriptions-item label="图像计数">{{
              imageCount
            }}</el-descriptions-item>
            <el-descriptions-item label="分屏模式">
              {{ screenLayout }}分屏
            </el-descriptions-item>
          </el-descriptions>
        </div>
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted, watch } from "vue";
import { useRoute } from "vue-router";
import { ElMessage } from "element-plus";
import * as deviceApi from "@/api/device";

const route = useRoute();
const deviceId = route.query.deviceId as string;
const isMonitoring = ref(false);
const lastUpdate = ref("-");
const imageCount = ref(0);
const deviceStatus = ref("offline");
const screenLayout = ref(1); // 默认全屏

// 为每个分屏创建独立的图片状态
const screenImages = ref<{ [key: number]: string }>({});

let pollingInterval: number | null = null;
let statusInterval: number | null = null;

// 计算分屏布局的类名
const gridLayoutClass = computed(() => {
  return {
    "grid-1": screenLayout.value === 1,
    "grid-2": screenLayout.value === 2,
    "grid-4": screenLayout.value === 4,
    "grid-8": screenLayout.value === 8,
    "grid-16": screenLayout.value === 16,
  };
});

// 获取最新图像
const fetchLatestImage = async () => {
  try {
    const response = await deviceApi.getRealtimeImage(deviceId);
    console.log("获取实时图像响应:", response.data);

    if (response.data && response.data.code === 0) {
      const imageData = response.data.data.imageData;
      if (imageData) {
        // 为每个分屏设置不同的延迟，模拟多个摄像头
        const base64Image = `data:image/jpeg;base64,${imageData}`;
        for (let i = 1; i <= screenLayout.value; i++) {
          setTimeout(() => {
            screenImages.value[i] = base64Image;
          }, i * 200); // 每个分屏延迟200ms更新
        }
        lastUpdate.value = new Date().toLocaleString();
        imageCount.value++;
      } else {
        console.warn("图像数据为空");
      }
    } else {
      console.warn("获取图像失败:", response.data?.message || "未知错误");
    }
  } catch (error) {
    console.error("获取实时图像失败:", error);
  }
};

// 获取设备状态
const fetchDeviceStatus = async () => {
  try {
    console.log("正在获取设备状态...");
    const response = await deviceApi.getDeviceStatus(deviceId);
    console.log("设备状态响应:", response.data);

    if (response.data && response.data.code === 0) {
      const newStatus = response.data.data.status;
      // 解析时间字符串，保持原始时区
      const lastActiveStr = response.data.data.lastActive;
      // 将时间字符串转换为本地时间
      const lastActive = new Date(lastActiveStr);
      const now = new Date();

      // 获取时区偏移量（毫秒）
      const tzOffset = now.getTimezoneOffset() * 60 * 1000;

      // 计算考虑时区的时间差（秒）
      const timeDiff =
        Math.abs(now.getTime() - (lastActive.getTime() - tzOffset)) / 1000;

      console.log(`原始最后活跃时间字符串: ${lastActiveStr}`);
      console.log(`解析后的最后活跃时间: ${lastActive.toISOString()}`);
      console.log(`当前时间: ${now.toISOString()}`);
      console.log(`时区偏移（分钟）: ${now.getTimezoneOffset()}`);
      console.log(`计算的时间差: ${timeDiff}秒`);

      // 如果最后活跃时间在15秒内，则认为设备在线
      const isDeviceActive = timeDiff <= 15;
      const effectiveStatus = isDeviceActive ? "online" : "offline";

      console.log(
        `计算得到的实际状态: ${effectiveStatus} (基于时间差 ${timeDiff}秒)`
      );

      // 如果设备从在线变为离线，自动停止监控
      if (
        deviceStatus.value === "online" &&
        effectiveStatus === "offline" &&
        isMonitoring.value
      ) {
        console.log("检测到设备离线，停止监控");
        ElMessage.warning("设备已离线，停止监控");
        stopMonitoring();
      }

      deviceStatus.value = effectiveStatus;
    } else {
      console.warn("获取设备状态失败:", response.data?.message || "未知错误");
    }
  } catch (error) {
    console.error("获取设备状态失败:", error);
  }
};

const toggleMonitor = () => {
  if (isMonitoring.value) {
    stopMonitoring();
  } else {
    startMonitoring();
  }
};

const startMonitoring = async () => {
  try {
    console.log("正在启动监控...");
    await fetchDeviceStatus();
    console.log("当前设备状态:", deviceStatus.value);

    if (deviceStatus.value !== "online") {
      console.warn("设备不在线，无法开始监控");
      ElMessage.warning("设备当前不在线，无法开始监控");
      return;
    }

    isMonitoring.value = true;
    console.log("开始获取第一帧图像...");
    // 立即获取一次图像
    await fetchLatestImage();
    console.log("开始轮询获取新图像...");
    // 开始轮询获取新图像
    pollingInterval = window.setInterval(fetchLatestImage, 1000);
    ElMessage.success("开始监控");
  } catch (error) {
    ElMessage.error("启动监控失败");
    console.error("启动监控失败:", error);
  }
};

const stopMonitoring = () => {
  isMonitoring.value = false;
  if (pollingInterval) {
    clearInterval(pollingInterval);
    pollingInterval = null;
  }
  // 清除所有分屏的图片
  screenImages.value = {};
  ElMessage.info("已停止监控");
};

// 启动状态定时更新
const startStatusPolling = () => {
  // 立即获取一次状态
  fetchDeviceStatus();
  // 每2秒更新一次状态
  statusInterval = window.setInterval(fetchDeviceStatus, 2000);
};

// 停止状态定时更新
const stopStatusPolling = () => {
  if (statusInterval) {
    clearInterval(statusInterval);
    statusInterval = null;
  }
};

// 监听分屏布局变化
watch(screenLayout, (newLayout) => {
  // 清除超出当前布局的图片
  const currentImages = { ...screenImages.value };
  for (const key in currentImages) {
    if (parseInt(key) > newLayout) {
      delete screenImages.value[key];
    }
  }
});

onMounted(async () => {
  if (!deviceId) {
    ElMessage.error("设备ID不能为空");
    return;
  }

  console.log("页面加载，设备ID:", deviceId);
  try {
    // 获取设备信息和状态
    console.log("正在获取设备信息...");
    const deviceInfo = await deviceApi.getDevice(deviceId);
    console.log("设备信息:", deviceInfo.data);

    console.log("正在获取初始状态...");
    await fetchDeviceStatus();
    console.log("启动状态定时更新...");
    // 启动状态定时更新
    startStatusPolling();
  } catch (error) {
    ElMessage.error("获取设备信息失败");
    console.error("初始化失败:", error);
  }
});

onUnmounted(() => {
  stopMonitoring();
  stopStatusPolling();
});
</script>

<style scoped>
.realtime-monitor {
  padding: 20px;
  height: calc(100vh - 40px); /* 减去padding的高度 */
}

.el-card {
  height: 100%;
  display: flex;
  flex-direction: column;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 10px 0;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 20px;
}

.layout-controls {
  margin-left: 20px;
}

.header-controls {
  display: flex;
  gap: 16px;
  align-items: center;
}

.device-status {
  min-width: 100px;
  text-align: center;
}

.monitor-content {
  flex: 1;
  display: flex;
  gap: 20px;
  overflow: hidden; /* 防止内容溢出 */
  padding: 10px 0;
}

.image-grid {
  flex: 1;
  display: grid;
  gap: 10px;
  height: 100%;
  padding: 10px;
}

.grid-item {
  position: relative;
  width: 100%;
  height: 100%;
  background-color: #f5f7fa;
  border-radius: 4px;
  overflow: hidden;
}

.grid-item-title {
  position: absolute;
  top: 10px;
  left: 10px;
  background-color: rgba(0, 0, 0, 0.5);
  color: white;
  padding: 4px 8px;
  border-radius: 4px;
  font-size: 12px;
}

.grid-1 {
  grid-template-columns: 1fr;
}

.grid-2 {
  grid-template-columns: repeat(2, 1fr);
}

.grid-4 {
  grid-template-columns: repeat(2, 1fr);
  grid-template-rows: repeat(2, 1fr);
}

.grid-8 {
  grid-template-columns: repeat(4, 1fr);
  grid-template-rows: repeat(2, 1fr);
}

.grid-16 {
  grid-template-columns: repeat(4, 1fr);
  grid-template-rows: repeat(4, 1fr);
}

.image-container {
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
}

.image-container img {
  width: 100%;
  height: 100%;
  object-fit: contain;
  transition: opacity 0.3s ease-in-out;
}

.image-container img:not([src]) {
  opacity: 0;
}

.image-container img[src] {
  opacity: 1;
}

.no-image {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 100%;
  height: 100%;
  color: #909399;
  font-size: 16px;
}

.info-panel {
  width: 300px;
  min-width: 300px;
  overflow-y: auto; /* 如果内容过多允许滚动 */
}

:deep(.el-descriptions) {
  padding: 10px;
}

:deep(.el-descriptions__body) {
  background-color: #fff;
}

/* 响应式布局调整 */
@media (max-width: 1200px) {
  .grid-8,
  .grid-16 {
    grid-template-columns: repeat(2, 1fr);
    grid-template-rows: auto;
  }
}

@media (max-width: 768px) {
  .header-left {
    flex-direction: column;
    align-items: flex-start;
    gap: 10px;
  }

  .layout-controls {
    margin-left: 0;
  }

  .monitor-content {
    flex-direction: column;
  }

  .image-grid {
    grid-template-columns: 1fr !important;
    grid-template-rows: auto !important;
  }

  .info-panel {
    width: 100%;
    min-width: unset;
  }
}
</style>
