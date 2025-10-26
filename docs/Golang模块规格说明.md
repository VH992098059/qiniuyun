# 后端 Golang 模块规格说明（backend/golang）

## 文档目标
- 说明 `backend/golang` 的模块边界、职责、接口与配置，便于协作与扩展。
- 明确与前端、Python OCR、Java LLM 的通信关系与依赖。

## 范围与定位
- 运行平台：Windows 优化（桌面自动化与应用启动能力）。
- 框架与工具：GoFrame v2、robotgo、gorilla/websocket。
- 职责聚焦：ASR/TTS/系统操作编排与 HTTP API 暴露。

## 模块总览
- 入口与路由：`internal/cmd/cmd.go`（启用 CORS，注册 `/ws` 与 `/client` 路由）
- 控制器层：`internal/controller/voice/*.go`（绑定 `api/voice/v1` 路由）
  - `voice_v1_voice.go` → `POST /client/tts`（文本转语音，返回 `audio/mpeg`）
  - `voice_v1_asr.go` → `POST /client/asr`（音频识别→自动化→TTS，返回 `audio/mpeg`）
  - `voice_v1_text.go` → `POST /client/text`（文本驱动的自动化，返回 `{res_text:"已完成"}`）
- API 定义：`api/voice/v1/voice.go`（请求/响应结构与元信息）
- 业务逻辑层：`internal/logic/*.go`
  - `voice.go`：TTS 请求构造与音频流返回
  - `asr.go`：Base64 音频校验/解码→`multipart/form-data` 上传至 ASR 微服务→解析识别文本→触发自动化→TTS
  - `robot.go`：应用启动、截图、OCR、坐标标注、动作执行（click/typestr/keytap/finish）
  - `app.go`：Windows 应用启动（注册表/`PATH` 搜索）
  - `model_service.go`：与 Java LLM 服务交互，拉取意图与动作计划
- 工具层：`utility/*.go`
  - `ocr.go`：`POST http://localhost:5000/ocr` 上传截图，解析结果
  - `json_search.go`、`new_json.go`、`synthesis_photo.go`：JSON/坐标处理与标注
  - `stream.go`：SSE 工具（当前未对外暴露）

## 进程与通信（端到端）
- 流程：`前端 → Go → OCR → Java LLM → Go → 前端`
  - 前端上传音频/文本到 Go：`POST /client/asr | /client/text | /client/tts`
  - Go 调用 OCR：`POST http://localhost:5000/ocr`（字段 `image`）
  - Go 调用 Java LLM：`POST http://localhost:8081/api/ai/chat/do`（文本意图与动作计划）
  - Go 执行桌面动作并回传音频或文本到前端（HTTP 响应或后续轮询/事件）

## 接口规格（v1）
- `POST /client/tts`
  - Req: `{ input: string }`
  - Res: `audio/mpeg`（直接写二进制）
- `POST /client/asr`
  - Req: `{ audio_base64: string, language?: "auto" }`（支持 Data URI；仅 `wav/flac`）
  - Res: `audio/mpeg`（识别文本经 TTS 返回）
- `POST /client/text`
  - Req: `{ text?: string }`
  - Res: `{ res_text: "已完成" }`

## 配置与环境
- 端口与路由
  - 默认监听：`:8000`（参考 `manifest/deploy/kustomize/overlays/develop/configmap.yaml`）
  - 路由分组：`/client`（HTTP）、`/ws`（WebSocket）
- TTS（CosyVoice2）
  - 从 `g.Cfg()` 读取 `chat.baseURL` 与 `chat.apiKey`
  - 目标接口：`POST <baseURL>/audio/speech`，设置 `Authorization: Bearer <apiKey>` 与 `Accept: audio/mpeg`
- ASR 微服务（SenseVoice/FunASR 自建）
  - 地址：`http://localhost:50000/api/v1/asr`
  - 上传：`multipart/form-data` 字段名 `files`，`language=auto`
- OCR 微服务（Flask + PaddleOCR）
  - 地址：`http://localhost:5000/ocr`
  - 上传：`multipart/form-data` 字段名 `image`

## 错误与返回规范
- `400`：参数缺失或格式错误（例如不支持的音频格式）
- `502`：上游服务不可用或异常（TTS/ASR/模型/OCR）
- `500`：内部错误（自动化执行失败、解析异常等）
- 音频响应统一设置 `Content-Type: audio/mpeg`

## 并发与性能建议
- 预热 Java LLM 与 ASR/TTS 服务，减少首呼延迟
- OCR 与意图分析可并行；动作执行串行保证一致性
- 大音频分片与流式识别（待规划）；结果缓存与重试策略

## 依赖与安装
- Go：`go mod tidy` 或 `go mod download`
- Windows 桌面自动化：`robotgo`（已在 `go.mod`）
- 外部服务：启动 Python OCR（`files/OCR_api.py`）、ASR（SenseVoice `api.py`）、Java LLM（`backend/java`）

## 测试与验证
- 单元测试：逻辑层（TTS/ASR 解析/坐标标注）
- 集成测试：
  - 启动各服务后，验证 `/client/tts` 与 `/client/asr` 音频返回
  - 使用示例截图验证 OCR→坐标标注→动作执行链路
- 回归检查：配置读取（`g.Cfg()`）与跨域响应（CORS）

## 扩展与约束
- 不接入第三方 agent，工具编排由自研有限规则实现
- 可新增仅文本返回的 ASR 端点，例如 `POST /client/asr/text`
- 规划：WebSocket 状态上报、SSE 流式输出、跨平台支持

## 路径索引
- 入口：`internal/cmd/cmd.go`
- 控制器：`internal/controller/voice/voice_v1_*.go`、`internal/controller/voice/voice_new.go`
- API：`api/voice/v1/voice.go`
- 逻辑：`internal/logic/*.go`
- 工具：`utility/*.go`
- OCR 微服务：`files/OCR_api.py`