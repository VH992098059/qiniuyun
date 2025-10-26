# Voice Controlled Desktop Assistant

一个基于 Electron + Vue3（Element Plus）与 Go 后端 + Spring AI 的桌面语音助理。前端使用 Vite 与 pnpm 管理；后端分层：
- Golang：负责本地自动化能力与系统集成（如录音 ASR、语音朗读 TTS、系统操作等）
- Java（Spring Boot + Spring AI）：负责与 LLM（大模型）交互，提供聊天/推理/工具调用等 AI 服务
- Python（微服务）：提供 OCR 图像文字识别（PaddleOCR + Flask），按需扩展其他能力

部分自动化能力依赖 Windows 平台的 MSYS2 环境。

## 概览
- 前端：`Vue 3 + Vite + Element Plus + Electron`
- 后端：
  - `Go (GoFrame gf)`：本地自动化/设备能力（ASR/TTS/系统操作）
  - `Java (Spring Boot + Spring AI)`：LLM 服务（Chat/Reasoning/Tools）
  - `Python (Flask + PaddleOCR)`：OCR 微服务
- 包管理：`pnpm`
- 运行平台：Windows 10/11（开发与目标平台）

## 环境要求
- Node.js 18+（建议 LTS）
- pnpm 8+
- Go 1.21+（建议使用官方安装包）
- JDK 17+（用于 Spring Boot）
- Python 3.9+（用于 OCR 微服务与可选 ASR）
- Windows（必需）
- MSYS2（用于自动化任务与原生工具链，见下文）

安装 pnpm：
```powershell
npm i -g pnpm
```

---

## 如何运行程序（完整流程）

1) 启动 Golang 后端（ASR/TTS/系统能力）
```powershell
cd backend/golang
go mod tidy
go run .
# 或编译
# go build -o dist/assistant.exe .
```
默认使用 GoFrame 命令体系（入口 `backend/golang/main.go`）。端口与模块配置见 `backend/golang/internal`。

2) 启动 Java LLM 服务（Spring AI）
```powershell
cd backend/java
# 使用 Maven Wrapper
mvnw.cmd spring-boot:run
# 或使用本机 Maven
# mvn spring-boot:run
```
- 端口默认 `8080`（可在 `application.yml`/`application.properties` 修改）。
- 开启 CORS 以允许 `http://localhost:5173` 来源（开发环境）。
- SSE 端点参考：`GET /ai/manus/chat?message=...`（下文“架构与模块规格”有细节）。

3) 启动 OCR 微服务（Python + Flask + PaddleOCR）
```powershell
cd backend/golang/files
# 建议创建虚拟环境
python -m venv .venv
.\.venv\Scripts\activate
pip install --upgrade pip
pip install flask paddleocr
# 启动服务
python OCR_api.py
```
- 默认监听 `http://0.0.0.0:5000/ocr`
- 使用方法与返回格式见下文“OCR 微服务安装与启动”

4) 启动前端（Vite + Electron）
```powershell
cd frontend
pnpm install
pnpm dev
```
- `dev:renderer`：仅启动 Vite（端口 `5173`）
- `dev:electron`：等待 `http://localhost:5173` 就绪后启动 Electron
- `dev`：并行启动两者（推荐）

5) 联调验证
- 在 Electron 窗口进行语音录制/识别，发送文本/语音消息
- 前端通过 HTTP/SSE 与 Java LLM 对话
- Golang 负责 ASR/TTS 播放与本地自动化
- OCR 服务可通过后端或前端调用其 HTTP 接口实现图像文字识别

> 若仅浏览器预览渲染产物，可在 `frontend` 执行 `pnpm build && pnpm preview`

---

## 架构设计与模块规格

### 总览
- UI 层（Electron + Vue3）：
  - 会话窗口、录音控制、消息展示、工具结果展示
  - 通过 HTTP/SSE 调用后端服务
