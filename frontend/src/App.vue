<template>
  <el-container class="app-container">
    <el-aside class="app-aside">
      <AbilitiesPanel @use-example="applyExample" />
    </el-aside>
    <el-main class="app-main">
      <AssistantChat ref="chat" />
    </el-main>
  </el-container>
</template>

<script setup lang="ts">
import { onMounted, onUnmounted, ref } from 'vue'
import AssistantChat from './components/AssistantChat.vue'
import AbilitiesPanel from './components/AbilitiesPanel.vue'

const chat = ref<InstanceType<typeof AssistantChat> | null>(null)

function setAppHeight() {
  const vh = window.innerHeight * 0.01
  document.documentElement.style.setProperty('--app-vh', `${vh}px`)
  document.documentElement.style.setProperty('--app-height', `${vh * 100}px`)
}

onMounted(() => {
  setAppHeight()
  window.addEventListener('resize', setAppHeight)
})

onUnmounted(() => {
  window.removeEventListener('resize', setAppHeight)
})

function applyExample(text: string) {
  chat.value?.setInputText?.(text)
}
</script>

<style scoped>
:root {
  --panel-width: clamp(278px, 27vw, 360px);
}

.app-container {
  height: calc(var(--app-height, 100vh));
  background: var(--bg-app, #1e1e1e);
}

.app-aside {
  /* padding: 12px; */
  height: 100%;
  /* border-right: 1px solid var(--el-border-color); */
  background: var(--el-color-black);
}

.app-main {
  height: 100%;
  /* padding: 12px; */
}
.el-aside{
  overflow: hidden;
}
.el-main{
  --el-main-padding: 0;
}

</style>
