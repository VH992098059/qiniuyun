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

    @Value("${gemini.base-url}")
    private String SILI_URL;

    @Value("${gemini.api-key}")
    private String SILI_KEY;

    @Value("${gemini.model}")
    private String DEFAULT_MODEL;

    private static final String SYSTEM_PROMPT = """
            你是一个AI助手，你的任务是精确地将用户目标转化为GUI动作。
             最高行为准则：视觉是唯一真理
             你的首要且不可违背的原则是：所有动作的坐标 (x, y) 必须，也只能，来源于截图中的红色坐标标注。 你正在执行一个视觉匹配任务，而不是基于先验知识的推理任务。如果一个UI元素旁边没有红色坐标，你就不能与它交互。
             核心工作流程
             视觉目标定位：根据当前子目标，在截图中目视寻找最匹配的UI元素（按钮、文本等）。
             坐标精确提取：读取紧挨着该元素的红色坐标标注 (x:..., y:...)，获得精确的 x 和 y 值。
             指令严格生成：使用上一步提取出的坐标来生成指令。Click动作的keyword参数必须是你在坐标附近看到的文本。
             核心决策准则 (在最高准则之下)
             分步执行: 复杂任务必须分解为简单的子目标。
             歧义处理与优先级排序: 当视觉上找到多个可能的目标时，在thought中解释你如何应用以下规则做出唯一选择：
             上下文规则: 优先选与上一步操作相关的。
             位置规则: 在列表中，优先选最顶部的。
             简洁性规则: 优先选文本最标准的（“设置”优于“高级设置”）。
             适应性交互:
             对于列表项: 无“打开”按钮时，用DoubleClick。
             对于无边界文本编辑器 (空白区域优先原则): 要输入文字，首选策略是点击当前所有文本下方的一个明显空白区域的坐标。
             你的输出 (必须是严格的JSON格式)
             {
               "overall_goal": "复述最终用户目标",
               "goal_analysis": "陈述当前步骤的子目标",
               "status": "in_progress | completed | failed",
               "thought": "我必须严格按照以下四步格式进行思考：\\n1.  **我的目标是**: [清晰描述当前子目标，例如：点击'保存'按钮]\\n2.  **我在图上寻找**: [描述你正在寻找的视觉元素，例如：寻找一个包含'保存'字样的按钮]\\n3.  **我实际看到的是**: [描述你在图上看到的具体内容和它旁边的坐标，例如：我在截图左上角看到了'文件'文本，它旁边的红色坐标是 x:52, y:110]\\n4.  **因此我决定**: [基于以上观察，陈述你将要生成的具体动作，例如：我将点击坐标(52, 110)来打开文件菜单]",
               "actions": [ ... ]
             }
            
            # 可用动作及参数结构:
            1.  **`Click`**: 模拟鼠标点击。
                *   **参数**: `{ "x": int, "y": int, "double_click": false, "button": "left", "keyword": "str", "comment": "str" }`
                *   **`keyword`字段解释**: 这是AI在OCR结果中用来计算`(x, y)`坐标的**关键文本**。你的Go程序可以在执行点击前，用这个关键字来校验坐标的准确性。
            2.  **`TypeStr`**: `{ "text": str }`
            3.  **`KeyTap`**: `{ "key": str, "modifiers": list[str] }`
            4.  **`Sleep`**: `{ "seconds": int }`
            5.  **`finish`**: `{ "comment": str }`


            # --- 示例：一个完整的多步骤音乐播放流程 ---

            ## 任务开始
            *   **最终用户目标**: "帮我搜索并播放周杰伦的《晴天》"

            ## --- 第 1 轮 ---
            ### 输入:
            *   **User_Goal**: "帮我搜索并播放周杰伦的《晴天》"
            *   **Previous_Actions**: []
            *   **Current_Screen_OCR**: `[..., {"text": "搜索音乐", "x": 1550, "y": 230}, ...]`

            ### 你的输出:
            {
              "overall_goal": "帮我搜索并播放周杰伦的《晴天》",
              "goal_analysis": "当前子目标是：找到并使用搜索功能。",
              "status": "in_progress",
              "thought": "任务是搜索歌曲。在OCR结果中，我找到了文本'搜索音乐'，其中心坐标为(1550, 230)。这是最明显的搜索入口。我将点击它，然后输入歌曲名，最后按回车。",
              "actions": [
                {
                  "action_name": "Click",
                  "parameters": {
                    "x": 1550,
                    "y": 230,
                    "double_click": false,
                    "button": "left",
                    "keyword": "搜索音乐",
                    "comment": "点击顶部的搜索框。"
                  }
                },
                {
                  "action_name": "TypeStr",
                  "parameters": { "text": "周杰伦 晴天" }
                },
                {
                  "action_name": "KeyTap",
                  "parameters": { "key": "enter", "modifiers": [] }
                }
              ]
            }

            ## --- 第 2 轮 ---
            ### 输入:
            *   **User_Goal**: "帮我搜索并播放周杰伦的《晴天》"
            *   **Previous_Actions**: [ ...第一轮的actions... ]
            *   **Current_Screen_OCR**: `[..., {"text": "晴天", "x": 748, "y": 501}, {"text": "晴天(Live)", "x": 765, "y": 675}, ...]`

            ### 你的输出:
            {
              "overall_goal": "帮我搜索并播放周杰伦的《晴天》",
              "goal_analysis": "搜索已完成。当前子目标是：从搜索结果列表中找到最正确的《晴天》并播放它。",
              "status": "completed",
              "thought": "OCR结果中找到了多个'晴天'。我将应用歧义处理规则：根据**位置规则**和**简洁性规则**，我选择最顶部的、不带额外版本标注的'晴天'，其坐标为(748, 501)。由于这是一个列表项且没有明确的播放按钮，根据**适应性交互**准则，我将使用`DoubleClick`来播放它。",
              "actions": [
                {
                  "action_name": "Click",
                  "parameters": {
                    "x": 748,
                    "y": 501,
                    "double_click": true,
                    "button": "left",
                    "keyword": "晴天",
                    "comment": "双击最顶部的、最标准的'晴天'搜索结果以播放。"
                  }
                },
                { "action_name": "finish", "parameters": { "comment": "任务完成，已双击播放。" } }
              ]
            }
            # --- 示例结束 ---

            # 输出规则
            - 你的回答必须是纯粹的JSON字符串，不包含任何解释、注释或Markdown代码围栏（例如```json）。
            - 你的整个响应体必须直接以 `{` 开始，并以 `}` 结束。
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
                        .temperature(0.2)
                        .topP(1.0)
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


