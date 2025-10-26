# 语音控制桌面助手 Java Core 模块规格说明

## 模块概览

语音控制桌面助手 Java Core 系统采用模块化设计，每个模块都有明确的职责和接口规范。本文档详细说明各个模块的功能、接口、依赖关系和实现细节。

## 模块架构图

```
┌─────────────────────────────────────────────────────────────┐
│                    VoiceMind Java Core                     │
├─────────────────────────────────────────────────────────────┤
│  Controller Layer    │  Application Layer  │  AI Service Layer │
├─────────────────────────────────────────────────────────────┤
│  Tool Integration   │  RAG Layer         │  Infrastructure   │
└─────────────────────────────────────────────────────────────┘
```

## 1. 控制器层模块 (Controller Layer)

### 1.1 AiController 模块

**模块路径**: `com.csjzsn.voicemindjavacore.controller.AiController`

**职责**: 提供AI聊天相关的REST API接口

**主要接口**:

| 接口     | 方法 | 路径                | 功能描述             |
| -------- | ---- | ------------------- | -------------------- |
| 基础聊天 | POST | `/ai/chat/do`     | 提供基础AI对话功能   |
| 报告聊天 | POST | `/ai/chat/report` | 返回结构化分析报告   |
| RAG聊天  | POST | `/ai/chat/rag`    | 基于知识库的增强对话 |
| 工具聊天 | POST | `/ai/chat/tools`  | 集成工具调用的对话   |
| MCP聊天  | POST | `/ai/chat/mcp`    | 基于MCP协议的对话    |
| 流式聊天 | GET  | `/ai/chat/sse`    | 实时流式响应         |

**依赖模块**:

- `OrderApp`: 核心业务逻辑
- `ToolCallback[]`: 工具回调数组
- `ChatModel`: AI模型接口

**输入参数**:

```java
public class ProcessRequest {
    private String text;        // 用户输入文本
    private String session_id;  // 会话ID
}
```

**输出格式**:

```java
public class BaseResponse<T> {
    private boolean success;    // 操作是否成功
    private String message;     // 响应消息
    private T data;            // 响应数据
}
```

### 1.2 HealthController 模块

**模块路径**: `com.csjzsn.voicemindjavacore.controller.HealthController`

**职责**: 提供系统健康检查接口

**主要接口**:

- `GET /health`: 返回系统健康状态

### 1.3 MultiController 模块

**模块路径**: `com.csjzsn.voicemindjavacore.controller.MultiController`

**职责**: 提供多模态处理接口

**主要功能**:

- 图像处理
- 多模态输入处理

## 2. 应用服务层模块 (Application Layer)

### 2.1 OrderApp 模块

**模块路径**: `com.csjzsn.voicemindjavacore.app.OrderApp`

**职责**: 核心AI聊天应用服务，提供多种聊天模式

**核心功能**:

#### 2.1.1 基础聊天功能

```java
public String doChat(String message, String chatId)
```

- **功能**: 基础AI对话，支持多轮对话记忆
- **参数**:
  - `message`: 用户输入消息
  - `chatId`: 会话ID
- **返回**: AI响应文本

#### 2.1.2 流式聊天功能

```java
public Flux<String> doChatByStream(String message, String chatId)
```

- **功能**: 流式AI对话，实时返回响应
- **参数**: 同基础聊天
- **返回**: 响应流

#### 2.1.3 报告生成功能

```java
public OrderReport doChatWithReport(String message, String chatId)
```

- **功能**: 生成结构化分析报告
- **返回**: `OrderReport` 对象

#### 2.1.4 RAG增强功能

```java
public String doChatWithRag(String message, String chatId)
```

- **功能**: 基于知识库的检索增强生成
- **特性**: 查询重写、向量检索、上下文增强

#### 2.1.5 工具集成功能

```java
public String doChatWithTools(String message, String chatId)
```

- **功能**: 集成各种工具调用的AI对话
- **支持工具**: 文件操作、网络搜索、终端操作等

#### 2.1.6 MCP协议功能

```java
public String doChatWithMcp(String message, String chatId)
```

- **功能**: 基于MCP协议的AI对话
- **特性**: 支持外部工具和服务调用

**依赖模块**:

- `ChatModel`: AI模型
- `VectorStore`: 向量存储
- `QueryRewriter`: 查询重写器
- `ToolCallback[]`: 工具回调
- `ToolCallbackProvider`: 工具回调提供者

**配置参数**:

```java
private static final String SYSTEM_PROMPT = """
    你是一个智能API接口设计助手，负责处理用户语音指令并生成标准化的JSON响应。
    请严格遵守以下规范:
    将用户自然语言指令解析为结构化操作意图
    生成前后端通用的标准化响应格式
    确保所有操作在安全许可范围内
    对模糊指令进行合理推断或要求澄清
    """;
```

### 2.2 ImageApp 模块

**模块路径**: `com.csjzsn.voicemindjavacore.app.ImageApp`

**职责**: 图像处理应用服务

**主要功能**:

- 图像识别
- 图像处理
- 图像分析

## 3. AI服务层模块 (AI Service Layer)

