package com.csjzsn.voicemindjavacore.common;

import com.csjzsn.voicemindjavacore.exceptions.CustRuntimeException;
import io.swagger.v3.oas.annotations.media.Schema;
import lombok.Data;

@Data
@Schema(description = "统一响应格式")
public class BaseResponse<T> {
    
    @Schema(description = "响应状态码", example = "200", required = true)
    private Integer code;
    
    @Schema(description = "响应消息", example = "success", required = true)
    private String message;
    
    @Schema(description = "响应数据", example = "具体的业务数据")
    private T data;
    
    @Schema(description = "响应时间戳", example = "1703123456789", required = true)
    private Long timestamp;

    public BaseResponse() {
        this.timestamp = System.currentTimeMillis();
    }

    /**
     * 创建成功响应
     * @param data 响应数据
     * @return 成功响应对象
     */
    public static <T> BaseResponse<T> success(T data) {
        BaseResponse<T> response = new BaseResponse<>();
        response.code = 200;
        response.message = "success";
        response.data = data;
        return response;
    }

    /**
     * 创建成功响应（无数据）
     * @return 成功响应对象
     */
    public static BaseResponse<?> success() {
        return success(null);
    }

    /**
     * 创建失败响应
     * @param code 错误码
     * @param message 错误消息
     * @return 失败响应对象
     */
    public static BaseResponse<?> fail(Integer code, String message) {
        BaseResponse<Object> response = new BaseResponse<>();
        response.code = code;
        response.message = message;
        return response;
    }

    /**
     * 根据业务异常创建失败响应
     * @param e 业务异常
     * @return 失败响应对象
     */
    public static BaseResponse<?> fail(CustRuntimeException e) {
        return fail(e.getCode(), e.getMessage());
    }
}