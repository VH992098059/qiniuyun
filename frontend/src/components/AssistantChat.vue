<script setup lang="ts">
import { ElButton, ElMessage} from 'element-plus';
import { ref, nextTick, onMounted, onUnmounted, computed } from 'vue'
import { transcribeBlob, pickText } from '../services/asr'
import { voiceService } from '../services/voice'

interface Msg { role: 'user' | 'assistant'; text: string; time: number }
const messages = ref<Msg[]>([{
  role: 'assistant',
  text: '你好，我是你的语音桌面助理。请按麦克风开始说话，或直接输入文本。',
  time: Date.now()
}])

const inputText = ref('')
const isRecording = ref(false)
const isTranscribing = ref(false)
const asrAvailable = !!(navigator.mediaDevices && 'MediaRecorder' in window)
let mediaRecorder: MediaRecorder | null = null
let recordedChunks: BlobPart[] = []
let mediaStream: MediaStream | null = null
const chatBox = ref<HTMLDivElement | null>(null)

function scrollToBottom() {
  nextTick(() => {
    const el = chatBox.value
    if (el) el.scrollTop = el.scrollHeight
  })
}

async function speak(text: string) {
  const lastMsg = messages.value[messages.value.length - 1]
  const msgId = String(lastMsg?.time ?? Date.now())
  try {
    await voiceService.readElMessage(msgId, text)
  } catch {
    // 回退：浏览器语音合成
    try {
      const u = new SpeechSynthesisUtterance(text)
      u.lang = 'zh-CN'
      speechSynthesis.speak(u)
    } catch {}
  }
}

async function sendText(text: string) {
  const trimmed = text.trim()
  if (!trimmed) return
  messages.value.push({ role: 'user', text: trimmed, time: Date.now() })
  inputText.value = ''
  scrollToBottom()
  // 调用 Electron 主进程进行模拟执行
  const api = (window as any).electronAPI
  if (api && typeof api.executeCommand === 'function') {
    const res = await api.executeCommand(trimmed)
    const reply = res?.reply ?? '（主进程未响应，前端占位文本）'
    messages.value.push({ role: 'assistant', text: reply, time: Date.now() })
    await speak(reply)
  } else {
    const reply = '（Electron 未连接，前端占位响应：已接收你的指令）'
    messages.value.push({ role: 'assistant', text: reply, time: Date.now() })
    await speak(reply)
  }
  scrollToBottom()
}

function useExample(text: string) {
  inputText.value = text
}

async function startRecording() {
  if (!asrAvailable) {
    ElMessage.warning('当前环境不支持录音')
    return
  }
  try {
    mediaStream = await navigator.mediaDevices.getUserMedia({ audio: true })
    recordedChunks = []
    mediaRecorder = new MediaRecorder(mediaStream, { mimeType: 'audio/webm' })
    mediaRecorder.ondataavailable = (e: BlobEvent) => {
      if (e.data && e.data.size > 0) recordedChunks.push(e.data)
    }
    mediaRecorder.onstart = () => { isRecording.value = true }
    mediaRecorder.onstop = async () => {
      isRecording.value = false
      const blob = new Blob(recordedChunks, { type: 'audio/webm' })
      recordedChunks = []
      // 停止所有轨道
      mediaStream?.getTracks().forEach(t => t.stop())
      mediaStream = null
      try {
        isTranscribing.value = true
        const resp = await transcribeBlob(blob, 'auto')
        const text = pickText(resp).trim()
        if (text) {
          await sendText(text)
        } else {
          ElMessage.info('未识别到有效语音')
        }
      } catch (err) {
        console.error('ASR 识别失败', err)
        ElMessage.error('语音识别失败')
      } finally {
        isTranscribing.value = false
      }
    }
    mediaRecorder.start()
  } catch (err) {
    console.error('获取麦克风失败', err)
    ElMessage.error('无法访问麦克风')
  }
}

