package com.csjzsn.voicemindjavacore.model.dto;

import io.swagger.v3.oas.annotations.media.Schema;
import jakarta.validation.constraints.NotBlank;
import lombok.Data;

/**
 *  AI聊天请求体
 */
@Data
@Schema(description = "AI聊天请求参数")
public class ProcessRequest {
    
    @NotBlank(message = "session_id不能为空")
    @Schema(description = "会话ID", example = "session_123456", required = true)
    private String session_id;

    @NotBlank(message = "消息不能为空")
    @Schema(description = "用户输入的消息内容", example = "你好，请帮我分析一下这个报告", required = true)
    private String text;

    @NotBlank(message = "系统类型不能为空")
    @Schema(description = "系统类型", example = "windows", required = true)
    private String systemType;
}
