# 语音控制桌面助手 （Voice Controlled Desktop Assistant）

一个基于 Electron + Vue3（Element Plus）与 Go 后端 + Spring AI 的桌面语音助理。前端使用 Vite 与 pnpm 管理；后端分层：

- Golang：负责本地自动化能力与系统集成（如录音 ASR、语音朗读 TTS、系统操作等）
- Java（Spring Boot + Spring AI）：负责与 LLM（大模型）交互，提供聊天/推理/工具调用等 AI 服务
- Python（微服务）：提供 OCR 图像文字识别（PaddleOCR + Flask），按需扩展其他能力

部分自动化能力依赖 Windows 平台的 MSYS2 环境。

演示视频：[七牛云demo_哔哩哔哩_bilibili](https://www.bilibili.com/video/BV1ZPszzJEoa/)

## 概览

- 前端：`Vue 3 + Vite + Element Plus + Electron`
- 后端：
  - `Go (GoFrame gf)`：本地自动化/设备能力（ASR/TTS/系统操作）
  - `Java (Spring Boot + Spring AI)`：LLM 服务（Chat/Reasoning/Tools）
  - `Python (Flask + PaddleOCR)`：OCR 微服务
- 包管理：`pnpm`
- 运行平台：Windows 10/11（开发与目标平台）

## 产品功能与优先级

- 核心功能（High）
  - 语音输入与识别（ASR）：麦克风音频→文本，作为指令与对话输入。
  - 文本到语音（TTS）：LLM 回复→语音播报，形成语音交互闭环。
  - LLM 对话（HTTP/SSE）：Java（Spring AI）承载推理，前端以 SSE 流式展示。
  - 桌面端 UI（Electron + Vue3）：录音控制、消息展示、状态与错误提示。
  - 配置与密钥管理：统一管理与切换模型提供方（Qwen/Gemini），开发环境 CORS/SSE。
- 增强功能（Medium）
  - OCR 能力（Python 微服务）：对截图/图片做文字识别，为指令与工具调用提供信息输入。
  - 基础系统操作桥接（Go）：在不引入第三方 agent 的前提下，按规则映射有限系统操作。
  - 会话管理：历史记录、复制/清空、导出（Markdown/JSON）。
  - 播放控制：TTS 中断、重播、音量与速率调整。
- 扩展功能（Low）
  - 多语言与方言支持（ASR/TTS）。
  - 离线/弱网模式（仅 ASR/TTS 的有限能力）。
  - 插件化工具接入（OCR、日历、剪贴板、截图）。
- 本次开发计划（MVP 范围）
  - 打通链路：ASR → Java LLM（SSE） → TTS 播放。
  - 双模型接入与切换：Qwen3-Max、Gemini-2.5-Pro-Thinking。
  - 集成 OCR 微服务：`POST /ocr`，结果可作为 LLM 上下文。
  - 前端基本交互：录音/停止、消息发送、流式渲染、错误与提示。

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

## 启动方式（推荐顺序）

1) MSYS2 安装与使用（Windows）

- 部分自动化操作（例如调用本地工具、脚本、打包或某些原生依赖）需要 MSYS2 环境。
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

2) SenseVoice 语音模型使用方法

- 作为 ASR 微服务，直接使用官方项目的 `api.py` 启动服务，详见下文「SenseVoice 语音模型使用方法」。
- 项目链接：`https://github.com/FunAudioLLM/SenseVoice`

3) OCR 微服务

- 在 `backend/golang/files/OCR_api.py` 目录启动 Flask + PaddleOCR 服务。
- 该服务基于 Flask + PaddleOCR，提供简单的图片文字识别 API。
- 安装

```powershell
cd backend/golang/files
python -m venv .venv
.\.venv\Scripts\activate
pip install --upgrade pip
pip install flask paddleocr
```

- 首次运行会下载 PaddleOCR 模型，需耐心等待。
- 启动

```powershell
python OCR_api.py
# 监听 0.0.0.0:5000
```

