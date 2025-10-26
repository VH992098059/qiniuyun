# 后端 Golang 架构设计文档

本文档描述 `backend/golang` 的整体架构、模块职责与关键流程，便于开发、运维与扩展。

- 运行平台：Windows（当前桌面自动化实现针对 Windows 优化）
- Web 框架：GoFrame v2（`github.com/gogf/gf/v2`）
- 路由前缀：`/client`（HTTP），`/ws`（WebSocket）
- 外部依赖：ASR 微服务、OCR 微服务、TTS 服务（CosyVoice）、桌面自动化（robotgo）

---

## 架构总览

- 应用入口
  - `internal/cmd/cmd.go`：注册并启动 HTTP 服务，启用 CORS，绑定路由分组：
    - `GET /ws`：WebSocket 通道，用于双向消息传输（示例 Echo）。
    - `Group /client`：绑定 `voice` 控制器（v1）。
- API 层
  - 定义：`api/voice/v1/voice.go`（请求/响应结构、路径与方法）
  - 控制器：`internal/controller/voice/*.go` 实现 `/tts`、`/asr`、`/text`。
- 业务逻辑层（Logic）
  - 语音合成（TTS）：`internal/logic/voice.go`，调用上游 TTS 服务返回 MP3 音频流。
  - 语音识别（ASR）：`internal/logic/asr.go`，将 Base64 音频转为 `multipart/form-data` 上传至本地 ASR 微服务，解析识别文本、触发桌面自动化，并将识别文本再做 TTS 返回音频。
  - 桌面自动化：`internal/logic/robot.go` 结合模型意图、OCR 标注与 `robotgo` 执行动作（点击/输入/快捷键）。
  - 应用启动：`internal/logic/app.go` 在 Windows 上通过注册表或 PATH 查找并启动目标应用。
  - 模型服务：`internal/logic/model_service.go` 调用本地模型服务（`http://localhost:8081/api/ai/chat/do`）获取意图/计划/动作参数。
- 工具层（Utility）
  - OCR 客户端：`utility/ocr.go` 将桌面截图发送至 Flask OCR 微服务。
  - JSON 解析与坐标计算：`utility/json_search.go`、`utility/new_json.go`、`utility/synthesis_photo.go` 对 OCR 输出进行处理并标注坐标。
  - 流式响应工具：`utility/stream.go`（SSE 工具，当前未在公开路由中使用）。
- 外部微服务（位于 `files/` 或独立进程）
  - OCR（Flask + PaddleOCR）：`files/OCR_api.py`，`POST /ocr` 接收 `image` 并返回 JSON 结果。
  - ASR（FastAPI/自定义）：文档内假定地址 `http://localhost:50000/api/v1/asr`。
  - TTS（CosyVoice2-0.5B）：按配置 `chat.baseURL`、`chat.apiKey` 调用。

---

## 模块与责任

- `internal/cmd/cmd.go`
  - 初始化并启动 HTTP 服务，配置 CORS 与路由分组。
  - 绑定 `voice.NewV1()` 控制器到 `/client`。
- `internal/controller/voice`
  - `voice_v1_voice.go`（`POST /client/tts`）：接收文本，调用 TTS，直接写回 `audio/mpeg`。
  - `voice_v1_asr.go`（`POST /client/asr`）：接收 Base64 音频，调用 ASR→桌面自动化→TTS，返回 `audio/mpeg`。
  - `voice_v1_text.go`（`POST /client/text`）：交给模型分析并触发桌面自动化，返回 `{res_text:"已完成"}`。
- `api/voice/v1/voice.go`
  - 请求与响应结构体定义，路径元信息：`/tts`、`/asr`、`/text`。
- `internal/logic`
  - `voice.go`：TTS 请求构造（`/audio/speech`），设置 `Authorization: Bearer <apiKey>` 与 `Accept: audio/mpeg`，直接读取音频流返回。
  - `asr.go`：
    - 支持纯 Base64 或 Data URI（`data:audio/wav;base64,...` / `data:audio/flac;base64,...`）。
    - 仅允许 `wav/flac`（避免上游不兼容）。
    - 上传字段名为 `files`（`multipart/form-data`）。
    - 解析响应中 `result`（对象或数组）与顶层 `clean_text/text/raw_text`。
    - 触发 `RobotAutoApplication` 并对识别文本做 TTS。
  - `robot.go`：
    - 调用 `ModelService` 获取 `SafetyCheck/Intention` 与动作计划。
    - 启动应用、截图、调用 OCR、解析动作列表并执行（click/typestr/keytap/finish）。
  - `app.go`：
    - Windows 下通过注册表 `Uninstall` 项查找安装路径；找不到则回退 `PATH`。
    - 调用 `exec.Command(...)` 启动目标程序。
  - `model_service.go`：
    - 文本意图服务：`POST http://localhost:8081/api/ai/chat/do`，外层结构 `modelChatRes{code,message,data,timestamp}`，其中 `data` 为 JSON 字符串，内层 `DataAnalysis`。
    - 动作分析服务：多部分表单上传图片与文本，解析 `overall_goal/goal_analysis/status/thought/actions`。