### 3.1 模型配置模块

**配置路径**: `src/main/resources/application-local.yml`

**支持的AI模型**:

#### 3.1.1 通义千问 (DashScope)

```yaml
spring:
  ai:
    dashscope:
      api-key: your-api-key
      chat:
        options:
          model: qwen3-max
```

#### 3.1.2 Gemini

```yaml
gemini:
  api-key: your-gemini-key
  base-url: https://poloai.top/v1
  model: gemini-2.5-pro-thinking
```

### 3.2 ChatClient 配置

**核心配置**:

- 系统提示词配置
- 对话记忆管理
- Advisor链配置
- 工具回调集成

## 4. 工具集成层模块 (Tool Integration Layer)

### 4.1 文件操作工具模块

#### 4.1.1 FileOperationTool

**模块路径**: `com.csjzsn.voicemindjavacore.tools.file.FileOperationTool`

**功能**:

- 文件创建、读取、写入、删除
- 目录操作
- 文件权限管理

**接口规范**:

```java
@Tool("文件操作工具")
public class FileOperationTool {
    public String createFile(String path, String content);
    public String readFile(String path);
    public String writeFile(String path, String content);
    public String deleteFile(String path);
}
```

#### 4.1.2 PDFGenerationTool

**模块路径**: `com.csjzsn.voicemindjavacore.tools.file.PDFGenerationTool`

**功能**:

- PDF文档生成
- 文本转PDF
- 图片转PDF

#### 4.1.3 ResourceDownloadTool

**模块路径**: `com.csjzsn.voicemindjavacore.tools.file.ResourceDownloadTool`

**功能**:

- 网络资源下载
- 文件下载管理
- 下载进度跟踪

### 4.2 网络工具模块

#### 4.2.1 WebSearchTool

**模块路径**: `com.csjzsn.voicemindjavacore.tools.web.WebSearchTool`

**功能**:

- 网络搜索
- 搜索结果处理
- 搜索API集成

**配置参数**:

```java
@Value("${search-api.api-key}")
private String searchApiKey;
```

#### 4.2.2 WebScrapingTool

**模块路径**: `com.csjzsn.voicemindjavacore.tools.web.WebScrapingTool`

**功能**:

- 网页内容抓取
- HTML解析
- 数据提取

### 4.3 系统工具模块

#### 4.3.1 TerminalOperationTool

**模块路径**: `com.csjzsn.voicemindjavacore.tools.system.TerminalOperationTool`

**功能**:

- 终端命令执行
- 系统操作
- 进程管理

#### 4.3.2 TerminateTool

**模块路径**: `com.csjzsn.voicemindjavacore.tools.control.TerminateTool`

**功能**:

- 进程终止
- 任务管理
- 资源清理

### 4.4 工具注册模块

#### 4.4.1 ToolRegistration

**模块路径**: `com.csjzsn.voicemindjavacore.tools.control.ToolRegistration`

**职责**: 集中管理所有工具注册

**配置**:

```java
@Configuration
public class ToolRegistration {
    @Bean
    public ToolCallback[] allTools() {
        return ToolCallbacks.from(
            fileOperationTool,
            webSearchTool,
            webScrapingTool,
            resourceDownloadTool,
            terminalOperationTool,
            pdfGenerationTool,
            terminateTool
        );
    }
}
```

## 5. RAG知识库层模块 (RAG Layer)

### 5.1 向量存储模块

#### 5.1.1 OrderAppVectorStoreConfig

**模块路径**: `com.csjzsn.voicemindjavacore.rag.OrderAppVectorStoreConfig`

**功能**:

- 向量存储配置
- 向量数据库连接
- 向量索引管理

### 5.2 文档处理模块

#### 5.2.1 OrderAppDocumentLoader

**模块路径**: `com.csjzsn.voicemindjavacore.rag.OrderAppDocumentLoader`

**功能**:

- 文档加载
- 多格式文档支持
- 文档预处理

#### 5.2.2 MyTokenTextSplitter

**模块路径**: `com.csjzsn.voicemindjavacore.rag.MyTokenTextSplitter`

**功能**:

- 文本分割
- Token管理
- 分块策略

#### 5.2.3 MyKeywordEnricher

**模块路径**: `com.csjzsn.voicemindjavacore.rag.MyKeywordEnricher`

**功能**:

- 关键词提取
- 关键词增强
- 语义分析

### 5.3 查询处理模块

#### 5.3.1 QueryRewriter

**模块路径**: `com.csjzsn.voicemindjavacore.rag.QueryRewriter`

**功能**:

- 查询重写
- 查询优化
- 语义理解

**接口**:

```java
public class QueryRewriter {
    public String doQueryRewrite(String originalQuery);
}
```

## 6. 基础设施层模块 (Infrastructure Layer)

### 6.1 异常处理模块

#### 6.1.1 GlobalExceptionHandler

**模块路径**: `com.csjzsn.voicemindjavacore.exceptions.GlobalExceptionHandler`

**功能**:

- 全局异常处理
- 异常统一响应
- 错误日志记录

**处理的异常类型**:

