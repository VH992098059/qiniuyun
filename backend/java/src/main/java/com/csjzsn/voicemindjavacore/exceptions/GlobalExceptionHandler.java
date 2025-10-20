package com.csjzsn.voicemindjavacore.exceptions;// 文件名: GlobalExceptionHandler.java

import com.csjzsn.voicemindjavacore.common.BaseResponse;
import jakarta.servlet.http.HttpServletRequest;
import lombok.extern.slf4j.Slf4j;
import org.springframework.http.HttpStatus;
import org.springframework.web.bind.MethodArgumentNotValidException;
import org.springframework.web.bind.annotation.ExceptionHandler;
import org.springframework.web.bind.annotation.ResponseStatus;
import org.springframework.web.bind.annotation.RestControllerAdvice;

@Slf4j
@RestControllerAdvice
public class GlobalExceptionHandler {

    // 处理自定义业务异常
    @ExceptionHandler(CustRuntimeException.class)
    public BaseResponse<?> handleAppException(CustRuntimeException e, HttpServletRequest request) {
        log.warn("业务异常: URL={}, Code={}, Msg={}",
                request.getRequestURI(), e.getCode(), e.getMessage());
        return BaseResponse.fail(e);
    }

    // 处理参数校验异常
    @ExceptionHandler(MethodArgumentNotValidException.class)
    public BaseResponse<?> handleValidationException(MethodArgumentNotValidException e) {
        String errorMsg = e.getBindingResult().getFieldErrors().stream()
                .map(error -> error.getField() + ": " + error.getDefaultMessage())
                .findFirst()
                .orElse("参数校验失败");
        return BaseResponse.fail(1002, errorMsg);
    }

    // 处理其他所有未捕获异常
    @ExceptionHandler(Exception.class)
    @ResponseStatus(HttpStatus.INTERNAL_SERVER_ERROR)
    public BaseResponse<?> handleException(Exception e, HttpServletRequest request) {
        log.error("系统异常: URL={}", request.getRequestURI(), e);
        return BaseResponse.fail(500, "系统繁忙，请稍后再试");
    }
}