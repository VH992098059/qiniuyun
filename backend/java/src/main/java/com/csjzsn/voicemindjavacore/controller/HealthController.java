package com.csjzsn.voicemindjavacore.controller;

import io.swagger.v3.oas.annotations.Operation;
import io.swagger.v3.oas.annotations.responses.ApiResponse;
import io.swagger.v3.oas.annotations.responses.ApiResponses;
import io.swagger.v3.oas.annotations.tags.Tag;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

@RestController
@RequestMapping("/health")
@Tag(name = "健康检查接口", description = "提供系统健康状态检查功能")
public class HealthController {

    @Operation(summary = "健康检查", description = "检查系统是否正常运行")
    @ApiResponses(value = {
            @ApiResponse(responseCode = "200", description = "系统正常运行", 
                    content = @io.swagger.v3.oas.annotations.media.Content(schema = @io.swagger.v3.oas.annotations.media.Schema(example = "ok"))),
            @ApiResponse(responseCode = "500", description = "系统异常")
    })
    @GetMapping
    public String healthCheck(){
        return "ok";
    }
}
