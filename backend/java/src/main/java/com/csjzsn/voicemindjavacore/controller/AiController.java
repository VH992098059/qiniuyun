package com.csjzsn.voicemindjavacore.controller;

import com.csjzsn.voicemindjavacore.app.OrderApp;
import com.csjzsn.voicemindjavacore.common.BaseResponse;
import com.csjzsn.voicemindjavacore.common.OrderReport;
import com.csjzsn.voicemindjavacore.model.dto.ProcessRequest;
import io.swagger.v3.oas.annotations.Operation;
import io.swagger.v3.oas.annotations.Parameter;
import io.swagger.v3.oas.annotations.media.Content;
import io.swagger.v3.oas.annotations.media.Schema;
import io.swagger.v3.oas.annotations.responses.ApiResponse;
import io.swagger.v3.oas.annotations.responses.ApiResponses;
import io.swagger.v3.oas.annotations.tags.Tag;
import jakarta.annotation.Resource;
import jakarta.validation.Valid;
import org.springframework.ai.chat.model.ChatModel;
import org.springframework.ai.tool.ToolCallback;
import org.springframework.http.MediaType;
import org.springframework.web.bind.annotation.*;
import reactor.core.publisher.Flux;

@RestController
@RequestMapping("/ai/chat")
@Tag(name = "AI聊天接口", description = "提供各种AI聊天功能的接口")
public class AiController {

    @Resource
    private OrderApp orderApp;


    @Resource
    private ToolCallback[] allTools;

    @Resource
    private ChatModel dashscopeChatModel;

    @Operation(summary = "基础AI聊天", description = "使用基础AI模型进行聊天对话")
    @ApiResponses(value = {
            @ApiResponse(responseCode = "200", description = "聊天成功", 
                    content = @Content(schema = @Schema(implementation = BaseResponse.class))),
            @ApiResponse(responseCode = "400", description = "请求参数错误"),
            @ApiResponse(responseCode = "500", description = "服务器内部错误")
    })
    @PostMapping("/do")
    public BaseResponse<?> chat(@Valid @RequestBody ProcessRequest processRequest){
        String result = orderApp.doChat(processRequest.getText(), processRequest.getSession_id());
        return BaseResponse.success(result);
    }

    @Operation(summary = "带报告的AI聊天", description = "使用AI模型进行聊天并返回详细的分析报告")
    @ApiResponses(value = {
            @ApiResponse(responseCode = "200", description = "聊天成功，返回分析报告", 
                    content = @Content(schema = @Schema(implementation = BaseResponse.class))),
            @ApiResponse(responseCode = "400", description = "请求参数错误"),
            @ApiResponse(responseCode = "500", description = "服务器内部错误")
    })
    @PostMapping("/report")
    public BaseResponse<?> chatWithReport(@Valid @RequestBody ProcessRequest processRequest){
        OrderReport report = orderApp.doChatWithReport(processRequest.getText(), processRequest.getSession_id());
        return BaseResponse.success(report);
    }

    @Operation(summary = "RAG增强聊天", description = "使用RAG（检索增强生成）技术进行聊天，提供更准确的上下文信息")
    @ApiResponses(value = {
            @ApiResponse(responseCode = "200", description = "RAG聊天成功", 
                    content = @Content(schema = @Schema(implementation = BaseResponse.class))),
            @ApiResponse(responseCode = "400", description = "请求参数错误"),
            @ApiResponse(responseCode = "500", description = "服务器内部错误")
    })
    @PostMapping("/rag")
    public BaseResponse<?> chatWithRag(@Valid @RequestBody ProcessRequest processRequest){
        String result = orderApp.doChatWithRag(processRequest.getText(), processRequest.getSession_id());
        return BaseResponse.success(result);
    }

    @Operation(summary = "工具增强聊天", description = "使用工具集成的AI聊天，可以调用各种工具来完成任务")
    @ApiResponses(value = {
            @ApiResponse(responseCode = "200", description = "工具聊天成功", 
                    content = @Content(schema = @Schema(implementation = BaseResponse.class))),
            @ApiResponse(responseCode = "400", description = "请求参数错误"),
            @ApiResponse(responseCode = "500", description = "服务器内部错误")
    })
    @PostMapping("/tools")
    public BaseResponse<?> chatWithTools(@Valid @RequestBody ProcessRequest processRequest){
        String result = orderApp.doChatWithTools(processRequest.getText(), processRequest.getSession_id());
        return BaseResponse.success(result);
    }

    @Operation(summary = "MCP协议聊天", description = "使用MCP（Model Context Protocol）协议进行聊天")
    @ApiResponses(value = {
            @ApiResponse(responseCode = "200", description = "MCP聊天成功", 
                    content = @Content(schema = @Schema(implementation = BaseResponse.class))),
            @ApiResponse(responseCode = "400", description = "请求参数错误"),
            @ApiResponse(responseCode = "500", description = "服务器内部错误")
    })
    @PostMapping("/mcp")
    public BaseResponse<?> chatWithMcp(@Valid @RequestBody ProcessRequest processRequest){
        String result = orderApp.doChatWithMcp(processRequest.getText(), processRequest.getSession_id());
        return BaseResponse.success(result);
    }

    @Operation(summary = "SSE流式聊天", description = "使用Server-Sent Events进行流式聊天，实时返回AI响应")
    @ApiResponses(value = {
            @ApiResponse(responseCode = "200", description = "流式聊天成功", 
                    content = @Content(mediaType = "text/event-stream")),
            @ApiResponse(responseCode = "400", description = "请求参数错误"),
            @ApiResponse(responseCode = "500", description = "服务器内部错误")
    })
    @GetMapping(value = "/sse", produces = MediaType.TEXT_EVENT_STREAM_VALUE)
    public Flux<String> doChatWithOrderAppSSE(
            @Parameter(description = "聊天消息内容", example = "你好，请帮我分析一下这个订单") 
            @RequestParam String message, 
            @Parameter(description = "聊天会话ID", example = "chat_123456") 
            @RequestParam String chatId) {
        return orderApp.doChatByStream(message, chatId);
    }


    
}
