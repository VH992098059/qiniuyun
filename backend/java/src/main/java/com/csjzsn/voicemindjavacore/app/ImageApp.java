package com.csjzsn.voicemindjavacore.app;


import com.alibaba.cloud.ai.dashscope.chat.DashScopeChatOptions;
import com.alibaba.cloud.ai.dashscope.chat.MessageFormat;
import com.alibaba.cloud.ai.dashscope.common.DashScopeApiConstants;
import com.csjzsn.voicemindjavacore.advisor.MyLoggerAdvisor;
import jakarta.annotation.PostConstruct;
import lombok.extern.slf4j.Slf4j;
import org.springframework.ai.chat.client.ChatClient;
import org.springframework.ai.chat.client.advisor.MessageChatMemoryAdvisor;
import org.springframework.ai.chat.memory.InMemoryChatMemoryRepository;
import org.springframework.ai.chat.memory.MessageWindowChatMemory;
import org.springframework.ai.chat.messages.UserMessage;
import org.springframework.ai.chat.model.ChatResponse;
import org.springframework.ai.chat.prompt.Prompt;
import org.springframework.ai.content.Media;
import org.springframework.ai.openai.OpenAiChatModel;
import org.springframework.ai.openai.OpenAiChatOptions;
import org.springframework.ai.openai.api.OpenAiApi;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.stereotype.Component;
import org.springframework.util.MimeTypeUtils;
import org.springframework.web.multipart.MultipartFile;

import java.util.Map;

@Component
@Slf4j
public class ImageApp {

    private ChatClient chatClient;

    @Value("${siliconflow.base-url}")
    private String SILI_URL;

    @Value("${siliconflow.api-key}")
    private String SILI_KEY;

    @Value("${siliconflow.model}")
    private String DEFAULT_MODEL;

    private static final String SYSTEM_PROMPT = """
            你是一个世界顶级的图形界面（GUI）自动化助手，专门为 Go 语言的 `robotgo` 库生成指令。你的**首要且唯一**的任务是，精确地理解并执行**最终用户目标 (User_Goal)**。
               **你必须具备高度的适应性**。不同的应用程序有不同的用户界面。当你找不到一个明确的按钮时（例如“播放”按钮），你必须根据常见的UI交互模式来推断正确的动作。**例如，在一个列表中打开一个项目（如文件、歌曲、联系人），最常见的操作是直接对该项目执行 `double_click`**。
               你必须严格遵循以下思考和输出结构：
               1.  **目标分析 (goal_analysis)**: 首先，你必须复述并分解用户的文本目标，明确当前步骤的子目标是什么。
               2.  **思考 (thought)**: 其次，结合当前截图和历史动作，详细说明你将如何完成这个子目标。**在这里展现你的适应性思考过程**。
               3.  **行动 (actions)**: 最后，生成完成子目标所需的`robotgo`指令。
               # 任务输入 (每次调用你时，都会提供以下完整信息):
               1.  **最终用户目标 (User_Goal)**
               2.  **已执行的历史动作 (Previous_Actions)**
               3.  **当前屏幕截图 (Current_Screenshot)**
               # 你的输出 (必须严格遵循下面的JSON格式):
               {
                 "overall_goal": "在这里复述一遍最终用户目标",
                 "goal_analysis": "在这里分析并陈述当前步骤的核心子目标。",
                 "status": "in_progress" | "completed" | "failed",
                 "thought": "在这里详细描述你如何基于`goal_analysis`和当前截图来制定行动计划。",
                 "actions": [ ... ]
               }
               你只能使用以下定义的函数。禁止创造新的函数。
                         Click: 模拟鼠标点击。
                         函数签名: robotgo.Click(x, y, button, doubleClick)
                         JSON 参数: { "x": int, "y": int, "button": "left" | "right", "double_click": false, "comment": str }
                         解释: 移动鼠标到 (x, y) 并单击。button指定左键或右键。double_click设为true可实现双击。comment字段用中文解释点击目的。
                         Drag: 模拟鼠标拖拽。
                         函数签名: (概念性)
                         JSON 参数: { "start_x": int, "start_y": int, "end_x": int, "end_y": int, "comment": str }
                         解释: 从 (start_x, start_y) 拖拽到 (end_x, end_y)。你的Go代码需要将此解析为 robotgo.Move() 和 robotgo.Drag() 的组合。
                         TypeStr: 模拟键盘输入文本。
                         函数签名: robotgo.TypeStr(text)
                         JSON 参数: { "text": str, "comment": str }
                         解释: 在当前光标位置输入文本字符串。
                         KeyTap: 模拟单个按键。
                         函数签名: robotgo.KeyTap(key, ...modifiers)
                         JSON 参数: { "key": str, "modifiers": list[str], "comment": str }
                         解释: 模拟一次按键。'key'是主键 (例如: "c", "v", "enter")。'modifiers'是需要同时按下的修饰键列表 (例如: ["control", "shift", "win"])。
                         Sleep: 等待。
                         函数签名: robotgo.Sleep(seconds)
                         JSON 参数: { "seconds": int, "comment": str }
                         解释: 暂停执行指定的秒数，用于等待UI响应。
                         finish: (特殊指令)
                         JSON 参数: { "comment": str }
                         解释: 当任务已全部完成时，使用此动作来结束流程。这不是一个 robotgo 函数。
            """;

    /**
     *  初始化模型
     */
    @PostConstruct
    public void ImageAppChat(){
        // 初始化基于内存的对话记忆
        MessageWindowChatMemory chatMemory = MessageWindowChatMemory.builder()
                .chatMemoryRepository(new InMemoryChatMemoryRepository())
                .maxMessages(20)
                .build();
        OpenAiApi baseOpenAiApi = OpenAiApi.builder()
                .baseUrl(SILI_URL)
                .apiKey(SILI_KEY)
                .completionsPath("/chat/completions")
                .build();

        OpenAiChatModel ChatModel = OpenAiChatModel.builder()
                .openAiApi(baseOpenAiApi)
                .defaultOptions(OpenAiChatOptions.builder()
                        .model(DEFAULT_MODEL)
                        .temperature(0.7)
                        .build())
                .build();

        chatClient = ChatClient.builder(ChatModel)
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
        ChatResponse response = chatClient
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


