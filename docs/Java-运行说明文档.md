# VoiceMind Java Core 运行说明文档

## 项目简介

VoiceMind Java Core 是一个基于Spring AI框架的智能语音助手后端系统，提供多种AI聊天模式、工具集成、RAG检索增强生成等功能。

## 环境要求

### 系统要求
- **操作系统**: Windows 10/11, macOS, Linux
- **Java版本**: Java 21 或更高版本
- **内存**: 建议 4GB 以上
- **磁盘空间**: 至少 1GB 可用空间

### 开发环境
- **IDE**: IntelliJ IDEA, VS Code
- **构建工具**: Maven 3.6+
- **版本控制**: Git

## 快速开始

### 1. 克隆项目

```bash
git clone <repository-url>
cd VoiceMind-Java-core
```

### 2. 配置环境变量

#### 2.1 配置AI模型API密钥

编辑 `src/main/resources/application-local.yml` 文件：

```yaml
spring:
  ai:
    dashscope:
      api-key: your-dashscope-api-key  # 替换为您的通义千问API密钥
gemini:
  api-key: your-gemini-api-key        # 替换为您的Gemini API密钥
  base-url: https://poloai.top/v1     # Gemini API基础URL
  model: gemini-2.5-pro-thinking      # 使用的模型名称
search-api:
  api-key: your-search-api-key        # 替换为您的搜索API密钥
```

#### 2.2 获取API密钥

