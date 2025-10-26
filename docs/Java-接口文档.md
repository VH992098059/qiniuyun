# VoiceMind Java Core API


**简介**:VoiceMind Java Core API


**HOST**:http://localhost:8081/api


**联系人**:CSJZSN


**Version**:1.0.0


**接口路径**:/api/v3/api-docs/default


[TOC]






# AI聊天接口


## 工具增强聊天


**接口地址**:`/api/ai/chat/tools`


**请求方式**:`POST`


**请求数据类型**:`application/x-www-form-urlencoded,application/json`


**响应数据类型**:`*/*`


**接口描述**:<p>使用工具集成的AI聊天，可以调用各种工具来完成任务</p>



**请求示例**:


```javascript
{
  "session_id": "session_123456",
  "text": "你好，请帮我分析一下这个报告",
  "systemType": "windows"
}
```


**请求参数**:


| 参数名称 | 参数说明 | 请求类型    | 是否必须 | 数据类型 | schema |
| -------- | -------- | ----- | -------- | -------- | ------ |
|processRequest|AI聊天请求参数|body|true|ProcessRequest|ProcessRequest|
|&emsp;&emsp;session_id|会话ID||true|string||
|&emsp;&emsp;text|用户输入的消息内容||true|string||
|&emsp;&emsp;systemType|系统类型||true|string||


**响应状态**:


| 状态码 | 说明 | schema |
| -------- | -------- | ----- | 
|200|工具聊天成功|BaseResponse|
|400|请求参数错误|BaseResponseObject|
|500|服务器内部错误|BaseResponseObject|


**响应状态码-200**:


**响应参数**:


| 参数名称 | 参数说明 | 类型 | schema |
| -------- | -------- | ----- |----- | 
|code|响应状态码|integer(int32)|integer(int32)|
|message|响应消息|string||
|data|响应数据|object||
|timestamp|响应时间戳|integer(int64)|integer(int64)|


**响应示例**:
```javascript
{
	"code": 200,
	"message": "success",
	"data": "具体的业务数据",
	"timestamp": 1703123456789
}
```


**响应状态码-400**:


**响应参数**:


| 参数名称 | 参数说明 | 类型 | schema |
| -------- | -------- | ----- |----- | 
|code|响应状态码|integer(int32)|integer(int32)|
|message|响应消息|string||
|data|响应数据|object||
|timestamp|响应时间戳|integer(int64)|integer(int64)|


**响应示例**:
```javascript
{
	"code": 200,
	"message": "success",
	"data": "具体的业务数据",
	"timestamp": 1703123456789
}
```


**响应状态码-500**:


**响应参数**:


| 参数名称 | 参数说明 | 类型 | schema |
| -------- | -------- | ----- |----- | 
|code|响应状态码|integer(int32)|integer(int32)|
|message|响应消息|string||
|data|响应数据|object||
|timestamp|响应时间戳|integer(int64)|integer(int64)|


**响应示例**:
```javascript
{
	"code": 200,
	"message": "success",
	"data": "具体的业务数据",
	"timestamp": 1703123456789
}
```


## 带报告的AI聊天


**接口地址**:`/api/ai/chat/report`


**请求方式**:`POST`


**请求数据类型**:`application/x-www-form-urlencoded,application/json`


**响应数据类型**:`*/*`


**接口描述**:<p>使用AI模型进行聊天并返回详细的分析报告</p>



**请求示例**:


```javascript
{
  "session_id": "session_123456",
  "text": "你好，请帮我分析一下这个报告",
  "systemType": "windows"
}
```


**请求参数**:


| 参数名称 | 参数说明 | 请求类型    | 是否必须 | 数据类型 | schema |
| -------- | -------- | ----- | -------- | -------- | ------ |
|processRequest|AI聊天请求参数|body|true|ProcessRequest|ProcessRequest|
|&emsp;&emsp;session_id|会话ID||true|string||
|&emsp;&emsp;text|用户输入的消息内容||true|string||
|&emsp;&emsp;systemType|系统类型||true|string||


**响应状态**:


| 状态码 | 说明 | schema |
| -------- | -------- | ----- | 
|200|聊天成功，返回分析报告|BaseResponse|
|400|请求参数错误|BaseResponseObject|
|500|服务器内部错误|BaseResponseObject|


**响应状态码-200**:


**响应参数**:


