package com.csjzsn.voicemindjavacore.model.dto;

import lombok.Data;

import java.util.List;

/**
 * 响应体
 */
@Data
public class ProcessResponse {
    private String session_id;
    private String response;
    private List<Action> actions;
}