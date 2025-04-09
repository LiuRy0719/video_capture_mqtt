<template>
  <div class="device-list">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>设备列表</span>
          <el-button type="primary" @click="handleAddDevice"
            >添加设备</el-button
          >
        </div>
      </template>
      <el-table :data="deviceList" style="width: 100%" v-loading="loading">
        <el-table-column prop="id" label="设备ID" width="180" />
        <el-table-column prop="name" label="设备名称" width="180" />
        <el-table-column prop="status" label="状态">
          <template #default="{ row }">
            <el-tag :type="row.status === 'online' ? 'success' : 'danger'">
              {{ row.status === "online" ? "在线" : "离线" }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="lastActive" label="最后活跃时间" />
        <el-table-column label="操作" width="280">
          <template #default="{ row }">
            <el-button type="primary" link @click="handleMonitor(row.id)"
              >实时监控</el-button
            >
            <el-button type="primary" link @click="handleHistory(row.id)"
              >历史记录</el-button
            >
            <el-button type="danger" link @click="handleDelete(row.id)"
              >删除</el-button
            >
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- 添加设备对话框 -->
    <el-dialog v-model="dialogVisible" title="添加设备" width="30%">
      <el-form :model="formData" label-width="80px">
        <el-form-item label="设备ID">
          <el-input v-model="formData.deviceId" placeholder="请输入设备ID" />
        </el-form-item>
        <el-form-item label="设备名称">
          <el-input v-model="formData.name" placeholder="请输入设备名称" />
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="dialogVisible = false">取消</el-button>
          <el-button type="primary" @click="submitDevice">确定</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted, reactive } from "vue";
import { useRouter } from "vue-router";
import { ElMessage, ElMessageBox } from "element-plus";
import * as deviceApi from "@/api/device";
import type { Device } from "@/api/device";

const router = useRouter();
const deviceList = ref<Device[]>([]);
const dialogVisible = ref(false);
const formData = reactive({
  deviceId: "",
  name: "",
});
const loading = ref(false);
let refreshInterval: number | null = null;

// 获取设备列表
const fetchDevices = async () => {
  try {
    loading.value = true;
    const response = await deviceApi.getDevices();
    if (response.data.code === 0) {
      deviceList.value = response.data.data;
    } else {
      ElMessage.error(response.data.message || "获取设备列表失败");
    }
  } catch (error) {
    console.error("获取设备列表失败:", error);
    ElMessage.error("获取设备列表失败");
  } finally {
    loading.value = false;
  }
};

// 添加设备
const handleAddDevice = () => {
  formData.deviceId = "";
  formData.name = "";
  dialogVisible.value = true;
};

const submitDevice = async () => {
  if (!formData.deviceId || !formData.name) {
    ElMessage.warning("请填写完整信息");
    return;
  }

  try {
    const response = await deviceApi.addDevice(
      formData.deviceId,
      formData.name
    );
    if (response.data.code === 0) {
      ElMessage.success("添加成功");
      dialogVisible.value = false;
      formData.deviceId = "";
      formData.name = "";
      await fetchDevices();
    } else {
      ElMessage.error(response.data.message || "添加失败");
    }
  } catch (error) {
    console.error("添加设备失败:", error);
    ElMessage.error("添加设备失败");
  }
};

// 删除设备
const handleDelete = async (deviceId: string) => {
  try {
    await ElMessageBox.confirm("确认删除该设备吗？", "提示", {
      type: "warning",
    });
    const response = await deviceApi.deleteDevice(deviceId);
    if (response.data.code === 0) {
      ElMessage.success("删除成功");
      await fetchDevices();
    } else {
      ElMessage.error(response.data.message || "删除失败");
    }
  } catch (error) {
    if (error !== "cancel") {
      console.error("删除设备失败:", error);
      ElMessage.error("删除设备失败");
    }
  }
};

const handleMonitor = (deviceId: string) => {
  router.push(`/monitor/realtime?deviceId=${deviceId}`);
};

const handleHistory = (deviceId: string) => {
  router.push(`/monitor/history?deviceId=${deviceId}`);
};

// 启动自动刷新
const startAutoRefresh = () => {
  // 每10秒刷新一次设备列表
  refreshInterval = window.setInterval(fetchDevices, 10000);
};

// 停止自动刷新
const stopAutoRefresh = () => {
  if (refreshInterval) {
    clearInterval(refreshInterval);
    refreshInterval = null;
  }
};

onMounted(async () => {
  await fetchDevices();
  startAutoRefresh();
});

onUnmounted(() => {
  stopAutoRefresh();
});
</script>

<style scoped>
.device-list {
  padding: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
}
</style>
