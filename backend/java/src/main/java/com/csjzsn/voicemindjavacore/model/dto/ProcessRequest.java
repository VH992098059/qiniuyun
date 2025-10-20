package com.csjzsn.voicemindjavacore.model.dto;

import jakarta.validation.constraints.NotBlank;
import lombok.Data;

/**
 *  请求体
 */
@Data
public class ProcessRequest {
    @NotBlank(message = "session_id不能为空")
    private String session_id;

    @NotBlank(message = "消息不能为空")
    private String text;

    @NotBlank(message = "系统类型不能为空")
    private String systemType;
}