* 调用示例

```bash
curl -F "image=@path/to/image.png" http://localhost:5000/ocr
```

- 返回为 JSON 数组，每个元素包含检测框坐标与文字等信息；服务会在 `output/` 目录保存识别标注图片与 JSON 文件。

4) 安装后端依赖（Golang 与 Java）

- Golang：在 `backend/golang` 执行 `go mod tidy` 或 `go mod download` 安装项目依赖。
- Java：在 `backend/java` 执行 `mvnw.cmd clean package -DskipTests` 下载并缓存 Maven 依赖。

5) 启动 Golang 后端

- 进入 `backend/golang` 并运行 `go run .` 或编译；后端默认调用 ASR/TTS 微服务。

6) 启动 Java LLM 服务

- 进入 `backend/java` 并运行 Spring Boot；前端通过 SSE 获取流式响应。

7) 启动前端

- 进入 `frontend` 运行 `pnpm install && pnpm dev`，浏览器访问提示的 `http://localhost:<port>/`。

---

## 如何运行程序（完整流程）

1) 启动 Golang 后端（ASR/TTS/系统能力）

```powershell
cd backend/golang
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

- 详细步骤与命令见下文「OCR 微服务安装与启动」。
- 默认监听 `http://0.0.0.0:5000/ocr`。

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

### MVP 验收标准

- 语音 → 文本：前端录音后，可成功识别并显示文本。
- 文本 → LLM：将文本发送到 LLM，前端以 SSE 流式展示。
- LLM → 语音：将回复合成并在桌面播报，可控制停止/重播。
- OCR：上传图片到 `POST /ocr`，结果在前端展示并可作为上下文输入。
- 双模型切换：通过配置选择 Qwen 或 Gemini，服务能正确返回结果。

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

## 实现挑战与应对

- Windows 音频采集与设备兼容性
  - 统一采样参数（如 16k/mono）、提供设备列表与默认设备回退；对错误做明确提示与日志。
- SSE 流式与 CORS
  - 开启全局 CORS；SSE 响应确保 `text/event-stream`、禁用缓存与 keep-alive；前端断线重连与错误提示。
- 双模型接入差异（Qwen 与 Gemini）
  - 在 Java 服务抽象 `LLMProvider`，分别实现 Qwen/Gemini 适配；以配置/环境变量切换，统一输出 SSE/JSON。
- 性能与资源使用
  - 使用流式；OCR 首次启动提示模型下载；对非核心路径使用队列与缓存；UI 显示加载进度与降级策略。
- 无第三方 agent 的编排
  - 采用“规则驱动”的轻量编排：约定 LLM 输出的工具触发格式（如 `tool: ocr` 或 JSON 标记），Java 服务检测并调用对应工具，并回填上下文；工具列表有限且接口稳定。
- 交互与可用性
  - 前端统一状态机，明确阶段（录音中/识别中/生成中/播报中）；失败场景提供重试与回退文案。

---

## SenseVoice 语音模型使用方法

项目链接：`https://github.com/FunAudioLLM/SenseVoice`

为便于集成，请直接启动 SenseVoice 项目中的 `api.py` 作为 ASR 微服务。

### 快速启动

```powershell
# 克隆并进入项目
git clone https://github.com/FunAudioLLM/SenseVoice.git
cd SenseVoice

# 建议使用虚拟环境
python -m venv .venv
.\.venv\Scripts\activate
pip install --upgrade pip
pip install -r requirements.txt

# 启动服务（默认即可）
python api.py
```

- 服务启动后，请确保监听端口与 Go 后端一致（推荐 `50000`），并提供 `POST /api/v1/asr` 接口（`multipart/form-data`，字段 `files`，可选参数 `language=auto`），以对齐后端的默认调用。
- 若 `api.py` 的默认路由或端口不同：
  - 方案 A：在 `api.py` 中添加兼容路由 `POST /api/v1/asr`；
  - 方案 B：调整 Go 后端 `internal/logic/asr.go` 中的 ASR 目标地址为实际路由。

