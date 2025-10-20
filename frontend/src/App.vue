<script setup lang="ts">
import { ref } from 'vue'
import AssistantChat from '../src/components/AssistantChat.vue'
import AbilitiesPanel from '../src/components/AbilitiesPanel.vue'

const chat = ref<InstanceType<typeof AssistantChat> | null>(null)
const applyExample = (t: string) => { chat.value?.useExample(t) }
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
  display: flex;
  flex-direction: column;
  gap: 16px;
  max-width: 1200px;
  margin: 0 auto;
  padding: 20px;
}
.app-header {
  text-align: center;
}
.app-header h1 {
  margin: 0 0 6px 0;
  letter-spacing: 0.5px;
}
.sub {
  font-size: 14px;
  color: #8aa0b6;
}
.app-main {
  display: grid;
  grid-template-columns: 300px 1fr;
  gap: 16px;
  align-items: start;
}
@media (max-width: 960px) {
  .app-main {
    grid-template-columns: 1fr;
  }
}
</style>
