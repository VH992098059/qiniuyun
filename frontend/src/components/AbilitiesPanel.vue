<template>
  <el-card class="panel-card" shadow="hover">
    <template #header>
      <div class="panel-header">能力面板</div>
    </template>

    <el-scrollbar class="panel-body">
      <div class="group">
        <div class="group-title">常用操作</div>
        <el-space wrap size="8">
          <el-button type="primary" plain @click="emitUse('打开记事本')">打开记事本</el-button>
          <el-button type="primary" plain @click="emitUse('打开控制面板')">打开控制面板</el-button>
          <el-button type="primary" plain @click="emitUse('打开资源管理器')">打开资源管理器</el-button>
          <el-button type="primary" plain @click="emitUse('关闭窗口')">关闭窗口</el-button>
        </el-space>
      </div>

      <div class="group">
        <div class="group-title">系统工具</div>
        <el-space wrap size="8">
          <el-button plain @click="emitUse('截图保存到桌面')">截图保存到桌面</el-button>
          <el-button plain @click="emitUse('查询系统信息')">查询系统信息</el-button>
          <el-button plain @click="emitUse('打开任务管理器')">打开任务管理器</el-button>
        </el-space>
      </div>

      <div class="group">
        <div class="group-title">示例</div>
        <el-space direction="vertical" alignment="stretch" size="6" class="examples">
          <el-button v-for="e in examples" :key="e" type="success" plain size="small" @click="emitUse(e)" style="width: 255px;">{{ e }}</el-button>
        </el-space>
      </div>
    </el-scrollbar>
  </el-card>
</template>

<script setup lang="ts">
import { ref } from 'vue'

const emit = defineEmits<{
  (e: 'use-example', text: string): void
}>()

const examples = ref<string[]>([
  '查询今天的待办并朗读',
  '打开微信并置顶窗口',
  '搜索并打开 Chrome 浏览器',
  '打开下载目录并显示最新文件',
])

function emitUse(t: string) {
  emit('use-example', t)
}
</script>

<style scoped>
.panel-card { height: 100%; display: flex; flex-direction: column; }
.panel-header { font-weight: 600; }
.panel-body { flex: 1 1 auto; min-height: 0; }
.group { margin-bottom: 12px; }
.group-title { color: var(--el-text-color-regular); margin-bottom: 6px; font-size: 13px; }
.examples { margin-top: 6px; }
</style>