package com.csjzsn.voicemindjavacore.exceptions;

import lombok.Getter;

/**
 * 应用自定义异常的基类
 */
@Getter
public class CustRuntimeException extends RuntimeException {
    private final Integer code;

    public CustRuntimeException(Integer code, String message) {
        super(message);
        this.code = code;
    }

    public CustRuntimeException(Integer code, String message, Throwable cause) {
        super(message, cause);
        this.code = code;
    }

}