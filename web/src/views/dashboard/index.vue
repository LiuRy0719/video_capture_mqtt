<template>
  <div class="dashboard">
    <el-row :gutter="20">
      <el-col :span="6">
        <el-card shadow="hover">
          <template #header>
            <div class="card-header">
              <span>总视频数</span>
              <el-icon><VideoCamera /></el-icon>
            </div>
          </template>
          <div class="card-content">
            <span class="number">0</span>
            <span class="unit">个</span>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="hover">
          <template #header>
            <div class="card-header">
              <span>今日采集</span>
              <el-icon><Calendar /></el-icon>
            </div>
          </template>
          <div class="card-content">
            <span class="number">0</span>
            <span class="unit">个</span>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="hover">
          <template #header>
            <div class="card-header">
              <span>存储空间</span>
              <el-icon><FolderOpened /></el-icon>
            </div>
          </template>
          <div class="card-content">
            <span class="number">0</span>
            <span class="unit">GB</span>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="hover">
          <template #header>
            <div class="card-header">
              <span>在线设备</span>
              <el-icon><Monitor /></el-icon>
            </div>
          </template>
          <div class="card-content">
            <span class="number">0</span>
            <span class="unit">台</span>
          </div>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup lang="ts">
// Dashboard component
</script>

<style scoped>
.dashboard {
  padding: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.card-content {
  text-align: center;
  padding: 20px 0;
}

.number {
  font-size: 36px;
  font-weight: bold;
  color: #303133;
  margin-right: 8px;
}

.unit {
  font-size: 14px;
  color: #909399;
}
</style>
