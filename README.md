# Voice Controlled Desktop Assistant

一个基于 Electron + Vue3（Element Plus）与 Go 后端 + Spring AI 的桌面语音助理。前端使用 Vite 与 pnpm 管理；后端分层：
- Golang：负责本地自动化能力与系统集成（如录音 ASR、语音朗读 TTS、系统操作等）
- Java（Spring Boot + Spring AI）：负责与 LLM（大模型）交互，提供聊天/推理/工具调用等 AI 服务

部分自动化能力依赖 Windows 平台的 MSYS2 环境。

## 概览
- 前端：`Vue 3 + Vite + Element Plus + Electron`
- 后端：
  - `Go (GoFrame gf)`：本地自动化/设备能力（ASR/TTS/系统操作）
  - `Java (Spring Boot + Spring AI)`：LLM 服务（Chat/Reasoning/Tools）
- 包管理：`pnpm`
- 运行平台：Windows 10/11（开发与目标平台）

## 环境要求
- Node.js 18+（建议 LTS）
- pnpm 8+
- Go 1.21+（建议使用官方安装包）
- JDK 17+（用于 Spring Boot）
- Windows（必需）
- MSYS2（用于自动化任务与原生工具链，见下文）

安装 pnpm：
```powershell
npm i -g pnpm
```

## 前端开发（pnpm）
进入前端目录并安装依赖：
```powershell
cd frontend
pnpm install
```

启动开发模式（同时启动 Vite 与 Electron）：
```powershell
pnpm dev
```
- `dev:renderer`：仅启动 Vite（端口 `5173`）
- `dev:electron`：等待 `http://localhost:5173` 就绪后启动 Electron
- `dev`：并行启动两者（推荐）

构建与预览（浏览器预览渲染产物）：
```powershell
pnpm build
pnpm preview
```
> 生产打包 Electron 需要额外配置打包流程；当前 `pnpm start` 适用于在开发服务器（5173）就绪时启动 Electron。

## 后端（Golang：ASR/TTS 与本地自动化）
进入后端 Go 目录并启动：
```powershell
cd backend/golang
go mod tidy
go run .
```
编译可执行文件：
```powershell
go build -o dist/assistant.exe .
```
> 后端入口 `backend/golang/main.go` 使用 GoFrame 的命令体系（`cmd.Main.Run`）。服务端口与模块配置请参考 `backend/golang/internal` 目录具体实现。

### 能力边界（Golang）
- 录音（ASR）：采集音频并调用识别接口/服务
- 朗读（TTS）：播放合成语音，提供状态回调（开始/结束/错误）
- 系统操作：如应用启动、窗口控制、文件处理等（按需实现）
- 与前端/Electron 通讯：通过 IPC/HTTP/SSE 等方式暴露能力

## 后端（Spring AI：LLM 大模型服务）
> Java 服务位于 `backend/java`，用于承载 LLM 聊天/推理能力，前端通过 HTTP/SSE 接入。

### 启动方式
使用 Maven Wrapper（已提供 `mvnw`、`mvnw.cmd`）：
```powershell
cd backend/java
mvnw.cmd spring-boot:run
```
或使用本机 Maven：
```powershell
mvn spring-boot:run
```
默认端口通常为 `8080`，可在 `application.yml`/`application.properties` 修改。

### Spring AI 配置（示例）
在环境变量或配置文件中设置模型提供方的 Key 等参数（以 OpenAI 为例）：
```properties
# application.properties 示例
spring.ai.openai.api-key=${SPRING_AI_OPENAI_API_KEY}
spring.ai.openai.chat.options.model=gpt-4o-mini
spring.ai.retry.max-attempts=3
spring.ai.retry.backoff=500ms
```
如使用其他厂商（Azure OpenAI、DeepSeek、Moonshot 等），请参考 Spring AI 官方配置键并设置 API Key 与 Base URL。

### API 约定（建议）
- `POST /ai/chat`：标准聊天接口（JSON 请求/响应）
- `GET  /ai/manus/chat?message=...`：SSE 流式聊天（示例，与前端 `utils/sse/sse.ts` 中的路径一致）
- 根据你的控制器映射实际实现，前端 `SSEClient` 默认会连接以 `/ai/**` 开头的端点进行流式消费。

示例：测试 SSE 流式输出
```bash
curl -N -H "Accept: text/event-stream" "http://localhost:8080/ai/manus/chat?message=你好"
```
> 若端点命名不同，请同步调整前端 `frontend/utils/sse/sse.ts` 的连接路径。

### CORS 与联调
- 前端在开发时运行于 `http://localhost:5173`，请在 Java 服务端开启 CORS 允许该来源（`@CrossOrigin` 或 WebFlux/Spring MVC 全局 CORS 配置）。
- 若使用 SSE，请确保响应头与连接保持（`text/event-stream`、禁用缓存、适当的心跳/断线重连）。

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

## 运行与联调流程
1. 启动 Golang 后端：在 `backend/golang` 执行 `go run .`。
2. 启动 Spring AI 服务：在 `backend/java` 执行 `mvnw.cmd spring-boot:run`。
3. 启动前端：在 `frontend` 执行 `pnpm dev`（会同时启动 Vite 与 Electron）。
4. 交互验证：在 Electron 窗口中进行语音录制与识别，并发送文本/语音消息；前端通过 HTTP/SSE 与 Spring AI 服务进行 LLM 对话，Golang 负责 ASR/TTS 播放与本地自动化。
5. 若需要仅浏览器预览渲染产物，可在 `frontend` 执行 `pnpm build && pnpm preview`。

## 目录结构（简要）
```
qiniuyun
├── backend
│   ├── golang           Go 后端（ASR/TTS/系统能力）
│   ├── internal         业务与命令实现
│   ├── main.go          入口（GoFrame）
│   └── java             Spring Boot + Spring AI（LLM 服务）
│       ├── pom.xml        依赖与构建
│       └── src            控制器/服务实现
│           └── main\java\com\csjzsn\voicemindjavacore
│               ├── tools  工具类
│               ├── advisor  顾问（类似拦截器）
│               ├── app  项目核心实现
│               ├── common  通用实体类
│               ├── config  配置
│               ├── constant  常量
│               ├── controller  接口层
│               ├── exceptions  异常处理器
│               ├── model  实体模型
│               ├── rag  rag实现
│               └── VoiceMindJavaCoreApplication.java  启动类
└── frontend
    ├── src              前端源码（Vue3 + Element Plus）
    ├── electron         Electron 主进程与预加载脚本
    └── package.json     前端脚本与依赖（pnpm）+ react-app

```

## 许可证
本项目采用开源许可证（见根目录 `LICENSE`）。

## 参考
- [0] MSYS2 官方站点与安装指南：https://www.msys2.org/