<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import AssistantChat from '../src/components/AssistantChat.vue'
import AbilitiesPanel from '../src/components/AbilitiesPanel.vue'

const chat = ref<InstanceType<typeof AssistantChat> | null>(null)
const applyExample = (t: string) => { chat.value?.useExample(t) }

// 根据浏览器/客户端窗口高度设置 CSS 变量，用于更精准的视口高度
function setAppHeight() {
  const h = window.innerHeight || document.documentElement.clientHeight
  document.documentElement.style.setProperty('--app-height', `${h}px`)
}

onMounted(() => {
  setAppHeight()
  window.addEventListener('resize', setAppHeight)
})

onUnmounted(() => {
  window.removeEventListener('resize', setAppHeight)
})
</script>

<template>
  <div class="app-root">
    
    <main class="app-main">
      <AbilitiesPanel @use-example="applyExample" />
      <AssistantChat ref="chat" />
    </main>
  </div>
</template>

<style scoped>
.app-root {
  /* 页面容器占满窗口，并提供内边距 */
  display: flex;
  flex-direction: column;
  gap: 20px;
  height: calc(var(--app-height, 100vh) - 26px); /* 基于客户端高度，默认回退到 100vh */
  box-sizing: border-box;
  max-width: none;
  margin: 0;
  /* padding: clamp(16px, 3vw, 32px); */
  text-align: left;
  /* 左栏固定宽度，右栏自适应 */
  --panel-width: clamp(280px, 22vw, 360px);
}
.app-main {
  display: grid;
  grid-template-columns: var(--panel-width) minmax(0, 1fr);
  gap: 20px;
  align-items: stretch; /* 让右侧内容拉伸填满 */
  height: 100%;
  min-height: 0; /* 允许子项内部滚动 */
}
@media (max-width: 960px) {
  .app-main {
    grid-template-columns: 1fr;
    grid-template-rows: auto 1fr; /* 左栏在上，右栏铺满下方 */
  }
}
</style>