1. **通义千问API密钥**:
   - 访问 [阿里云DashScope控制台](https://dashscope.console.aliyun.com/)
   - 注册/登录账号
   - 创建API密钥

2. **Gemini API密钥**:
   - 访问 [Google AI Studio](https://aistudio.google.com/)
   - 获取API密钥

3. **搜索API密钥**:
   - 根据使用的搜索服务提供商获取相应密钥

### 3. 编译项目

```bash
# 使用Maven编译
mvn clean compile

# 或者使用Maven Wrapper (推荐)
./mvnw clean compile  # Linux/macOS
mvnw.cmd clean compile  # Windows
```

### 4. 运行项目

#### 4.1 开发模式运行

```bash
# 使用Maven运行
mvn spring-boot:run

# 或者使用Maven Wrapper
./mvnw spring-boot:run  # Linux/macOS
mvnw.cmd spring-boot:run  # Windows
```

#### 4.2 打包运行

```bash
# 打包项目
mvn clean package

# 运行JAR文件
java -jar target/VoiceMind-Java-core-0.0.1-SNAPSHOT.jar
```

#### 4.3 IDE运行

1. 在IDE中打开项目
2. 找到 `VoiceMindJavaCoreApplication.java` 主类
3. 右键选择 "Run" 或 "Debug"

### 5. 验证运行状态

#### 5.1 检查服务状态

访问健康检查接口：
```
GET http://localhost:8081/api/health
```

#### 5.2 查看API文档

访问Swagger UI文档：
```
http://localhost:8081/api/swagger-ui.html
```

## 配置说明

### 1. 应用配置

#### 1.1 基础配置 (`application.yml`)

```yaml
spring:
  application:
    name: VoiceMind-Java-core
  profiles:
    active: local  # 激活的配置文件
  servlet:
    multipart:
      max-file-size: 10MB      # 最大文件上传大小
      max-request-size: 10MB   # 最大请求大小

server:
  port: 8081                   # 服务端口
  servlet:
    context-path: /api         # 上下文路径

# API文档配置
springdoc:
  swagger-ui:
    path: /swagger-ui.html
  api-docs:
    path: /v3/api-docs
  group-configs:
    - group: 'default'
      paths-to-match: '/**'
      packages-to-scan: com.csjzsn.voicemindjavacore.controller

# Knife4j增强配置
knife4j:
  enable: true
  setting:
    language: zh_cn
```

#### 1.2 本地配置 (`application-local.yml`)

```yaml
spring:
  ai:
    dashscope:
      api-key: your-api-key
      chat:
        options:
          model: qwen3-max
gemini:
  api-key: your-gemini-key
  base-url: https://poloai.top/v1
  model: gemini-2.5-pro-thinking
search-api:
  api-key: your-search-key
```

### 2. 环境变量配置

可以通过环境变量覆盖配置文件中的设置：

```bash
# Windows
set SPRING_AI_DASHSCOPE_API_KEY=your-api-key
set SERVER_PORT=8081

# Linux/macOS
export SPRING_AI_DASHSCOPE_API_KEY=your-api-key
export SERVER_PORT=8081
```

## API接口使用

### 1. 基础AI聊天

```bash
curl -X POST "http://localhost:8081/api/ai/chat/do" \
  -H "Content-Type: application/json" \
  -d '{
    "text": "你好，请介绍一下自己",
    "session_id": "chat_001"
  }'
```

### 2. 带报告的AI聊天

```bash
curl -X POST "http://localhost:8081/api/ai/chat/report" \
  -H "Content-Type: application/json" \
  -d '{
    "text": "帮我分析一下这个订单",
    "session_id": "chat_002"
  }'
```

### 3. RAG增强聊天

```bash
curl -X POST "http://localhost:8081/api/ai/chat/rag" \
  -H "Content-Type: application/json" \
  -d '{
    "text": "查询相关文档信息",
    "session_id": "chat_003"
  }'
```

### 4. 工具增强聊天

```bash
curl -X POST "http://localhost:8081/api/ai/chat/tools" \
  -H "Content-Type: application/json" \
  -d '{
    "text": "帮我搜索一下最新的AI技术",
    "session_id": "chat_004"
  }'
```

### 5. SSE流式聊天

```bash
curl -X GET "http://localhost:8081/api/ai/chat/sse?message=你好&chatId=chat_005"
```

## 常见问题

### 1. 启动失败

**问题**: 应用启动失败，提示端口被占用
**解决方案**:
```bash
# 检查端口占用
netstat -ano | findstr :8081  # Windows
lsof -i :8081                 # Linux/macOS

# 修改端口配置
# 在application.yml中修改server.port
```

**问题**: API密钥配置错误
**解决方案**:
- 检查 `application-local.yml` 中的API密钥配置
- 确认API密钥有效且有足够额度
- 检查网络连接

### 2. 依赖问题

**问题**: Maven依赖下载失败
**解决方案**:
```bash
# 清理Maven缓存
mvn dependency:purge-local-repository

# 重新下载依赖
mvn clean install

# 使用阿里云镜像（在pom.xml中已配置）
```

### 3. 内存不足

**问题**: 运行时出现OutOfMemoryError
**解决方案**:
```bash
# 增加JVM内存
java -Xms512m -Xmx2048m -jar target/VoiceMind-Java-core-0.0.1-SNAPSHOT.jar
```

### 4. 工具调用失败

**问题**: AI工具调用失败
**解决方案**:
- 检查工具权限配置
- 确认系统环境支持（如终端操作需要相应权限）
- 查看日志文件获取详细错误信息

## 日志配置

### 1. 日志级别配置

在 `application.yml` 中添加：

```yaml
logging:
  level:
    com.csjzsn.voicemindjavacore: DEBUG
    org.springframework.ai: INFO
    root: INFO
```

### 2. 日志文件配置

```yaml
logging:
  file:
    name: logs/voicemind.log
  pattern:
    file: "%d{yyyy-MM-dd HH:mm:ss} [%thread] %-5level %logger{36} - %msg%n"
```

## 性能调优

### 1. JVM参数调优

```bash
# 生产环境推荐参数
java -Xms1g -Xmx4g -XX:+UseG1GC -XX:MaxGCPauseMillis=200 \
  -jar target/VoiceMind-Java-core-0.0.1-SNAPSHOT.jar
```

### 2. 应用参数调优

```yaml
# 在application.yml中配置
spring:
  ai:
    dashscope:
      chat:
        options:
          temperature: 0.7
          max-tokens: 2000
```

## 部署建议

### 1. 开发环境
- 使用IDE直接运行
- 开启热重载功能
- 使用内存存储

### 2. 测试环境
- 使用Docker容器部署
- 配置外部数据库
- 启用详细日志

### 3. 生产环境
- 使用Docker或Kubernetes部署
- 配置负载均衡
- 使用外部存储
- 启用监控和告警

## 监控和维护

### 1. 健康检查

定期检查以下端点：
- `/api/health` - 应用健康状态
- `/api/swagger-ui.html` - API文档

### 2. 日志监控

关注以下日志：
- 错误日志
- 性能日志
- 安全日志

### 3. 资源监控

监控以下指标：
- CPU使用率
- 内存使用率
- 磁盘空间
- 网络连接

## 联系支持

如遇到问题，请：
1. 查看本文档的常见问题部分
2. 检查项目日志文件
3. 提交Issue到项目仓库
4. 联系开发团队

---

**注意**: 请确保在生产环境中妥善保管API密钥，不要将包含敏感信息的配置文件提交到版本控制系统。

