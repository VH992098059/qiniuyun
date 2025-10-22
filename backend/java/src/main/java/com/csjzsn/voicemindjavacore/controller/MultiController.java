package com.csjzsn.voicemindjavacore.controller;

import com.csjzsn.voicemindjavacore.app.ImageApp;
import com.csjzsn.voicemindjavacore.common.BaseResponse;
import jakarta.annotation.Resource;
import org.springframework.http.MediaType;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RequestParam;
import org.springframework.web.bind.annotation.RestController;
import org.springframework.web.multipart.MultipartFile;

@RestController
@RequestMapping("/ai/multi")
public class MultiController {

    @Resource
    private ImageApp imageApp;

    /**
     * 处理上传的图片并识别播放按钮位置
     */
    @PostMapping(value = "/analyze", consumes = MediaType.MULTIPART_FORM_DATA_VALUE)
    public BaseResponse<?> analyzeImage(
            @RequestParam("file") MultipartFile file,
            @RequestParam(value = "prompt", required = false) String prompt) throws Exception {
        String doMulti = imageApp.doMulti(file, prompt);
        return BaseResponse.success(doMulti);

    }
}
