import { createRouter, createWebHistory } from "vue-router";
import Layout from "../views/layout/index.vue";

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: "/login",
      name: "Login",
      component: () => import("../views/login/index.vue"),
      meta: { title: "登录" },
    },
    {
      path: "/",
      component: Layout,
      redirect: "/dashboard",
      children: [
        {
          path: "dashboard",
          name: "Dashboard",
          component: () => import("../views/dashboard/index.vue"),
          meta: { title: "仪表盘", requiresAuth: true },
        },
        {
          path: "device",
          name: "Device",
          component: () => import("../views/device/list.vue"),
          meta: { title: "设备管理", requiresAuth: true },
        },
        {
          path: "monitor/realtime",
          name: "Realtime",
          component: () => import("../views/monitor/realtime.vue"),
          meta: { title: "实时监控", requiresAuth: true },
        },
        {
          path: "monitor/history",
          name: "History",
          component: () => import("../views/monitor/history.vue"),
          meta: { title: "历史记录", requiresAuth: true },
        },
      ],
    },
  ],
});

// 路由守卫
router.beforeEach((to, from, next) => {
  const token = localStorage.getItem("token");
  const isLoginPage = to.path === "/login";

  if (!token && !isLoginPage) {
    // 没有token且不是登录页，重定向到登录页
    next({ path: "/login", query: { redirect: to.fullPath } });
  } else if (token && isLoginPage) {
    // 如果已登录且在登录页，继续访问登录页
    next();
  } else {
    next();
  }
});

export default router;