| 参数名称 | 参数说明 | 类型 | schema |
| -------- | -------- | ----- |----- | 
|code|响应状态码|integer(int32)|integer(int32)|
|message|响应消息|string||
|data|响应数据|object||
|timestamp|响应时间戳|integer(int64)|integer(int64)|


**响应示例**:
```javascript
{
	"code": 200,
	"message": "success",
	"data": "具体的业务数据",
	"timestamp": 1703123456789
}
```


**响应状态码-400**:


**响应参数**:


| 参数名称 | 参数说明 | 类型 | schema |
| -------- | -------- | ----- |----- | 
|code|响应状态码|integer(int32)|integer(int32)|
|message|响应消息|string||
|data|响应数据|object||
|timestamp|响应时间戳|integer(int64)|integer(int64)|


**响应示例**:
```javascript
{
	"code": 200,
	"message": "success",
	"data": "具体的业务数据",
	"timestamp": 1703123456789
}
```


**响应状态码-500**:


**响应参数**:


| 参数名称 | 参数说明 | 类型 | schema |
| -------- | -------- | ----- |----- | 
|code|响应状态码|integer(int32)|integer(int32)|
|message|响应消息|string||
|data|响应数据|object||
|timestamp|响应时间戳|integer(int64)|integer(int64)|


**响应示例**:
```javascript
{
	"code": 200,
	"message": "success",
	"data": "具体的业务数据",
	"timestamp": 1703123456789
}
```


## RAG增强聊天


**接口地址**:`/api/ai/chat/rag`


**请求方式**:`POST`


**请求数据类型**:`application/x-www-form-urlencoded,application/json`


**响应数据类型**:`*/*`


**接口描述**:<p>使用RAG（检索增强生成）技术进行聊天，提供更准确的上下文信息</p>



**请求示例**:


```javascript
{
  "session_id": "session_123456",
  "text": "你好，请帮我分析一下这个报告",
  "systemType": "windows"
}
```


**请求参数**:


| 参数名称 | 参数说明 | 请求类型    | 是否必须 | 数据类型 | schema |
| -------- | -------- | ----- | -------- | -------- | ------ |
|processRequest|AI聊天请求参数|body|true|ProcessRequest|ProcessRequest|
|&emsp;&emsp;session_id|会话ID||true|string||
|&emsp;&emsp;text|用户输入的消息内容||true|string||
|&emsp;&emsp;systemType|系统类型||true|string||


**响应状态**:


| 状态码 | 说明 | schema |
| -------- | -------- | ----- | 
|200|RAG聊天成功|BaseResponse|
|400|请求参数错误|BaseResponseObject|
|500|服务器内部错误|BaseResponseObject|


**响应状态码-200**:


**响应参数**:


| 参数名称 | 参数说明 | 类型 | schema |
| -------- | -------- | ----- |----- | 
|code|响应状态码|integer(int32)|integer(int32)|
|message|响应消息|string||
|data|响应数据|object||
|timestamp|响应时间戳|integer(int64)|integer(int64)|


**响应示例**:
```javascript
{
	"code": 200,
	"message": "success",
	"data": "具体的业务数据",
	"timestamp": 1703123456789
}
```


**响应状态码-400**:


**响应参数**:


| 参数名称 | 参数说明 | 类型 | schema |
| -------- | -------- | ----- |----- | 
|code|响应状态码|integer(int32)|integer(int32)|
|message|响应消息|string||
|data|响应数据|object||
|timestamp|响应时间戳|integer(int64)|integer(int64)|


**响应示例**:
```javascript
{
	"code": 200,
	"message": "success",
	"data": "具体的业务数据",
	"timestamp": 1703123456789
}
```


**响应状态码-500**:


**响应参数**:


| 参数名称 | 参数说明 | 类型 | schema |
| -------- | -------- | ----- |----- | 
|code|响应状态码|integer(int32)|integer(int32)|
|message|响应消息|string||
|data|响应数据|object||
|timestamp|响应时间戳|integer(int64)|integer(int64)|


**响应示例**:
```javascript
{
	"code": 200,
	"message": "success",
	"data": "具体的业务数据",
	"timestamp": 1703123456789
}
```


## MCP协议聊天


**接口地址**:`/api/ai/chat/mcp`


**请求方式**:`POST`


**请求数据类型**:`application/x-www-form-urlencoded,application/json`


**响应数据类型**:`*/*`


**接口描述**:<p>使用MCP（Model Context Protocol）协议进行聊天</p>



**请求示例**:


