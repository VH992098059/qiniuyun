# Golang接口文档

本文件描述 `backend/golang/api/voice/v1` 提供的接口。路由前缀为 `/client`，默认服务地址为 `http://localhost:8000`（可通过配置调整）。

- 版本：v1
- 前缀：`/client`
- 鉴权：无（默认开启 CORS）
- 通用请求头：`Content-Type: application/json`
- 音频响应头：`Content-Type: audio/mpeg`

## 依赖与前置条件

- TTS：使用 `CosyVoice2-0.5B`，需在配置中设置 `chat.baseURL` 与 `chat.apiKey`。
- ASR：调用本地 ASR 微服务 `http://localhost:50000/api/v1/asr`（服务需先启动）。
- 桌面自动化：`/text` 与 `/asr` 会触发桌面自动化流程（应用启动、点击、输入等），依赖 OCR 服务与 `robotgo`。

---

## 1) 文本转语音（TTS）

- 路径：`POST /client/tts`
- 描述：将输入文本合成为 MP3 音频并直接返回音频字节流。
- 请求体：
  ```json
  {
    "input": "string" // 必填，要合成的文本
  }
  ```
- 响应：
  - 200：`audio/mpeg`（MP3 音频二进制）
  - 4xx/5xx：JSON 错误（由 GoFrame 中间件统一返回）
- 示例：
  ```bash
  curl -X POST \
    http://localhost:8000/client/tts \
    -H "Content-Type: application/json" \
    -d '{"input":"你好，世界"}' \
    --output tts.mp3
  ```
- 可能错误：
  - 400：缺少 `input`
  - 502：上游 TTS 服务错误或配置缺失（`chat.baseURL`/`chat.apiKey`）
  - 500：其他服务内部错误

---

## 2) 语音识别并朗读结果（ASR→TTS）

- 路径：`POST /client/asr`
- 描述：接收 Base64 编码的音频，识别文本后：

  1) 启动桌面自动化（按识别的意图控制应用）；
  2) 将识别结果再做 TTS，返回 MP3 音频。
- 请求体：

  ```json
  {
    "audio_base64": "string", // 必填，Base64 音频；支持两种形式
    "language": "auto"         // 选填，默认 auto
  }
  ```

  - 支持的 `audio_base64` 形式：
    - 纯 Base64（WAV/FLAC）
    - Data URI：例如 `data:audio/wav;base64,<BASE64>` 或 `data:audio/flac;base64,<BASE64>`
  - 不支持：`webm/opus` 等格式（会返回 415 类错误）。
- 响应：

  - 200：`audio/mpeg`（MP3 音频二进制，为识别文本的合成语音）
  - 4xx/5xx：JSON 错误
- 示例：

  ```bash
  curl -X POST \
    http://localhost:8000/client/asr \
    -H "Content-Type: application/json" \
    -d '{
      "audio_base64": "data:audio/wav;base64,PD94bW...",
      "language": "auto"
    }' \
    --output asr_result.mp3
  ```
- 处理流程概要（内部实现）：

  1) 解码 Base64，校验格式（仅 WAV/FLAC）。
  2) 转为 `multipart/form-data` 上传至 `http://localhost:50000/api/v1/asr`（字段名 `files`）。
  3) 解析响应，优先取 `clean_text -> text -> raw_text`。
  4) 调用桌面自动化（应用启动/点击/输入等）。
  5) 将识别文本进行 TTS，返回音频。
- 可能错误：

  - 400：缺少 `audio_base64`
  - 415：不支持的音频格式（请使用 `audio/wav` 或 `audio/flac`）
  - 502：上游 ASR 服务错误或未启动
  - 500：其他服务内部错误

---

## 3) 文本指令入口（触发桌面自动化）

- 路径：`POST /client/text`
- 描述：将文本交给模型分析并触发桌面自动化（应用启动、OCR 分析、点击、键盘输入等）。
- 请求体：
  ```json
  {
    "text": "string" // 选填（允许空），为空时会按默认策略处理
  }
  ```
- 响应：
  ```json
  { "res_text": "已完成" }
  ```
- 处理流程概要（内部实现）：
  1) 调用 `ModelService` 获取意图与任务计划；
  2) 启动应用并截图，走 OCR 标注与动作解析；
  3) 通过 `robotgo` 执行动作（点击、输入、快捷键、完成）。
- 可能错误：
  - 502：模型或 OCR/自动化依赖未就绪
  - 500：其他服务内部错误

---

## 数据结构（定义于 `api/voice/v1/voice.go`）

- `VoiceReq`: `{ input: string }`
- `VoiceRes`: 音频响应（`mime:"audio/mpeg"`）
- `AsrReq`: `{ audio_base64: string, language?: string="auto" }`
- `AsrRes`: 音频响应（`mime:"audio/mpeg"`）
- `TextReq`: `{ text?: string }`
- `TextRes`: `{ res_text: string }`

---

## 常见问题

- 为什么 `/client/asr` 返回的是音频而不是文本？
  - 该端点在识别文本后会触发桌面自动化，并同时将识别结果做 TTS 返回音频，便于直接播放提醒/反馈。
- 如何仅获取识别文本？
  - 当前 v1 接口未提供仅文本返回的 ASR 端点。如需此能力，可新增一个返回 `application/json` 的接口（例如 `/client/asr/text`）。
- TTS 不返回音频或报 502？
  - 检查配置 `chat.baseURL`、`chat.apiKey` 是否正确，确保上游 TTS 服务可用。
- ASR 报格式错误？
  - 请确保音频为 `audio/wav` 或 `audio/flac`，并以 Base64（或 Data URI）形式传入。
