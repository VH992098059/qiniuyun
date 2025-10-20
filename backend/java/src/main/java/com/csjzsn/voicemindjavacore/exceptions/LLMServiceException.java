package com.csjzsn.voicemindjavacore.exceptions;

/**
 * LLM服务调用异常（如OpenAI API调用失败）
 */
public class LLMServiceException extends CustRuntimeException {
    public LLMServiceException(String message) {
        super(1001, "LLM服务错误: " + message);
    }

    public LLMServiceException(String message, Throwable cause) {
        super(1001, "LLM服务错误: " + message, cause);
    }
}