function stopRecording() {
  if (mediaRecorder && mediaRecorder.state !== 'inactive') {
    try { mediaRecorder.stop() } catch {}
  }
  isRecording.value = false
}

onMounted(() => {
  scrollToBottom()
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
  voiceService.destroy()
  if (mediaRecorder && mediaRecorder.state !== 'inactive') {
    try { mediaRecorder.stop() } catch {}
  }
  mediaStream?.getTracks().forEach(t => t.stop())
})

// 暴露给父组件使用示例填充
defineExpose({ useExample })

const srLabel = computed(() => {
  if (isRecording.value) return '录音中…'
  if (isTranscribing.value) return '识别中…'
  return asrAvailable ? '语音识别 API 已接入' : '语音识别不可用'
})
</script>

<template>
  <section class="chat">
    <div class="status">
      <span class="dot" :class="{ recording: isRecording }"></span>
      {{ srLabel }}
    </div>
    <div class="chat-box" ref="chatBox">
      <div v-for="(m, i) in messages" :key="i" class="msg" :class="m.role">
        <div class="bubble">
          <p class="text">{{ m.text }}</p>
          <span class="time">{{ new Date(m.time).toLocaleTimeString() }}</span>
        </div>
      </div>
    </div>
    <div class="toolbar">
      <ElButton class="mic" @click="isRecording ? stopRecording() : startRecording()">
        {{ isRecording ? '⏹️ 停止' : '🎤 开始录音' }}
      </ElButton>
      <input v-model="inputText" placeholder="在此输入或按麦克风说话" @keyup.enter="sendText(inputText)" />
      <ElButton class="send" @click="sendText(inputText)">📨 发送</ElButton>
    </div>
  </section>
</template>

<style scoped>
.chat {
  border: 1px solid #3a3a3a;
  border-radius: 12px;
  padding: 12px;
  background: rgba(255,255,255,0.04);
  display: flex; /* 让内部按列布局以充满高度 */
  flex-direction: column;
  height: 100%;
  min-height: 0; /* 允许内部滚动 */
}
.status {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 13px;
  color: #7f8ea3;
  margin-bottom: 8px;
}
.dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background: #7f8ea3;
}
.dot.recording {
  background: #e53935;
  box-shadow: 0 0 0 0 rgba(229,57,53, 0.7);
  animation: pulse 1.5s infinite;
}
@keyframes pulse {
  0% { box-shadow: 0 0 0 0 rgba(229,57,53, 0.7); }
  70% { box-shadow: 0 0 0 10px rgba(229,57,53, 0); }
  100% { box-shadow: 0 0 0 0 rgba(229,57,53, 0); }
}
.chat-box {
  /* 改为弹性填充高度 */
  flex: 1;
  min-height: 0;
  overflow: auto;
  display: flex;
  flex-direction: column;
  gap: 10px;
  background: linear-gradient(180deg, rgba(0,0,0,0.05), rgba(0,0,0,0.08));
  border-radius: 8px;
  padding: 12px;
}
.bubble {
  max-width: 80%;
  padding: 10px 12px;
  border-radius: 14px;
  color: #fff;
  box-shadow: 0 6px 16px rgba(0,0,0,0.15);
  white-space: pre-wrap;
  word-break: break-word;
}
.msg.assistant .bubble { background: linear-gradient(135deg, #1e293b, #334155); }
.msg.user .bubble { background: linear-gradient(135deg, #2563eb, #3b82f6); }
.text { margin: 0 0 6px 0; }
.time { font-size: 12px; opacity: 0.9; }
.toolbar {
  display: grid;
  grid-template-columns: 200px 1fr 140px;
  gap: 10px;
  margin-top: 12px;
}
.toolbar input {
  padding: 10px;
  border-radius: 10px;
  border: 1px solid #4a4a4a;
}
.mic {
  background: #3949ab;
  border-radius: 10px;
}
.send {
    background: #1e88e5;
    border-radius: 10px;
  }
  @media (max-width: 600px) {
    .toolbar {
      grid-template-columns: 1fr;
    }
  }
</style>