- 能力层（Go）：
  - 提供 ASR/TTS/系统操作能力，封装为可供前端调用的 API/IPC
  - 可与 Python 微服务交互（如 OCR）
- 智能层（Java + Spring AI）：
  - 统一 LLM 对话与推理能力，支持流式（SSE）输出
  - 可编排工具调用（如 OCR、系统操作等）
- 微服务（Python）：
  - 独立部署的 OCR 服务，面向 HTTP 调用；可拓展其他 CV/ASR 能力

### 进程与通信
- 前端（`localhost:5173`）⇄ Java LLM（`localhost:8080`）：
  - `GET /ai/manus/chat?message=...`（SSE，`text/event-stream`）
  - `POST /ai/chat`（JSON，同步）
- 前端/Go ⇄ Python OCR（`localhost:5000`）：
  - `POST /ocr`（`multipart/form-data`，字段 `image`）
- 前端 ⇄ Go：
  - IPC/HTTP（按项目实现，将 ASR/TTS/系统操作能力暴露为 API）

### 模块规格（示例）
- LLM 控制器（Java）：
  - `POST /ai/chat`：请求体 `{ message: string, ... }`；返回标准 JSON 回复
  - `GET /ai/manus/chat?message=...`：SSE 流式输出，保持连接、禁用缓存
- OCR 服务（Python）：
  - `POST /ocr`：上传图片；返回 PaddleOCR 结构化结果数组（含检测框与文字）
- Go 能力接口（示例约定）：
  - `POST /audio/asr`：提交音频或触发录音，返回识别文本
  - `POST /audio/tts`：提交文本，播放合成语音，返回播放状态
  - `POST /system/exec`：触发系统操作（启动应用、窗口控制等）

> 以上接口为建议与示例；请根据实际实现同步调整前端 `utils/sse/sse.ts` 与调用路径。

---

## SenseVoice 语音模型使用方法（示例）

本项目的 ASR（语音识别）可选使用「SenseVoice」模型进行本地或服务化推理。以下为在 Windows 的示例方案（基于 ModelScope/FunASR 生态）：

### 安装与准备
```powershell
# 建议使用独立虚拟环境
python -m venv sensevoice-venv
.sensevoice-venv\Scripts\activate
pip install --upgrade pip
# 安装依赖（示例）
pip install modelscope funasr
```

### 本地推理（示例脚本）
```python
from modelscope.pipelines import pipeline
from modelscope.utils.constant import Tasks

# 选择合适的 SenseVoice 模型（示例名称，需根据可用模型调整）
asr_pipeline = pipeline(task=Tasks.auto_speech_recognition,
                        model='iic/sensevoice_small',  # 示例模型名
                        model_revision='v1.0.0')

res = asr_pipeline({'audio': 'sample.wav'})
print(res)
```
- 将识别逻辑封装到 Go 的 `POST /audio/asr` 接口中，或以 Python 微服务形式暴露 HTTP 接口供 Go 调用。
- 若使用流式识别，可参考 FunASR 的 streaming/server 示例（命令可能为 `funasr-server` 或相应 Python 启动脚本，按官方文档配置模型与端口）。

> 模型名称与启动命令可能因版本/厂商不同而异，请以 ModelScope/FunASR/SenseVoice 官方文档为准，并确保在 Windows 环境安装所需的 Runtime 依赖（如 `onnxruntime`）。

### 与本项目集成建议
- Go 层通过 HTTP 调用本地 SenseVoice 微服务，或直接在 Go 中以子进程方式调用 Python 脚本并返回识别结果。
- 前端录音文件通过 Go 上传给 ASR 接口，最终文本回传给前端页面。

---

## OCR 微服务安装与启动（backend/golang/files/OCR_api.py）

该服务基于 Flask + PaddleOCR，提供简单的图片文字识别 API。

