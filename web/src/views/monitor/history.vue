<template>
  <div class="history-monitor">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>历史记录 - 设备 {{ deviceId }}</span>
          <div class="header-controls">
            <el-date-picker
              v-model="dateRange"
              type="daterange"
              range-separator="至"
              start-placeholder="开始日期"
              end-placeholder="结束日期"
              value-format="YYYY-MM-DD HH:mm:ss"
              @change="handleDateChange"
            />
            <el-select
              v-model="sortOrder"
              placeholder="排序方式"
              style="width: 120px"
            >
              <el-option label="最新在前" value="desc" />
              <el-option label="最早在前" value="asc" />
            </el-select>
            <el-button type="primary" @click="handleSearch" :loading="loading">
              查询
            </el-button>
          </div>
        </div>
      </template>
      <div class="history-list" v-loading="loading">
        <el-empty v-if="displayList.length === 0" description="暂无历史记录" />
        <el-timeline v-else>
          <el-timeline-item
            v-for="item in displayList"
            :key="item.timestamp"
            :timestamp="item.timestamp"
            placement="top"
          >
            <el-card>
              <div class="history-item">
                <img :src="item.imageUrl" alt="历史图像" />
                <div class="item-info">
                  <div class="time-info">
                    <p>
                      <strong>采集时间：</strong
                      >{{ formatDateTime(item.timestamp) }}
                    </p>
                    <p>
                      <strong>距现在：</strong>{{ getTimeAgo(item.timestamp) }}
                    </p>
                  </div>
                  <el-button type="primary" link @click="handlePreview(item)">
                    查看大图
                  </el-button>
                </div>
              </div>
            </el-card>
          </el-timeline-item>
        </el-timeline>

        <!-- 分页控件 -->
        <div class="pagination-container" v-if="historyList.length > 0">
          <el-pagination
            v-model:current-page="currentPage"
            v-model:page-size="pageSize"
            :page-sizes="[10, 20, 50, 100]"
            :total="historyList.length"
            layout="total, sizes, prev, pager, next"
            @size-change="handleSizeChange"
            @current-change="handleCurrentChange"
          />
        </div>
      </div>
    </el-card>

    <!-- 图片预览对话框 -->
    <el-dialog v-model="previewVisible" title="图片预览" width="80%">
      <img :src="previewImage" alt="预览图像" style="width: 100%" />
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from "vue";
import { useRoute } from "vue-router";
import { ElMessage } from "element-plus";
import { deviceApi } from "@/api/device";

interface HistoryImage {
  timestamp: string;
  imageUrl: string;
}

type DateRangeValue = [string, string] | null;

const route = useRoute();
const deviceId = route.query.deviceId as string;
const dateRange = ref<DateRangeValue>(null);
const previewVisible = ref(false);
const previewImage = ref("");
const loading = ref(false);
const historyList = ref<HistoryImage[]>([]);
const sortOrder = ref<"asc" | "desc">("desc");
const currentPage = ref(1);
const pageSize = ref(20);

// 格式化日期时间
const formatDateTime = (timestamp: string) => {
  const date = new Date(timestamp);
  const pad = (n: number) => n.toString().padStart(2, "0");
  return `${date.getFullYear()}-${pad(date.getMonth() + 1)}-${pad(
    date.getDate()
  )} ${pad(date.getHours())}:${pad(date.getMinutes())}:${pad(
    date.getSeconds()
  )}`;
};

// 计算距离现在的时间
const getTimeAgo = (timestamp: string) => {
  const now = new Date().getTime();
  const past = new Date(timestamp).getTime();
  const diff = now - past;

  const minutes = Math.floor(diff / (1000 * 60));
  const hours = Math.floor(minutes / 60);
  const days = Math.floor(hours / 24);

  if (days > 0) return `${days}天前`;
  if (hours > 0) return `${hours}小时前`;
  if (minutes > 0) return `${minutes}分钟前`;
  return "刚刚";
};

// 排序和分页后的列表
const displayList = computed(() => {
  const sorted = [...historyList.value].sort((a, b) => {
    const timeA = new Date(a.timestamp).getTime();
    const timeB = new Date(b.timestamp).getTime();
    return sortOrder.value === "desc" ? timeB - timeA : timeA - timeB;
  });

  const start = (currentPage.value - 1) * pageSize.value;
  const end = start + pageSize.value;
  return sorted.slice(start, end);
});

const handleDateChange = (val: DateRangeValue) => {
  dateRange.value = val;
};

const handleSizeChange = (size: number) => {
  pageSize.value = size;
  currentPage.value = 1; // 重置到第一页
};

const handleCurrentChange = (page: number) => {
  currentPage.value = page;
};

const handleSearch = async () => {
  if (!dateRange.value) {
    ElMessage.warning("请选择日期范围");
    return;
  }

  loading.value = true;
  try {
    const [startTime, endTime] = dateRange.value;
    console.log("发送查询请求:", { deviceId, startTime, endTime });

    const response = await deviceApi.getHistoryImages(
      deviceId,
      startTime,
      endTime
    );

    console.log("收到历史记录响应:", response.data);

    if (
      response.data &&
      response.data.code === 0 &&
      Array.isArray(response.data.data)
    ) {
      console.log("找到历史记录数量:", response.data.data.length);
      historyList.value = response.data.data.map((item: any) => ({
        timestamp: item.timestamp,
        imageUrl: `data:image/jpeg;base64,${item.imageData}`,
      }));

      currentPage.value = 1; // 重置到第一页

      if (historyList.value.length === 0) {
        ElMessage.warning("未找到历史记录");
      } else {
        ElMessage.success(`找到 ${historyList.value.length} 条历史记录`);
      }
    } else {
      console.warn("响应数据格式不正确:", response.data);
      historyList.value = [];
      ElMessage.warning("未找到历史记录");
    }
  } catch (error) {
    console.error("获取历史记录失败:", error);
    ElMessage.error("获取历史记录失败");
    historyList.value = [];
  } finally {
    loading.value = false;
  }
};

const handlePreview = (item: HistoryImage) => {
  previewImage.value = item.imageUrl;
  previewVisible.value = true;
};

onMounted(async () => {
  if (!deviceId) {
    ElMessage.error("设备ID不能为空");
    return;
  }

  // 默认加载最近24小时的历史记录
  const now = new Date();
  const yesterday = new Date(now.getTime() - 24 * 60 * 60 * 1000);

  // 使用本地时间格式化函数
  const formatDate = (date: Date) => {
    const pad = (n: number) => n.toString().padStart(2, "0");
    return `${date.getFullYear()}-${pad(date.getMonth() + 1)}-${pad(
      date.getDate()
    )} ${pad(date.getHours())}:${pad(date.getMinutes())}:${pad(
      date.getSeconds()
    )}`;
  };

  dateRange.value = [formatDate(yesterday), formatDate(now)];

  console.log("查询时间范围:", dateRange.value);
  await handleSearch();
});
</script>

<style scoped>
.history-monitor {
  padding: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.header-controls {
  display: flex;
  gap: 10px;
  align-items: center;
}

.history-list {
  margin-top: 20px;
  min-height: 200px;
}

.history-item {
  display: flex;
  gap: 20px;
}

.history-item img {
  width: 200px;
  height: 150px;
  object-fit: cover;
  border-radius: 4px;
}

.item-info {
  flex: 1;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
}

.time-info p {
  margin: 5px 0;
  color: #606266;
}

.pagination-container {
  margin-top: 20px;
  display: flex;
  justify-content: center;
}
</style>
