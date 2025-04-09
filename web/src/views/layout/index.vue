<template>
  <el-container class="layout-container">
    <el-aside :width="isCollapse ? '64px' : '200px'" class="aside">
      <div class="logo" :class="{ 'logo-collapse': isCollapse }">
        <img src="../../assets/police-badge.svg" class="logo-img" alt="警徽" />
        <div class="logo-text" v-show="!isCollapse">
          <h1>公安视频采集平台</h1>
        </div>
      </div>
      <el-menu
        :default-active="route.path"
        class="el-menu-vertical"
        :collapse="isCollapse"
        background-color="#304156"
        text-color="#bfcbd9"
        active-text-color="#409EFF"
        router
      >
        <el-menu-item index="/dashboard">
          <el-icon><Odometer /></el-icon>
          <template #title>仪表盘</template>
        </el-menu-item>

        <el-menu-item index="/device">
          <el-icon><Monitor /></el-icon>
          <template #title>设备管理</template>
        </el-menu-item>
      </el-menu>
    </el-aside>

    <el-container class="main-container">
      <el-header class="header">
        <div class="header-left">
          <el-icon
            class="fold-btn"
            :class="{ 'is-active': isCollapse }"
            @click="toggleSidebar"
          >
            <Fold v-if="!isCollapse" />
            <Expand v-else />
          </el-icon>

          <el-breadcrumb separator="/">
            <el-breadcrumb-item :to="{ path: '/' }">首页</el-breadcrumb-item>
            <el-breadcrumb-item>{{ route.meta.title }}</el-breadcrumb-item>
          </el-breadcrumb>
        </div>

        <div class="header-right">
          <el-tooltip content="全屏" placement="bottom">
            <el-icon class="action-icon" @click="toggleFullScreen">
              <FullScreen v-if="!isFullscreen" />
              <Aim v-else />
            </el-icon>
          </el-tooltip>

          <el-dropdown>
            <span class="user-info">
              <el-avatar
                :size="32"
                src="https://cube.elemecdn.com/0/88/03b0d39583f48206768a7534e55bcpng.png"
              />
              <span>管理员</span>
            </span>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item>
                  <el-icon><User /></el-icon>个人信息
                </el-dropdown-item>
                <el-dropdown-item divided @click="handleLogout">
                  <el-icon><SwitchButton /></el-icon>退出登录
                </el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </div>
      </el-header>

      <tags-view />

      <el-main>
        <router-view v-slot="{ Component }">
          <transition name="fade-transform" mode="out-in">
            <component :is="Component" />
          </transition>
        </router-view>
      </el-main>
    </el-container>
  </el-container>
</template>

<script setup lang="ts">
import { ref } from "vue";
import { useRoute, useRouter } from "vue-router";
import TagsView from "../../components/TagsView.vue";
import { ElMessageBox } from "element-plus";

const route = useRoute();
const router = useRouter();
const isCollapse = ref(false);
const isFullscreen = ref(false);

const toggleSidebar = () => {
  isCollapse.value = !isCollapse.value;
};

const toggleFullScreen = () => {
  if (!document.fullscreenElement) {
    document.documentElement.requestFullscreen();
    isFullscreen.value = true;
  } else {
    if (document.exitFullscreen) {
      document.exitFullscreen();
      isFullscreen.value = false;
    }
  }
};

const handleLogout = () => {
  ElMessageBox.confirm("确认退出登录吗？", "提示", {
    confirmButtonText: "确定",
    cancelButtonText: "取消",
    type: "warning",
  }).then(() => {
    localStorage.removeItem("token");
    router.push("/login");
  });
};
</script>

<style scoped>
.layout-container {
  height: 100vh;
}

.aside {
  background-color: #304156;
  transition: width 0.3s;
  overflow: hidden;
}

.logo {
  height: 56px;
  display: flex;
  align-items: center;
  padding: 0 16px;
  background: linear-gradient(to right, #001529 0%, #002140 100%);
  color: #fff;
  transition: all 0.3s;
  white-space: nowrap;
  overflow: hidden;
  box-shadow: inset 0 -1px 0 0 rgba(255, 255, 255, 0.1);
}

.logo-img {
  width: 32px;
  height: 32px;
  margin-right: 10px;
  transition: all 0.3s;
  filter: drop-shadow(0 0 3px rgba(64, 158, 255, 0.6));
  flex-shrink: 0;
  animation: pulse 2s infinite;
}

.logo-text {
  flex: 1;
  min-width: 0;
  overflow: hidden;
  position: relative;
}

.logo-text h1 {
  font-size: 15px;
  margin: 0;
  font-weight: 600;
  line-height: 1.2;
  white-space: nowrap;
  overflow: hidden;
  position: relative;
  animation: shine 4s linear infinite;
  background: linear-gradient(
    90deg,
    #ffffff 0%,
    #ffffff 35%,
    #409eff 50%,
    #ffffff 65%,
    #ffffff 100%
  );
  background-size: 400% 100%;
  -webkit-background-clip: text;
  background-clip: text;
  color: transparent;
}

.logo-text h1::after {
  content: "";
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: linear-gradient(
    90deg,
    transparent 0%,
    rgba(64, 158, 255, 0.15) 25%,
    rgba(64, 158, 255, 0.3) 50%,
    rgba(64, 158, 255, 0.15) 75%,
    transparent 100%
  );
  animation: shine 4s linear infinite;
  background-size: 400% 100%;
}

@keyframes shine {
  0% {
    background-position: 100% 0;
  }
  100% {
    background-position: -300% 0;
  }
}

@keyframes pulse {
  0% {
    filter: drop-shadow(0 0 3px rgba(64, 158, 255, 0.6));
    transform: scale(1);
  }
  50% {
    filter: drop-shadow(0 0 8px rgba(64, 158, 255, 0.8));
    transform: scale(1.05);
  }
  100% {
    filter: drop-shadow(0 0 3px rgba(64, 158, 255, 0.6));
    transform: scale(1);
  }
}

.logo-collapse {
  padding: 0;
  justify-content: center;
}

.logo-collapse .logo-img {
  margin-right: 0;
  transform: scale(0.9);
}

.el-menu-vertical {
  border-right: none;
}

.main-container {
  transition: margin-left 0.3s;
}

.header {
  background-color: #fff;
  border-bottom: 1px solid #e6e6e6;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 20px;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 20px;
}

.header-right {
  display: flex;
  align-items: center;
  gap: 20px;
}

.fold-btn {
  font-size: 20px;
  cursor: pointer;
  transition: transform 0.3s;
}

.fold-btn.is-active {
  transform: rotate(180deg);
}

.action-icon {
  font-size: 20px;
  cursor: pointer;
  padding: 8px;
  border-radius: 50%;
  transition: background-color 0.3s;
}

.action-icon:hover {
  background-color: #f5f7fa;
}

.user-info {
  display: flex;
  align-items: center;
  cursor: pointer;
  padding: 5px 10px;
  border-radius: 4px;
  transition: background-color 0.3s;
}

.user-info:hover {
  background-color: #f5f7fa;
}

.user-info span {
  margin-left: 8px;
}

/* 路由过渡动画 */
.fade-transform-enter-active,
.fade-transform-leave-active {
  transition: all 0.5s;
}

.fade-transform-enter-from {
  opacity: 0;
  transform: translateX(-30px);
}

.fade-transform-leave-to {
  opacity: 0;
  transform: translateX(30px);
}
</style>