### 对接示例（本地文件）

```powershell
# 将音频文件以 multipart/form-data 上传到 SenseVoice 服务
curl -X POST http://localhost:50000/api/v1/asr `
  -F "files=@sample.wav" `
  -F "language=auto"
```

- Go 后端会以相同方式调用该服务；前端录音通过后端转发到 ASR，结果以文本返回并进入后续 LLM/TTS 流程。
- 推荐将 SenseVoice 服务与 OCR 服务并行启动，提升端到端响应速度。

---

## LLM 模型采纳与对比

- 采纳方案
  - 首选：Qwen3-Max（通义千问）——中文能力强、响应速度与成本均衡、国内网络可用性更佳，工具调用能力好。
  - 增强：Gemini-2.5-Pro-Thinking（Google）——复杂推理与结构化思考能力更强，作为高难度问题的 fallback。
- 对比维度（简要）
  - 语言能力：Qwen 在中文语境与工具调用表达更优；Gemini 在英语/跨领域推理更强。
  - 成本与可用性：Qwen 在国内可用性与成本更友好；Gemini 在复杂推理场景价值更高。
  - 接入与生态：二者均支持 REST/SSE；服务端以提供方适配器方式统一对外接口。
- 具体接入（配置建议）
  - Java 服务层提供 `LLM_PROVIDER` 配置（`QWEN`/`GEMINI`）。
  - Qwen：环境变量 `DASHSCOPE_API_KEY`。
  - Gemini：环境变量 `GOOGLE_API_KEY`。

---

## 分工建议（角色与职责）

- 前端工程（Electron/Vue）：UI/交互、消息流、与后端通信（HTTP/SSE）、录音上传与结果渲染
- 后端工程（Golang）：ASR/TTS/系统操作能力实现与 API 暴露、与 Python/Java 服务编排
- 智能工程（Java + Spring AI）：LLM 对话/推理、工具调用编排、SSE 流式输出与 CORS 配置
- 算法/ML 工程（Python）：SenseVoice/FunASR 语音模型的服务化与优化，OCR/其他 CV 能力微服务

### 人员与分工

- Golang（本地自动化/ASR/TTS）：VH992098059（详见 [docs/Golang架构设计文档.md](docs/Golang架构设计文档.md)）
- Java（LLM 服务与控制器）：CSJZSN（详见 [backend/java/架构设计文档.md](backend/java/架构设计文档.md)）

---

## 未来规划

- 减少处理等待时间：全链路优化（ASR/LLM/TTS）、模型预热与缓存、并发与异步管线。
- 唤醒词与持续监听（可选）：提升免手动语音交互体验，提供能耗与隐私控制开关。
- 本地化模型选项（ASR/TTS 轻量版）：弱网场景下保持基本功能可用，增强隐私与成本可控。
- 上下文记忆与轻量知识库：支持本地笔记/文件检索与上下文补充。
- 插件化工具接口：标准化工具（OCR、Screenshot、Calendar、Clipboard）以 HTTP/IPC 方式接入。
- 跨平台打包（Win/Mac）：与 Electron 打包流程打通。
- 安全与隐私：更细粒度的日志与开关，明确数据流向与敏感信息最小化原则。

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

## 免责声明与约束

- 不调用第三方 agent 能力，所有工具编排由自研的“有限规则”实现。
- 仅允许调用 LLM 模型、ASR 与 TTS 能力；扩展其他工具需遵守相同约束与接口规范。
- 模型与服务的选择需依据实际密钥与可用性；文档中的模型名与配置以官方文档为准。

---

## 许可证

本项目采用开源许可证（见根目录 `LICENSE`）。

## 参考

- [0] MSYS2 官方站点与安装指南：https://www.msys2.org/
- FunASR 与 SenseVoice（ModelScope）：https://modelscope.cn/
- PaddleOCR：https://github.com/PaddlePaddle/PaddleOCR
