<script setup lang="ts">
import { defineEmits } from 'vue'

const emit = defineEmits<{
  (e: 'use-example', text: string): void
}>()

const examples = [
  '播放音乐',
  '写一篇关于人工智能的文章，300字',
  '打开记事本',
  '组合任务：搜索一首舒缓音乐并播放，然后新建文稿写开场白',
]
</script>

<template>
  <aside class="panel">
    <h2>能力面板</h2>
    <ul class="list">
      <li class="item"><span class="icon">🎵</span><span class="label">播放音乐</span></li>
      <li class="item"><span class="icon">📝</span><span class="label">写文章</span></li>
      <li class="item"><span class="icon">🗂️</span><span class="label">打开应用</span></li>
      <li class="item"><span class="icon">🧩</span><span class="label">组合任务</span></li>
    </ul>
    <h3>示例</h3>
    <div class="examples">
      <button v-for="ex in examples" :key="ex" @click="emit('use-example', ex)">{{ ex }}</button>
    </div>
  </aside>
</template>

<style scoped>
.panel {
  border: 1px solid #3a3a3a;
  border-radius: 12px;
  padding: 12px;
  background: rgba(255,255,255,0.04);
  height: 100%; /* 让左侧面板在栅格里拉伸填满 */
  display: flex;
  flex-direction: column;
}

.panel .header {
  font-size: 14px;
  color: #a0afc2;
  margin-bottom: 8px;
}

.panel .content {
  flex: 1;
  min-height: 0; /* 允许内部滚动 */
  overflow: auto;
}
.panel h2, .panel h3 {
  margin: 6px 0 10px 0;
  color: #c9d7e5;
  font-weight: 600;
}
.list {
  list-style: none;
  padding: 0;
  margin: 0 0 12px 0;
  display: flex;
  flex-direction: column;
  gap: 8px;
}
.item {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 10px 12px;
  border-radius: 12px;
  border: 1px solid rgba(255,255,255,0.1);
  background: #1f2937; /* slate-800 */
  transition: border-color .2s ease, background .2s ease, transform .15s ease;
}
.item:hover {
  border-color: #4b5563; /* slate-600 */
  background: #233143;
  transform: translateY(-1px);
}
.icon { font-size: 18px; }
.label { font-size: 14px; color: #e2e8f0; }
.examples {
  display: grid;
  grid-template-columns: 1fr;
  gap: 8px;
}
.examples button {
  text-align: left;
  border-radius: 10px;
  border: 1px solid rgba(255,255,255,0.1);
  padding: 8px 12px;
  background: #263246;
  color: #e3e9f2;
  transition: border-color .2s ease, background .2s ease, transform .15s ease;
}
.examples button:hover {
  border-color: #4b5563;
  background: #2d3b52;
  transform: translateY(-1px);
}
.examples button:active {
  transform: translateY(0) scale(0.99);
}
</style>