package com.csjzsn.voicemindjavacore.controller;

import com.csjzsn.voicemindjavacore.app.ImageApp;
import com.csjzsn.voicemindjavacore.common.BaseResponse;
import io.swagger.v3.oas.annotations.Operation;
import io.swagger.v3.oas.annotations.Parameter;
import io.swagger.v3.oas.annotations.media.Content;
import io.swagger.v3.oas.annotations.media.Schema;
import io.swagger.v3.oas.annotations.responses.ApiResponse;
import io.swagger.v3.oas.annotations.responses.ApiResponses;
import io.swagger.v3.oas.annotations.tags.Tag;
import jakarta.annotation.Resource;
import org.springframework.http.MediaType;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RequestParam;
import org.springframework.web.bind.annotation.RestController;
import org.springframework.web.multipart.MultipartFile;

@RestController
@RequestMapping("/ai/multi")
@Tag(name = "多模态AI接口", description = "提供图像分析和多模态AI处理功能")
public class MultiController {

    @Resource
    private ImageApp imageApp;

    @Operation(summary = "图像分析", description = "上传图片并使用AI进行图像分析，识别播放按钮位置等元素")
    @ApiResponses(value = {
            @ApiResponse(responseCode = "200", description = "图像分析成功", 
                    content = @Content(schema = @Schema(implementation = BaseResponse.class))),
            @ApiResponse(responseCode = "400", description = "请求参数错误或文件格式不支持"),
            @ApiResponse(responseCode = "500", description = "服务器内部错误")
    })
    @PostMapping(value = "/analyze", consumes = MediaType.MULTIPART_FORM_DATA_VALUE)
    public BaseResponse<?> analyzeImage(
            @Parameter(description = "上传的图片文件", required = true)
            @RequestParam("file") MultipartFile file,
            @Parameter(description = "可选的提示词，用于指导AI分析", example = "请识别图片中的播放按钮")
            @RequestParam(value = "prompt", required = false) String prompt) throws Exception {
        String doMulti = imageApp.doMulti(file, prompt);
        return BaseResponse.success(doMulti);
    }
}