### 安装
```powershell
cd backend/golang/files
python -m venv .venv
.\.venv\Scripts\activate
pip install --upgrade pip
pip install flask paddleocr
```
- 首次运行会下载 PaddleOCR 模型，需耐心等待。

### 启动
```powershell
python OCR_api.py
# 监听 0.0.0.0:5000
```

### 调用示例
```bash
curl -F "image=@path/to/image.png" http://localhost:5000/ocr
```
- 返回为 JSON 数组，每个元素包含检测框坐标与文字等信息；服务会在 `output/` 目录保存识别标注图片与 JSON 文件。

---

## Spring AI 配置（示例）
在环境变量或配置文件中设置模型提供方的 Key 等参数（以 OpenAI 为例）：
```properties
# application.properties 示例
spring.ai.openai.api-key=${SPRING_AI_OPENAI_API_KEY}
spring.ai.openai.chat.options.model=gpt-4o-mini
spring.ai.retry.max-attempts=3
spring.ai.retry.backoff=500ms
```
如使用其他厂商（Azure OpenAI、DeepSeek、Moonshot 等），请参考 Spring AI 官方配置键并设置 API Key 与 Base URL。

### SSE 测试示例
```bash
curl -N -H "Accept: text/event-stream" "http://localhost:8080/ai/manus/chat?message=你好"
```

### CORS 与联调
- 前端在开发时运行于 `http://localhost:5173`，请在 Java 服务端开启 CORS 允许该来源（`@CrossOrigin` 或 WebFlux/Spring MVC 全局 CORS 配置）。
- 若使用 SSE，请确保响应头与连接保持（`text/event-stream`、禁用缓存、适当的心跳/断线重连）。

---

## MSYS2 安装与使用（Windows）
部分自动化操作（例如调用本地工具、脚本、打包或某些原生依赖）需要 MSYS2 环境。

- 下载与安装：从官方网站下载安装程序并按照指引安装。[0]
- 安装位置建议：选择简短的 ASCII 路径、位于 NTFS 分区、避免空格与符号链接（例如 `C:\msys64`）。[0]
- 初次打开会启动 UCRT64 终端环境（建议使用该环境）。[0]
- 更新系统与包：
```bash
pacman -Syu
```
- 安装 MinGW-w64 工具链（示例）：
```bash
pacman -S --needed base-devel mingw-w64-x86_64-toolchain
```
- 将 `C:\msys64\ucrt64\bin`（或实际安装路径）加入系统 `PATH`，以便在脚本或 Go 程序中调用相关工具。

> 以上 MSYS2 信息来源于其官方说明与安装指南。[0]

---

## 分工建议（角色与职责）
- 前端工程（Electron/Vue）：UI/交互、消息流、与后端通信（HTTP/SSE）、录音上传与结果渲染
- 后端工程（Golang）：ASR/TTS/系统操作能力实现与 API 暴露、与 Python/Java 服务编排
- 智能工程（Java + Spring AI）：LLM 对话/推理、工具调用编排、SSE 流式输出与 CORS 配置
- 算法/ML 工程（Python）：SenseVoice/SenseTime/FunASR 等语音模型的服务化与优化，OCR/其他 CV 能力微服务

---

## 目录结构（简要）
```
backend/
  golang/           Go 后端（ASR/TTS/系统能力）
    internal/       业务与命令实现
    main.go         入口（GoFrame）
  java/             Spring Boot + Spring AI（LLM 服务）
    pom.xml         依赖与构建
    src/            控制器/服务实现
frontend/
  src/              前端源码（Vue3 + Element Plus）
  electron/         Electron 主进程与预加载脚本
  package.json      前端脚本与依赖（pnpm）
```

---

## 许可证
本项目采用开源许可证（见根目录 `LICENSE`）。

## 参考
- [0] MSYS2 官方站点与安装指南：https://www.msys2.org/
- FunASR 与 SenseVoice（ModelScope）：https://modelscope.cn/
- PaddleOCR：https://github.com/PaddlePaddle/PaddleOCR