<template>
  <div class="tags-view">
    <el-scrollbar class="tags-view-wrapper">
      <router-link
        v-for="tag in visitedViews"
        :key="tag.path"
        :class="isActive(tag) ? 'active' : ''"
        class="tags-view-item"
        :to="{ path: tag.path, query: tag.query }"
        @contextmenu.prevent="openMenu($event, tag)"
      >
        <el-icon class="tag-icon" v-if="tag.meta?.icon">
          <component :is="tag.meta.icon" />
        </el-icon>
        {{ tag.meta?.title }}
        <el-icon
          class="close-icon"
          v-if="visitedViews.length > 1"
          @click.prevent.stop="closeSelectedTag(tag)"
        >
          <Close />
        </el-icon>
      </router-link>
    </el-scrollbar>

    <!-- 右键菜单 -->
    <ul
      v-show="visible"
      :style="{ left: left + 'px', top: top + 'px' }"
      class="contextmenu"
    >
      <li @click="refreshSelectedTag(selectedTag)">刷新页面</li>
      <li @click="closeSelectedTag(selectedTag)">关闭当前</li>
      <li @click="closeOthersTags(selectedTag)">关闭其他</li>
      <li @click="closeAllTags">关闭所有</li>
    </ul>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch, onMounted, onBeforeUnmount } from "vue";
import { useRoute, useRouter } from "vue-router";
import type { RouteLocationNormalized } from "vue-router";

const route = useRoute();
const router = useRouter();

// 访问过的页面
const visitedViews = ref<RouteLocationNormalized[]>([]);

// 右键菜单相关
const visible = ref(false);
const top = ref(0);
const left = ref(0);
const selectedTag = ref<RouteLocationNormalized | null>(null);

// 判断是否是激活的标签
const isActive = (tag: RouteLocationNormalized) => {
  return tag.path === route.path;
};

// 添加标签
const addVisitedView = (view: RouteLocationNormalized) => {
  if (visitedViews.value.some((v) => v.path === view.path)) return;
  if (view.meta?.title) {
    visitedViews.value.push(Object.assign({}, view));
  }
};

// 关闭标签
const closeSelectedTag = (view: RouteLocationNormalized) => {
  const index = visitedViews.value.findIndex((v) => v.path === view.path);
  if (index > -1) {
    visitedViews.value.splice(index, 1);
  }
  if (isActive(view)) {
    toLastView();
  }
};

// 关闭其他标签
const closeOthersTags = (view: RouteLocationNormalized) => {
  visitedViews.value = visitedViews.value.filter((v) => v.path === view.path);
  if (!isActive(view)) {
    router.push(view);
  }
};

// 关闭所有标签
const closeAllTags = () => {
  visitedViews.value = [];
  router.push("/");
};

// 刷新页面
const refreshSelectedTag = (view: RouteLocationNormalized) => {
  router.replace({ path: "/redirect" + view.path });
};

// 跳转到最后一个标签
const toLastView = () => {
  const lastView = visitedViews.value.slice(-1)[0];
  if (lastView) {
    router.push(lastView);
  } else {
    router.push("/");
  }
};

// 打开右键菜单
const openMenu = (e: MouseEvent, tag: RouteLocationNormalized) => {
  const menuMinWidth = 105;
  const offsetLeft = document.documentElement.clientWidth - e.clientX;
  const offsetTop = document.documentElement.clientHeight - e.clientY;
  visible.value = true;
  selectedTag.value = tag;

  if (offsetLeft <= menuMinWidth) {
    left.value = e.clientX - menuMinWidth;
  } else {
    left.value = e.clientX;
  }
  top.value = e.clientY;
};

// 关闭右键菜单
const closeMenu = () => {
  visible.value = false;
};

// 监听路由变化
watch(
  () => route.path,
  () => {
    addVisitedView(route);
  }
);

// 监听点击事件关闭右键菜单
onMounted(() => {
  document.addEventListener("click", closeMenu);
});

onBeforeUnmount(() => {
  document.removeEventListener("click", closeMenu);
});

// 初始化时添加当前路由
addVisitedView(route);
</script>

<style scoped>
.tags-view {
  height: 34px;
  width: 100%;
  background: #fff;
  border-bottom: 1px solid #d8dce5;
  box-shadow: 0 1px 3px 0 rgba(0, 0, 0, 0.12);
  position: relative;
}

.tags-view-wrapper :deep(.el-scrollbar__wrap) {
  height: 34px;
}

.tags-view-item {
  display: inline-flex;
  align-items: center;
  position: relative;
  height: 26px;
  line-height: 26px;
  border: 1px solid #d8dce5;
  color: #495060;
  background: #fff;
  padding: 0 8px;
  font-size: 12px;
  margin-left: 5px;
  margin-top: 4px;
  text-decoration: none;
  border-radius: 3px;
  transition: all 0.3s cubic-bezier(0.645, 0.045, 0.355, 1);
}

.tags-view-item:first-of-type {
  margin-left: 15px;
}

.tags-view-item .tag-icon {
  margin-right: 4px;
  width: 14px;
  height: 14px;
}

.tags-view-item .close-icon {
  width: 14px;
  height: 14px;
  margin-left: 4px;
  border-radius: 50%;
  text-align: center;
  transition: all 0.3s cubic-bezier(0.645, 0.045, 0.355, 1);
  transform-origin: 100% 50%;
}

.tags-view-item .close-icon:hover {
  background-color: #b4bccc;
  color: #fff;
}

.tags-view-item:hover {
  background-color: #f5f7fa;
}

.tags-view-item.active {
  background-color: #42b983;
  color: #fff;
  border-color: #42b983;
}

.tags-view-item.active::before {
  content: "";
  background: #fff;
  display: inline-block;
  width: 8px;
  height: 8px;
  border-radius: 50%;
  position: relative;
  margin-right: 4px;
}

.tags-view-item.active .close-icon:hover {
  background-color: #fff;
  color: #42b983;
}

.contextmenu {
  margin: 0;
  background: #fff;
  z-index: 3000;
  position: absolute;
  list-style-type: none;
  padding: 5px 0;
  border-radius: 4px;
  font-size: 12px;
  font-weight: 400;
  color: #333;
  box-shadow: 2px 2px 3px 0 rgba(0, 0, 0, 0.3);
}

.contextmenu li {
  margin: 0;
  padding: 7px 16px;
  cursor: pointer;
}

.contextmenu li:hover {
  background: #eee;
}
</style>