```javascript
{
  "session_id": "session_123456",
  "text": "你好，请帮我分析一下这个报告",
  "systemType": "windows"
}
```


**请求参数**:


| 参数名称 | 参数说明 | 请求类型    | 是否必须 | 数据类型 | schema |
| -------- | -------- | ----- | -------- | -------- | ------ |
|processRequest|AI聊天请求参数|body|true|ProcessRequest|ProcessRequest|
|&emsp;&emsp;session_id|会话ID||true|string||
|&emsp;&emsp;text|用户输入的消息内容||true|string||
|&emsp;&emsp;systemType|系统类型||true|string||


**响应状态**:


| 状态码 | 说明 | schema |
| -------- | -------- | ----- | 
|200|MCP聊天成功|BaseResponse|
|400|请求参数错误|BaseResponseObject|
|500|服务器内部错误|BaseResponseObject|


**响应状态码-200**:


**响应参数**:


| 参数名称 | 参数说明 | 类型 | schema |
| -------- | -------- | ----- |----- | 
|code|响应状态码|integer(int32)|integer(int32)|
|message|响应消息|string||
|data|响应数据|object||
|timestamp|响应时间戳|integer(int64)|integer(int64)|


**响应示例**:
```javascript
{
	"code": 200,
	"message": "success",
	"data": "具体的业务数据",
	"timestamp": 1703123456789
}
```


**响应状态码-400**:


**响应参数**:


| 参数名称 | 参数说明 | 类型 | schema |
| -------- | -------- | ----- |----- | 
|code|响应状态码|integer(int32)|integer(int32)|
|message|响应消息|string||
|data|响应数据|object||
|timestamp|响应时间戳|integer(int64)|integer(int64)|


**响应示例**:
```javascript
{
	"code": 200,
	"message": "success",
	"data": "具体的业务数据",
	"timestamp": 1703123456789
}
```


**响应状态码-500**:


**响应参数**:


| 参数名称 | 参数说明 | 类型 | schema |
| -------- | -------- | ----- |----- | 
|code|响应状态码|integer(int32)|integer(int32)|
|message|响应消息|string||
|data|响应数据|object||
|timestamp|响应时间戳|integer(int64)|integer(int64)|


**响应示例**:
```javascript
{
	"code": 200,
	"message": "success",
	"data": "具体的业务数据",
	"timestamp": 1703123456789
}
```


## 基础AI聊天


**接口地址**:`/api/ai/chat/do`


**请求方式**:`POST`


**请求数据类型**:`application/x-www-form-urlencoded,application/json`


**响应数据类型**:`*/*`


**接口描述**:<p>使用基础AI模型进行聊天对话</p>



**请求示例**:


```javascript
{
  "session_id": "session_123456",
  "text": "你好，请帮我分析一下这个报告",
  "systemType": "windows"
}
```


**请求参数**:


| 参数名称 | 参数说明 | 请求类型    | 是否必须 | 数据类型 | schema |
| -------- | -------- | ----- | -------- | -------- | ------ |
|processRequest|AI聊天请求参数|body|true|ProcessRequest|ProcessRequest|
|&emsp;&emsp;session_id|会话ID||true|string||
|&emsp;&emsp;text|用户输入的消息内容||true|string||
|&emsp;&emsp;systemType|系统类型||true|string||


**响应状态**:


| 状态码 | 说明 | schema |
| -------- | -------- | ----- | 
|200|聊天成功|BaseResponse|
|400|请求参数错误|BaseResponseObject|
|500|服务器内部错误|BaseResponseObject|


**响应状态码-200**:


**响应参数**:


| 参数名称 | 参数说明 | 类型 | schema |
| -------- | -------- | ----- |----- | 
|code|响应状态码|integer(int32)|integer(int32)|
|message|响应消息|string||
|data|响应数据|object||
|timestamp|响应时间戳|integer(int64)|integer(int64)|


**响应示例**:
```javascript
{
	"code": 200,
	"message": "success",
	"data": "具体的业务数据",
	"timestamp": 1703123456789
}
```


**响应状态码-400**:


**响应参数**:


| 参数名称 | 参数说明 | 类型 | schema |
| -------- | -------- | ----- |----- | 
|code|响应状态码|integer(int32)|integer(int32)|
|message|响应消息|string||
|data|响应数据|object||
|timestamp|响应时间戳|integer(int64)|integer(int64)|


