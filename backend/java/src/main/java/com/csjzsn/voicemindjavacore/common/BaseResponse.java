package com.csjzsn.voicemindjavacore.common;

import com.csjzsn.voicemindjavacore.exceptions.CustRuntimeException;
import lombok.Data;

@Data
public class BaseResponse<T> {
    private Integer code;
    private String message;
    private T data;
    private Long timestamp;

    public BaseResponse() {
        this.timestamp = System.currentTimeMillis();
    }

    // 成功响应
    public static <T> BaseResponse<T> success(T data) {
        BaseResponse<T> response = new BaseResponse<>();
        response.code = 200;
        response.message = "success";
        response.data = data;
        return response;
    }

    public static BaseResponse<?> success() {
        return success(null);
    }

    // 失败响应
    public static BaseResponse<?> fail(Integer code, String message) {
        BaseResponse<Object> response = new BaseResponse<>();
        response.code = code;
        response.message = message;
        return response;
    }

    // 快速创建业务异常响应
    public static BaseResponse<?> fail(CustRuntimeException e) {
        return fail(e.getCode(), e.getMessage());
    }
}