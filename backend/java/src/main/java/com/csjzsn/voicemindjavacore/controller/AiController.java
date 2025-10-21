package com.csjzsn.voicemindjavacore.controller;

import com.csjzsn.voicemindjavacore.common.BaseResponse;
import com.csjzsn.voicemindjavacore.common.OrderReport;
import com.csjzsn.voicemindjavacore.app.OrderApp;
import com.csjzsn.voicemindjavacore.model.dto.ProcessRequest;
import jakarta.annotation.Resource;
import jakarta.validation.Valid;
import org.springframework.ai.chat.model.ChatModel;
import org.springframework.ai.tool.ToolCallback;
import org.springframework.http.MediaType;
import org.springframework.web.bind.annotation.*;
import org.springframework.web.servlet.mvc.method.annotation.SseEmitter;
import reactor.core.publisher.Flux;

@RestController
@RequestMapping("/ai")
public class AiController {

    @Resource
    private OrderApp orderApp;


    @Resource
    private ToolCallback[] allTools;

    @Resource
    private ChatModel dashscopeChatModel;

    @PostMapping("/do")
    public BaseResponse<?> chat(@Valid @RequestBody ProcessRequest processRequest){
        String result = orderApp.doChat(processRequest.getText(), processRequest.getSession_id());
        return BaseResponse.success(result);
    }

    @PostMapping("/report")
    public BaseResponse<?> chatWithReport(@Valid @RequestBody ProcessRequest processRequest){
        OrderReport report = orderApp.doChatWithReport(processRequest.getText(), processRequest.getSession_id());
        return BaseResponse.success(report);
    }

    @PostMapping("/rag")
    public BaseResponse<?> chatWithRag(@Valid @RequestBody ProcessRequest processRequest){
        String result = orderApp.doChatWithRag(processRequest.getText(), processRequest.getSession_id());
        return BaseResponse.success(result);
    }

    @PostMapping("/tools")
    public BaseResponse<?> chatWithTools(@Valid @RequestBody ProcessRequest processRequest){
        String result = orderApp.doChatWithTools(processRequest.getText(), processRequest.getSession_id());
        return BaseResponse.success(result);
    }

    @PostMapping("/mcp")
    public BaseResponse<?> chatWithMcp(@Valid @RequestBody ProcessRequest processRequest){
        String result = orderApp.doChatWithMcp(processRequest.getText(), processRequest.getSession_id());
        return BaseResponse.success(result);
    }

    @GetMapping(value = "/sse", produces = MediaType.TEXT_EVENT_STREAM_VALUE)
    public Flux<String> doChatWithOrderAppSSE(String message, String chatId) {
        return orderApp.doChatByStream(message, chatId);
    }


    
}
