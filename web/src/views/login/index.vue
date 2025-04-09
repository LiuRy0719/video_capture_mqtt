<template>
  <div class="login-container">
    <el-form
      ref="loginFormRef"
      :model="loginForm"
      :rules="loginRules"
      class="login-form"
      autocomplete="on"
      label-position="left"
    >
      <div class="title-container">
        <img src="../../assets/police-badge.svg" class="logo-img" alt="警徽" />
        <h3 class="title">公安视频采集平台</h3>
      </div>

      <el-form-item prop="username">
        <el-input
          v-model="loginForm.username"
          placeholder="用户名"
          type="text"
          tabindex="1"
          :prefix-icon="User"
        />
      </el-form-item>

      <el-form-item prop="password">
        <el-input
          v-model="loginForm.password"
          placeholder="密码"
          :type="passwordVisible ? 'text' : 'password'"
          tabindex="2"
          :prefix-icon="Lock"
          :suffix-icon="passwordVisible ? View : Hide"
          @click-suffix="passwordVisible = !passwordVisible"
        />
      </el-form-item>

      <el-button
        :loading="loading"
        type="primary"
        style="width: 100%; margin-bottom: 30px"
        @click="handleLogin"
      >
        登录
      </el-button>
    </el-form>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from "vue";
import { useRouter } from "vue-router";
import { ElMessage } from "element-plus";
import { User, Lock, View, Hide } from "@element-plus/icons-vue";

const router = useRouter();
const loading = ref(false);
const passwordVisible = ref(false);

const loginForm = reactive({
  username: "",
  password: "",
});

const loginRules = {
  username: [{ required: true, message: "请输入用户名", trigger: "blur" }],
  password: [{ required: true, message: "请输入密码", trigger: "blur" }],
};

const handleLogin = () => {
  loading.value = true;
  // 这里模拟登录，实际项目中需要调用后端接口
  setTimeout(() => {
    if (loginForm.username === "admin" && loginForm.password === "admin") {
      // 先清除之前的token
      localStorage.removeItem("token");
      // 设置新的token
      localStorage.setItem("token", "mock-token");
      router.push("/");
      ElMessage.success("登录成功");
    } else {
      ElMessage.error("用户名或密码错误");
    }
    loading.value = false;
  }, 1000);
};
</script>

<style scoped>
.login-container {
  min-height: 100vh;
  width: 100%;
  background: linear-gradient(to bottom, #001529, #003366);
  overflow: hidden;
  display: flex;
  align-items: center;
  justify-content: center;
}

.login-form {
  position: relative;
  width: 420px;
  max-width: 100%;
  padding: 35px;
  margin: 0 auto;
  overflow: hidden;
  background: rgba(255, 255, 255, 0.9);
  border-radius: 8px;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
}

.title-container {
  position: relative;
  text-align: center;
  margin-bottom: 40px;
}

.logo-img {
  width: 64px;
  height: 64px;
  margin-bottom: 20px;
  animation: pulse 2s infinite;
}

.title {
  font-size: 26px;
  color: #001529;
  margin: 0;
  font-weight: bold;
  background: linear-gradient(45deg, #001529, #409eff);
  -webkit-background-clip: text;
  background-clip: text;
  color: transparent;
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

:deep(.el-input__wrapper) {
  background-color: transparent;
}

:deep(.el-input__inner) {
  caret-color: #409eff;
}
</style>