- `utility/ocr.go`
  - `POST http://localhost:5000/ocr` 字段 `image` 上传截图，返回 JSON 数组，后续取首元素并写文件以供坐标标注。

---

## 关键流程

- `/client/text`
  1) `TextReq{text}` → `RobotAutoApplication(ctx,text)`。
  2) 调用 `ModelService` 获取意图（可能为 `action`）。
  3) 启动应用并轮询激活 → 截图保存 `files/photos/screenshot.png`。
  4) 发送截图至 OCR → 解析坐标 → 生成标注图。
  5) 根据动作列表执行：鼠标移动与点击、键入文本、按下快捷键、完成收尾。
  6) 返回 `{res_text:"已完成"}`。
- `/client/asr`
  1) 校验与解码 Base64 音频（仅 `wav/flac`）。
  2) 上传至 ASR 微服务 → 解析识别文本。
  3) 调用 `RobotAutoApplication` 执行桌面动作。
  4) 对识别文本调用 TTS → 返回 `audio/mpeg`。
- `/client/tts`
  1) 直接将 `input` 文本请求上游 TTS → 返回 `audio/mpeg`。
- `/ws`
  - 使用 `gorilla/websocket` 实现简单 Echo，当前用于演示或后续扩展。

---

## 配置与环境

- 端口与路径
  - 默认地址：`:8000`（见 `manifest/deploy/kustomize/overlays/develop/configmap.yaml`）。
  - 忽略的本地配置路径：`.gitignore` 中屏蔽 `**/config/config.yaml`。
- TTS 配置
  - 从 `g.Cfg()` 读取：`chat.baseURL` 与 `chat.apiKey`。
  - 请求 `POST <baseURL>/audio/speech`，返回 MP3。
- ASR 微服务
  - 目标地址：`http://localhost:50000/api/v1/asr`。
  - 接口要求：`multipart/form-data`，字段 `files` 为音频文件，`language=auto`。
- OCR 微服务
  - 目标地址：`http://localhost:5000/ocr`，字段 `image`。
  - 服务实现位于 `files/OCR_api.py`（Flask + PaddleOCR）。

---

## 错误与返回规范

- 控制器统一通过 GoFrame 响应中间件处理错误；常见：
  - `400`：参数缺失或格式错误（例如 ASR 音频格式不支持）。
  - `502`：上游服务不可用或返回异常（TTS/ASR/模型/OCR）。
  - `500`：服务内部错误。
- 音频响应统一设置 `Content-Type: audio/mpeg` 并直接写二进制。

---

## 安全与权限

- CORS：默认允许跨域（`cmd.go` 中 `r.Response.CORSDefault()`）。
- 桌面自动化：需在本机运行并具有必要权限；请谨慎使用，并确保只执行受控动作。
- 密钥：`chat.apiKey` 不应提交到版本库；建议通过环境变量或安全配置管理。

---

## 日志与监控

- 运行日志：使用标准库 `log` 输出关键过程与错误。
- GoFrame 日志：可通过配置调整 `level/stdout`（参考 `manifest/.../configmap.yaml`）。
- 建议：后续接入结构化日志与错误追踪（OpenTelemetry 已在 `go.mod` 依赖中）。

---

## 扩展规划

- 增加仅文本返回的 ASR 端点（例如 `/client/asr/text`）。
- WebSocket 管道用于实时状态上报与指令流式控制。
- SSE 工具与流式模型输出集成（`utility/stream.go`）。
- 增强跨平台兼容（macOS/Linux 的应用启动与窗口管理）。
- 性能优化：预热模型、并行 OCR/意图分析、端到端流式推理，减少处理等待时间。

---

## 参考路径索引

- 入口与路由：`internal/cmd/cmd.go`
- 控制器：`internal/controller/voice/voice_v1_*.go`
- API 定义：`api/voice/v1/voice.go`
- 逻辑：`internal/logic/*.go`
- 工具：`utility/*.go`
- OCR 微服务：`files/OCR_api.py`