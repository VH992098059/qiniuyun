<script setup lang="ts">
import { ref, nextTick, onMounted, computed } from 'vue'

interface Msg { role: 'user' | 'assistant'; text: string; time: number }
const messages = ref<Msg[]>([{
  role: 'assistant',
  text: '你好，我是你的语音桌面助理。请按麦克风开始说话，或直接输入文本。',
  time: Date.now()
}])

const inputText = ref('')
const isRecording = ref(false)
const srAvailable = typeof (window as any).webkitSpeechRecognition !== 'undefined' || typeof (window as any).SpeechRecognition !== 'undefined'
let recognition: any = null
const chatBox = ref<HTMLDivElement | null>(null)

function scrollToBottom() {
  nextTick(() => {
    const el = chatBox.value
    if (el) el.scrollTop = el.scrollHeight
  })
}

function speak(text: string) {
  try {
    const u = new SpeechSynthesisUtterance(text)
    u.lang = 'zh-CN'
    speechSynthesis.speak(u)
  } catch {}
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
    speak(reply)
  } else {
    const reply = '（Electron 未连接，前端占位响应：已接收你的指令）'
    messages.value.push({ role: 'assistant', text: reply, time: Date.now() })
    speak(reply)
  }
  scrollToBottom()
}

function useExample(text: string) {
  inputText.value = text
}

function startRecording() {
  if (!srAvailable) {
    isRecording.value = true
    return
  }
  const SR = (window as any).webkitSpeechRecognition || (window as any).SpeechRecognition
  recognition = new SR()
  recognition.continuous = false
  recognition.interimResults = true
  recognition.lang = 'zh-CN'
  const finalChunks: string[] = []
  recognition.onresult = (event: any) => {
    for (let i = event.resultIndex; i < event.results.length; i++) {
      const res = event.results[i]
      const transcript = res[0].transcript
      if (res.isFinal) finalChunks.push(transcript)
    }
  }
  recognition.onstart = () => { isRecording.value = true }
  recognition.onerror = () => { isRecording.value = false }
  recognition.onend = async () => {
    isRecording.value = false
    const text = finalChunks.join(' ').trim()
    if (text) await sendText(text)
  }
  recognition.start()
}

function stopRecording() {
  if (recognition) try { recognition.stop() } catch {}
  isRecording.value = false
  // 若无 WebSpeech，模拟一次简单输入（可选）
}

onMounted(() => scrollToBottom())

// 暴露给父组件使用示例填充
defineExpose({ useExample })

const srLabel = computed(() => isRecording.value ? '录音中…' : (srAvailable ? '语音识别可用' : '语音识别不可用'))
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
      <button class="mic" @click="isRecording ? stopRecording() : startRecording()">
        {{ isRecording ? '⏹️ 停止' : '🎤 开始录音' }}
      </button>
      <input v-model="inputText" placeholder="在此输入或按麦克风说话" @keyup.enter="sendText(inputText)" />
      <button class="send" @click="sendText(inputText)">📨 发送</button>
    </div>
  </section>
</template>

<style scoped>
.chat {
  border: 1px solid #3a3a3a;
  border-radius: 12px;
  padding: 12px;
  background: rgba(255,255,255,0.04);
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
  height: 420px;
  overflow: auto;
  display: flex;
  flex-direction: column;
  gap: 10px;
  background: linear-gradient(180deg, rgba(0,0,0,0.05), rgba(0,0,0,0.08));
  border-radius: 8px;
  padding: 12px;
}
.bubble {
  max-width: 70%;
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
  grid-template-columns: 180px 1fr 120px;
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
</style>