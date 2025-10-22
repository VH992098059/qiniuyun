package com.csjzsn.voicemindjavacore.model.dto;

import lombok.Data;

/**
 * 用于响应中的动作指令
 */
@Data
public class Action {
    private String type; // exec, write_text, etc.
    private String command;
    private String content; // 用于write_text等内容
    private String path;    // 用于文件操作路径
}