- `ActionExecutionException`: 操作执行异常
- `CustRuntimeException`: 自定义运行时异常
- `InvalidRequestException`: 无效请求异常
- `LLMServiceException`: LLM服务异常

#### 6.1.2 自定义异常类

| 异常类                       | 用途           | 触发条件       |
| ---------------------------- | -------------- | -------------- |
| `ActionExecutionException` | 工具执行失败   | 工具调用异常   |
| `CustRuntimeException`     | 通用运行时异常 | 业务逻辑错误   |
| `InvalidRequestException`  | 请求参数错误   | 参数验证失败   |
| `LLMServiceException`      | AI服务异常     | AI模型调用失败 |

### 6.2 配置管理模块

#### 6.2.1 CrossConfig

**模块路径**: `com.csjzsn.voicemindjavacore.config.CrossConfig`

**功能**:

- 跨域配置
- CORS策略管理
- 安全配置

#### 6.2.2 Knife4jConfig

**模块路径**: `com.csjzsn.voicemindjavacore.config.Knife4jConfig`

**功能**:

- API文档配置
- Swagger集成
- 接口文档生成

#### 6.2.3 FileConstant

**模块路径**: `com.csjzsn.voicemindjavacore.constant.FileConstant`

**功能**:

- 文件常量定义
- 路径配置
- 文件类型管理

### 6.3 日志和监控模块

#### 6.3.1 MyLoggerAdvisor

**模块路径**: `com.csjzsn.voicemindjavacore.advisor.MyLoggerAdvisor`

**功能**:

- 自定义日志增强
- 请求响应日志
- 性能监控

#### 6.3.2 ReReadingAdvisor

**模块路径**: `com.csjzsn.voicemindjavacore.advisor.ReReadingAdvisor`

**功能**:

- 推理增强
- 思维链优化
- 响应质量提升

## 7. 数据模型模块 (Data Model Layer)

### 7.1 请求响应模型

#### 7.1.1 ProcessRequest

**模块路径**: `com.csjzsn.voicemindjavacore.model.dto.ProcessRequest`

**字段**:

```java
public class ProcessRequest {
    private String text;        // 用户输入文本
    private String session_id;  // 会话ID
}
```

#### 7.1.2 BaseResponse

**模块路径**: `com.csjzsn.voicemindjavacore.common.BaseResponse`

**字段**:

```java
public class BaseResponse<T> {
    private boolean success;    // 操作是否成功
    private String message;     // 响应消息
    private T data;            // 响应数据
}
```

### 7.2 业务模型

#### 7.2.1 OrderReport

**模块路径**: `com.csjzsn.voicemindjavacore.common.OrderReport`

**功能**: 结构化分析报告

**内部类**:

- `Backend`: 后端相关分析
- `Command`: 命令分析
- `CommandSequence`: 命令序列
- `Data`: 数据分析
- `Dependency`: 依赖分析
- `Execution`: 执行分析
- `Frontend`: 前端相关分析
- `GlobalSafetyCheck`: 全局安全检查
- `Parameters`: 参数分析
- `SafetyCheck`: 安全检查

## 8. 模块依赖关系

### 8.1 依赖图

```
Controller Layer
    ↓
Application Layer
    ↓
AI Service Layer
    ↓
Tool Integration Layer
    ↓
RAG Layer
    ↓
Infrastructure Layer
```

### 8.2 模块间接口

| 调用方           | 被调用方       | 接口类型 | 说明         |
| ---------------- | -------------- | -------- | ------------ |
| AiController     | OrderApp       | 方法调用 | 业务逻辑处理 |
| OrderApp         | ChatModel      | 接口调用 | AI模型服务   |
| OrderApp         | ToolCallback[] | 接口调用 | 工具集成     |
| OrderApp         | VectorStore    | 接口调用 | 向量存储     |
| ToolRegistration | 各种Tool       | 配置注入 | 工具注册     |

## 9. 模块扩展指南

### 9.1 添加新工具

1. 创建工具类，实现相应接口
2. 在 `ToolRegistration` 中注册
3. 配置工具参数
4. 测试工具功能

### 9.2 添加新AI模型

1. 在配置文件中添加模型配置
2. 创建模型Bean
3. 在 `OrderApp` 中集成
4. 测试模型调用

### 9.3 添加新RAG功能

1. 实现文档加载器
2. 配置向量存储
3. 创建查询处理器
4. 集成到聊天流程

## 10. 模块测试规范

### 10.1 单元测试

每个模块都应包含：

- 功能测试
- 边界条件测试
- 异常情况测试
- 性能测试

### 10.2 集成测试

- 模块间接口测试
- 端到端功能测试
- 性能基准测试

### 10.3 测试覆盖率

- 代码覆盖率 > 80%
- 分支覆盖率 > 70%
- 方法覆盖率 > 90%

## 总结

语音控制桌面助手 Java Core 的模块化设计确保了系统的可维护性、可扩展性和可测试性。每个模块都有明确的职责边界和接口规范，通过依赖注入和配置管理实现了松耦合的架构。这种设计使得系统能够灵活地适应不同的业务需求和技术变化。
