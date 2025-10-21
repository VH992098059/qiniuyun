<template>
  <el-card class="chat-card" shadow="hover">
    <template #header>
      <div class="card-header">
        <span>助手对话</span>
        <el-tag v-if="isRecording" type="danger" size="small">录音中</el-tag>
        <el-tag v-else-if="isTranscribing" type="warning" size="small">识别中</el-tag>
        <el-tag v-else-if="isSpeaking" type="success" size="small">朗读中</el-tag>
      </div>
    </template>

    <el-scrollbar class="chat-scroll" ref="scrollRef">
      <div
        v-for="(m, i) in messages"
        :key="i"
        class="msg"
        :class="m.role"
        style=""
      >
        <div class="bubble" >
          {{ m.text }}
        </div>
      </div>
    </el-scrollbar>

    <div class="toolbar">
      <el-button
        class="mic"
        type="primary"
        :loading="isTranscribing"
        :disabled="isTranscribing"
        @click="toggleRecording"
      >
        {{ srLabel }}
      </el-button>

      <el-input
        v-model="inputText"
        placeholder="输入指令，按回车或点击发送"
        clearable
        @keyup.enter="sendText(inputText)"
      />

      <el-button type="success" :disabled="!inputText" @click="sendText(inputText)">
        发送
      </el-button>
    </div>
  </el-card>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted, nextTick } from 'vue'
import { ElMessage } from 'element-plus'
import { transcribeBlob, pickText } from '../services/asr'
import { voiceService } from '../services/voice'

interface Message {
  role: 'user' | 'assistant'
  text: string
}

const messages = ref<Message[]>([])
const inputText = ref('')
const isRecording = ref(false)
const isTranscribing = ref(false)
const isSpeaking = ref(false)
const mediaRecorder = ref<MediaRecorder | null>(null)
let audioChunks: BlobPart[] = []
const scrollRef = ref()

// 供父组件设置输入框
function setInputText(t: string) {
  inputText.value = t
}

defineExpose({ setInputText })

const srLabel = computed(() => {
  if (isRecording.value) return '停止录音'
  if (isTranscribing.value) return '识别中…'
  return '开始录音'
})

function pushMessage(role: Message['role'], text: string) {
  messages.value.push({ role, text })
  nextTick(() => {
    // 自动滚动到底部
    try {
      const el = (scrollRef.value as any)?.wrapRef
      el && (el.scrollTop = el.scrollHeight)
    } catch {}
  })
}

async function sendText(text: string) {
  const t = text.trim()
  if (!t) return
  pushMessage('user', t)
  inputText.value = ''

  // 模拟后端响应，或在此处调用真实接口
  const reply = `收到：${t}`
  pushMessage('assistant', reply)
  await speak(reply)
}

async function speak(text: string) {
  isSpeaking.value = true
  try {
    await voiceService.readElMessage(String(Date.now()),text)
  } catch (err) {
    // 回退到浏览器 SpeechSynthesis
    try {
      const utter = new SpeechSynthesisUtterance(text)
      utter.lang = 'zh-CN'
      await new Promise((resolve) => {
        utter.onend = resolve as any
        speechSynthesis.speak(utter)
      })
    } catch (e) {
      console.error(e)
    }
  } finally {
    isSpeaking.value = false
  }
}

function toggleRecording() {
  if (isRecording.value) {
    stopRecording()
  } else {
    startRecording()
  }
}

async function startRecording() {
  if (isRecording.value) return
  try {
    const stream = await navigator.mediaDevices.getUserMedia({ audio: true })
    const rec = new MediaRecorder(stream)
    mediaRecorder.value = rec
    audioChunks = []

    rec.ondataavailable = (e) => {
      if (e.data.size > 0) audioChunks.push(e.data)
    }
    rec.onstop = async () => {
      const blob = new Blob(audioChunks, { type: 'audio/webm' })
      isTranscribing.value = true
      try {
        const raw = await transcribeBlob(blob)
        const text = pickText(raw)
        if (text) {
          pushMessage('user', text)
          await sendText(text)
        } else {
          ElMessage.info('未识别到有效语音内容')
        }
      } catch (e: any) {
        console.error(e)
        ElMessage.error(e?.message || '语音识别失败')
      } finally {
        isTranscribing.value = false
        // 释放资源
        stream.getTracks().forEach((t) => t.stop())
      }
    }

    rec.start()
    isRecording.value = true
  } catch (e: any) {
    console.error(e)
    ElMessage.error(e?.message || '无法开始录音')
    isRecording.value = false
  }
}

function stopRecording() {
  if (!isRecording.value) return
  try {
    mediaRecorder.value?.stop()
  } finally {
    isRecording.value = false
  }
}

onMounted(() => {
  voiceService.setCallbacks({
    onStateChange: () => {},
    onLoadStart: () => {},
    onCanPlay: () => {},
    onEnded: () => {},
    onError: () => {},
    onAbort: () => {},
  })
})

onUnmounted(() => {
  try { voiceService.destroy?.() } catch {}
  try { stopRecording() } catch {}
})
</script>

<style scoped>
.chat-card {
  display: flex;
  flex-direction: column;
}
.card-header {
  display: flex;
  align-items: center;
  gap: 8px;
}
.chat-scroll {
  flex: 1 1 auto;
  min-height: 0;
}
/* 提高选择器优先级，覆盖可能的固定高度 */
.el-scrollbar.chat-scroll {
  height: calc(100vh - 206px) !important;
  overflow-y:auto !important;
  padding: 30px !important; 
}
.toolbar {  
  display: grid;
  grid-template-columns: 120px 1fr 96px;
  gap: 8px;
  padding-top: 10px;
  border-top: 1px solid var(--el-border-color);
}
.msg {
  display: flex;
  margin: 6px 0;
}
.msg.user { justify-content: flex-end; }
.msg.assistant { justify-content: flex-start; }
.bubble {
  max-width: 72ch;
  padding: 10px 12px;
  border-radius: 10px;
  line-height: 1.6;
  background: var(--el-color-primary-light-9);
  color: var(--el-text-color-primary);
}
.msg.user .bubble { background: var(--el-color-success-light-9); }

</style>