**响应示例**:
```javascript
{
	"code": 200,
	"message": "success",
	"data": "具体的业务数据",
	"timestamp": 1703123456789
}
```


**响应状态码-500**:


**响应参数**:


| 参数名称 | 参数说明 | 类型 | schema |
| -------- | -------- | ----- |----- | 
|code|响应状态码|integer(int32)|integer(int32)|
|message|响应消息|string||
|data|响应数据|object||
|timestamp|响应时间戳|integer(int64)|integer(int64)|


**响应示例**:
```javascript
{
	"code": 200,
	"message": "success",
	"data": "具体的业务数据",
	"timestamp": 1703123456789
}
```


## SSE流式聊天


**接口地址**:`/api/ai/chat/sse`


**请求方式**:`GET`


**请求数据类型**:`application/x-www-form-urlencoded`


**响应数据类型**:`text/event-stream`


**接口描述**:<p>使用Server-Sent Events进行流式聊天，实时返回AI响应</p>



**请求参数**:


| 参数名称 | 参数说明 | 请求类型    | 是否必须 | 数据类型 | schema |
| -------- | -------- | ----- | -------- | -------- | ------ |
|message|聊天消息内容|query|true|string||
|chatId|聊天会话ID|query|true|string||


**响应状态**:


| 状态码 | 说明 | schema |
| -------- | -------- | ----- | 
|200|流式聊天成功||
|400|请求参数错误||
|500|服务器内部错误||


**响应参数**:


暂无


**响应示例**:
```javascript

```


# 健康检查接口


## 健康检查


**接口地址**:`/api/health`


**请求方式**:`GET`


**请求数据类型**:`application/x-www-form-urlencoded`


**响应数据类型**:`*/*`


**接口描述**:<p>检查系统是否正常运行</p>



**请求参数**:


暂无


**响应状态**:


| 状态码 | 说明 | schema |
| -------- | -------- | ----- | 
|200|系统正常运行||
|500|系统异常||


**响应参数**:


暂无


**响应示例**:
```javascript

```


# 多模态AI接口


## 图像分析


**接口地址**:`/api/ai/multi/analyze`


**请求方式**:`POST`


**请求数据类型**:`application/x-www-form-urlencoded,multipart/form-data`


**响应数据类型**:`*/*`


**接口描述**:<p>上传图片并使用AI进行图像分析，识别播放按钮位置等元素</p>



**请求参数**:


| 参数名称 | 参数说明 | 请求类型    | 是否必须 | 数据类型 | schema |
| -------- | -------- | ----- | -------- | -------- | ------ |
|file|上传的图片文件|query|true|file||
|prompt|可选的提示词，用于指导AI分析|query|false|string||


**响应状态**:


| 状态码 | 说明 | schema |
| -------- | -------- | ----- | 
|200|图像分析成功|BaseResponse|
|400|请求参数错误或文件格式不支持|BaseResponseObject|
|500|服务器内部错误|BaseResponseObject|


**响应状态码-200**:


**响应参数**:


| 参数名称 | 参数说明 | 类型 | schema |
| -------- | -------- | ----- |----- | 
|code|响应状态码|integer(int32)|integer(int32)|
|message|响应消息|string||
|data|响应数据|object||
|timestamp|响应时间戳|integer(int64)|integer(int64)|


**响应示例**:
```javascript
{
	"code": 200,
	"message": "success",
	"data": "具体的业务数据",
	"timestamp": 1703123456789
}
```


**响应状态码-400**:


**响应参数**:


| 参数名称 | 参数说明 | 类型 | schema |
| -------- | -------- | ----- |----- | 
|code|响应状态码|integer(int32)|integer(int32)|
|message|响应消息|string||
|data|响应数据|object||
|timestamp|响应时间戳|integer(int64)|integer(int64)|


**响应示例**:
```javascript
{
	"code": 200,
	"message": "success",
	"data": "具体的业务数据",
	"timestamp": 1703123456789
}
```


**响应状态码-500**:


**响应参数**:


| 参数名称 | 参数说明 | 类型 | schema |
| -------- | -------- | ----- |----- | 
|code|响应状态码|integer(int32)|integer(int32)|
|message|响应消息|string||
|data|响应数据|object||
|timestamp|响应时间戳|integer(int64)|integer(int64)|


**响应示例**:
```javascript
{
	"code": 200,
	"message": "success",
	"data": "具体的业务数据",
	"timestamp": 1703123456789
}
```