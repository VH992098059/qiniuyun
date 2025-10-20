package com.csjzsn.voicemindjavacore.exceptions;

/**
 * 动作执行异常（Go端执行失败时，Java端收到的异常）
 */
public class ActionExecutionException extends CustRuntimeException {
    public ActionExecutionException(String message) {
        super(1003, "动作执行失败: " + message);
    }
}