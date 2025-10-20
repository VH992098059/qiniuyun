package com.csjzsn.voicemindjavacore.exceptions;

/**
 * 请求参数校验失败异常
 */
public class InvalidRequestException extends CustRuntimeException {
    public InvalidRequestException(String message) {
        super(1002, "请求参数错误: " + message);
    }
}