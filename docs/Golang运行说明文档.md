# 后端 Golang 运行说明文档（backend/golang）

## 概述
- 项目位置：`backend/golang`
- 职责：本地自动化与系统集成（ASR/TTS/系统操作）、统一对外 API（`/client/*`）、与 OCR/LLM 编排。
- 运行平台：Windows 10/11（当前自动化实现面向 Windows 优化）。
- 框架：GoFrame v2。

## 环境要求
- Go 版本：建议 `Go 1.20+`。
- Windows 依赖：部分桌面自动化能力可能需要安装 MSYS2（详见根 `README.md`）。
- 可选外部服务：
  - OCR 微服务（Flask + PaddleOCR）：`files/OCR_api.py` 默认 `http://localhost:5000/ocr`
  - ASR 微服务（SenseVoice/FunASR，`api.py`）：默认 `http://localhost:50000/api/v1/asr`
  - Java LLM 服务（Spring AI）：默认 `http://localhost:8081/api`

## 依赖安装
```powershell
# 进入 Go 后端目录
cd W:\go_gf\Voice_Controlled_Desktop_Assistant\backend\golang

# 安装依赖（任选其一）
go mod tidy       # 推荐，清理并拉取
# 或
go mod download   # 仅下载依赖
```

## 本地配置（推荐）
在 `backend/golang/config/config.yaml`（该路径已在 `.gitignore` 中忽略）创建配置文件：

```yaml
server:
  address:     ":8000"      # 本地监听端口
  openapiPath: "/api.json"  # OpenAPI 文档（可选）
  swaggerPath: "/swagger"   # Swagger 文档（可选）

logger:
  level:  "all"
  stdout: true

# TTS 服务配置（CosyVoice2 等）
chat:
  baseURL: "http://localhost:30000"   # 例：你的 TTS 服务地址（需支持 POST /audio/speech）
  apiKey:  "YOUR_API_KEY"             # 你的 TTS 服务密钥
```

说明：
- `server.address` 未配置时，可能使用框架默认端口；建议明确设置为 `:8000`。
- `chat.baseURL` 与 `chat.apiKey` 为必需，`internal/logic/voice.go` 会从配置读取并调用 `POST <baseURL>/audio/speech`，返回 `audio/mpeg`。
- ASR 地址当前在代码中写死为 `http://localhost:50000/api/v1/asr`（见 `internal/logic/asr.go`）。如需修改，请相应调整代码或后续加入配置项。

## 启动方式
```powershell
# 方法一：直接运行
cd W:\go_gf\Voice_Controlled_Desktop_Assistant\backend\golang
go run .

# 方法二：构建并运行（Windows）
go build -o main.exe
./main.exe
```

启动后：
- 默认监听 `http://localhost:8000/`（若按上文配置）。
- 已启用 CORS（见 `internal/cmd/cmd.go`，`r.Response.CORSDefault()`）。
- 路由分组：`/client`（HTTP），`/ws`（WebSocket）。

## 接口校验（curl 示例）

- 1) 文本到语音（TTS）
```powershell
curl.exe -X POST "http://localhost:8000/client/tts" \
  -H "Content-Type: application/json" \
  -d "{\"input\":\"你好，世界\"}" \
  --output tts.mp3
```
期望：返回 MP3 音频文件（`audio/mpeg`）。若 502，请检查 `chat.baseURL` 与 `chat.apiKey`。

- 2) 文本指令自动化（TEXT）
```powershell
curl.exe -X POST "http://localhost:8000/client/text" \
  -H "Content-Type: application/json" \
  -d "{\"text\":\"打开记事本\"}"
```
期望：返回 `{ "res_text": "已完成" }`，同时在控制台日志看到自动化执行流程（应用启动、截图、OCR、动作）。

- 3) 语音识别（ASR → 自动化 → TTS）
```powershell
# 将本地 WAV/FLAC 文件转为 Base64（PowerShell）
$bytes = [IO.File]::ReadAllBytes("sample.wav")
$b64   = [Convert]::ToBase64String($bytes)
$body  = @{ audio_base64 = $b64; language = "auto" } | ConvertTo-Json

# 调用 ASR 接口
curl.exe -X POST "http://localhost:8000/client/asr" \
  -H "Content-Type: application/json" \
  -d $body \
  --output asr.mp3
```
期望：ASR 识别文本后触发桌面自动化，并将识别文本经 TTS 返回 MP3 音频。
提示：仅支持 `audio/wav` 或 `audio/flac`（见 `internal/logic/asr.go`）。

## 联调流程（端到端）
- 数据流：`前端 → Go → OCR → Java LLM → Go → 前端`
- 服务就绪检查：
  - OCR：`python files/OCR_api.py`（监听 `http://localhost:5000/ocr`，字段 `image`）
  - ASR：在 SenseVoice/FunASR 项目中运行 `api.py`（请确保监听 `http://localhost:50000/api/v1/asr`）
  - Java LLM：`cd backend/java && mvnw.cmd spring-boot:run`（默认 `http://localhost:8081/api`）
- Go 启动：如上（`go run .` 或 `main.exe`）。
- 验证：依次调用 `/client/text`、`/client/tts`、`/client/asr`，观察日志与返回。

## 常见问题与排查
- `400` 参数错误：
  - ASR 音频格式不支持（仅 `wav/flac`）；或 JSON 字段缺失。
- `502` 上游不可用：
  - TTS：`chat.baseURL`/`chat.apiKey` 配置错误或服务未启动。
  - ASR：`http://localhost:50000/api/v1/asr` 未就绪或返回异常。
  - OCR/LLM：对应服务未启动或接口变更。
- `500` 内部错误：
  - 自动化执行失败、截图/OCR 处理异常、模型服务返回结构不符合预期等。

## 路由与文件索引
- 入口：`internal/cmd/cmd.go`（CORS、`/ws`、`/client`）
- 控制器：`internal/controller/voice/voice_v1_*.go`
  - `voice_v1_voice.go` → `POST /client/tts`
  - `voice_v1_asr.go`   → `POST /client/asr`
  - `voice_v1_text.go`  → `POST /client/text`
- API 定义：`api/voice/v1/voice.go`
- 逻辑层：`internal/logic/*.go`
  - `voice.go`（TTS）、`asr.go`（ASR → 自动化 → TTS）、`robot.go`（动作执行）
  - `app.go`（应用启动）、`model_service.go`（Java LLM 交互）
- 工具层：`utility/*.go`（OCR 客户端与坐标标注、SSE 工具等）
- 外部微服务：`files/OCR_api.py`

## 备注
- 已开启 CORS，便于前端开发环境直接调用。
- 若需开放 OpenAPI/Swagger，请确保配置中设置 `openapiPath` 与 `swaggerPath` 并引入相应中间件。
- 后续可增加 ASR 地址为可配置项，避免硬编码。