package com.csjzsn.voicemindjavacore.app;


import com.alibaba.cloud.ai.dashscope.chat.DashScopeChatOptions;
import com.alibaba.cloud.ai.dashscope.chat.MessageFormat;
import com.alibaba.cloud.ai.dashscope.common.DashScopeApiConstants;
import com.csjzsn.voicemindjavacore.advisor.MyLoggerAdvisor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.ai.chat.client.ChatClient;
import org.springframework.ai.chat.client.advisor.MessageChatMemoryAdvisor;
import org.springframework.ai.chat.memory.InMemoryChatMemoryRepository;
import org.springframework.ai.chat.memory.MessageWindowChatMemory;
import org.springframework.ai.chat.messages.UserMessage;
import org.springframework.ai.chat.model.ChatModel;
import org.springframework.ai.chat.model.ChatResponse;
import org.springframework.ai.chat.prompt.Prompt;
import org.springframework.ai.content.Media;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.stereotype.Component;
import org.springframework.util.MimeTypeUtils;
import org.springframework.web.multipart.MultipartFile;

import java.util.Map;

@Component
@Slf4j
public class ImageApp {

    private final ChatClient dashScopeChatClient;

    @Value("${embedding.options.model}")
    private String DEFAULT_MODEL;

    private static final String SYSTEM_PROMPT = """
            识别图像中的UI元素，按以下格式返回JSON示例：
               {"elements":[{"type":"button|input|menu","subtype":"play|close|search","x":100,"y":200,"w":40,"h":40,"text":""}],"img_size":[1920,1080]}
               识别优先级：
               1. 播放/暂停按钮（▶/❚❚）
               2. 关闭按钮（×）
               3. 搜索框（🔍或"搜索"文字）
               4. 菜单按钮（☰）
               规则：只返回置信度>0.7的元素，忽略尺寸<10x10的元素，无元素时返回 {"elements": []}
            """;

    /**
     *  初始化模型
     * @param chatModel ChatModel
     */
    public ImageApp(ChatModel chatModel){
        // 初始化基于内存的对话记忆
        MessageWindowChatMemory chatMemory = MessageWindowChatMemory.builder()
                .chatMemoryRepository(new InMemoryChatMemoryRepository())
                .maxMessages(20)
                .build();
        dashScopeChatClient=ChatClient.builder(chatModel)
                .defaultSystem(SYSTEM_PROMPT)
                .defaultAdvisors(
                        MessageChatMemoryAdvisor.builder(chatMemory).build(),
                        // 自定义日志 Advisor，可按需开启
                        new MyLoggerAdvisor()
                        // 自定义推理增强 Advisor，可按需开启
//                        ,new ReReadingAdvisor()
                )
                .build();

    }

    public String doMulti(MultipartFile file, String message) {
        Media media = new Media(MimeTypeUtils.parseMimeType(file.getContentType()), file.getResource());

        // 3. 构建多模态请求
        UserMessage userMessage = UserMessage.builder()
                .text(message)
                .media(media)
                .metadata(Map.of(
                        DashScopeApiConstants.MESSAGE_FORMAT, MessageFormat.IMAGE,
                        "response_format", "json_object" // 要求返回JSON
                ))
                .build();

        // 4. 调用DashScope API
        ChatResponse response = dashScopeChatClient
                .prompt(new Prompt(userMessage, DashScopeChatOptions.builder()
                        .withModel(DEFAULT_MODEL)
                        .withMultiModel(true)
                        .build()))
                .call()
                .chatResponse();

        // 5. 解析AI返回的JSON
        return response.getResult().getOutput().getText();
    }

}


