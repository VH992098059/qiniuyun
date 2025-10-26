package com.csjzsn.voicemindjavacore.app;


import com.csjzsn.voicemindjavacore.advisor.MyLoggerAdvisor;
import com.csjzsn.voicemindjavacore.advisor.ReReadingAdvisor;
import com.csjzsn.voicemindjavacore.common.OrderReport;
import com.csjzsn.voicemindjavacore.rag.QueryRewriter;
import jakarta.annotation.Resource;
import lombok.extern.slf4j.Slf4j;
import org.springframework.ai.chat.client.ChatClient;
import org.springframework.ai.chat.client.advisor.MessageChatMemoryAdvisor;
import org.springframework.ai.chat.client.advisor.vectorstore.QuestionAnswerAdvisor;
import org.springframework.ai.chat.memory.ChatMemory;
import org.springframework.ai.chat.memory.InMemoryChatMemoryRepository;
import org.springframework.ai.chat.memory.MessageWindowChatMemory;
import org.springframework.ai.chat.model.ChatModel;
import org.springframework.ai.chat.model.ChatResponse;
import org.springframework.ai.tool.ToolCallback;
import org.springframework.ai.tool.ToolCallbackProvider;
import org.springframework.ai.vectorstore.VectorStore;
import org.springframework.stereotype.Component;
import reactor.core.publisher.Flux;

@Component
@Slf4j
public class OrderApp {

    private final ChatClient chatClient;

    private static final String SYSTEM_PROMPT = """
            你是一个智能API接口设计助手，负责处理用户语音指令并生成标准化的JSON响应。请严格遵守以下规范:
            将用户自然语言指令解析为结构化操作意图
            生成前后端通用的标准化响应格式
            确保所有操作在安全许可范围内
            对模糊指令进行合理推断或要求澄清
            需要你判断用户操作内容，如果只是启动应用则为“launch”，如果需要操作则为“action”
            以下是响应格式规范示例
            {"safetyCheck":{"level":"yellow","message":"将在桌面创建新文件","intention":"action"},"data":{"Notepad":["记事本","notepad.exe"]}}
            """;

    /**
     * 初始化 ChatClient
     *
     * @param dashscopeChatModel 模型
     */
    public OrderApp(ChatModel dashscopeChatModel) {
        // 初始化基于内存的对话记忆
        MessageWindowChatMemory chatMemory = MessageWindowChatMemory.builder()
                .chatMemoryRepository(new InMemoryChatMemoryRepository())
                .maxMessages(20)
                .build();
        chatClient = ChatClient.builder(dashscopeChatModel)
                .defaultSystem(SYSTEM_PROMPT)
                .defaultAdvisors(
                        MessageChatMemoryAdvisor.builder(chatMemory).build(),
                        // 自定义日志 Advisor，可按需开启
                        new MyLoggerAdvisor()
                        // 自定义推理增强 Advisor，可按需开启
                       ,new ReReadingAdvisor()
                )
                .build();
    }


    /**
     * AI 基础对话（支持多轮对话记忆）
     *
     * @param message 消息
     * @param chatId  会话id
     * @return  字符串
     */
    public String doChat(String message, String chatId) {
        ChatResponse chatResponse = chatClient
                .prompt()
                .user(message)
                .advisors(spec -> spec.param(ChatMemory.CONVERSATION_ID, chatId))
                .call()
                .chatResponse();
        String content = chatResponse.getResult().getOutput().getText();
        log.info("content: {}", content);
        return content;
    }

    /**
     * AI 基础对话（支持多轮对话记忆，SSE 流式传输）
     *
     * @param message   消息
     * @param chatId    会话ID
     * @return  Flux
     */
    public Flux<String> doChatByStream(String message, String chatId) {
        return chatClient
                .prompt()
                .user(message)
                .advisors(spec -> spec.param(ChatMemory.CONVERSATION_ID, chatId))
                .stream()
                .content();
    }

    /**
     * AI 指令报告功能（实战结构化输出）
     *
     * @param message   消息
     * @param chatId    会话ID
     * @return  OrderReport
     */
    public OrderReport doChatWithReport(String message, String chatId) {
        OrderReport OrderReport = chatClient
                .prompt()
                .system(SYSTEM_PROMPT + "每次对话后都要生成指令结果，标题为{标题}的指令报告，内容为建议列表")
                .user(message)
                .advisors(spec -> spec.param(ChatMemory.CONVERSATION_ID, chatId))
                .call()
                .entity(OrderReport.class);
        log.info("OrderReport: {}", OrderReport);
        return OrderReport;
    }

    @Resource
    private VectorStore OrderAppVectorStore;

    @Resource
    private QueryRewriter queryRewriter;

    /**
     *  RAG 知识库进行对话
     *
     * @param message     消息
     * @param chatId    会话ID
     * @return  String
     */
    public String doChatWithRag(String message, String chatId) {
        // 查询重写
        String rewrittenMessage = queryRewriter.doQueryRewrite(message);
        ChatResponse chatResponse = chatClient
                .prompt()
                // 使用改写后的查询
                .user(rewrittenMessage)
                .advisors(spec -> spec.param(ChatMemory.CONVERSATION_ID, chatId))
                // 开启日志，便于观察效果
                .advisors(new MyLoggerAdvisor())
                // 应用 RAG 知识库问答
                .advisors(new QuestionAnswerAdvisor(OrderAppVectorStore))
                // 应用 RAG 检索增强服务（基于云知识库服务）
//                .advisors(OrderAppRagCloudAdvisor)
                // 应用 RAG 检索增强服务（基于 PgVector 向量存储）
//                .advisors(new QuestionAnswerAdvisor(pgVectorVectorStore))
                // 应用自定义的 RAG 检索增强服务（文档查询器 + 上下文增强器）
//                .advisors(
//                        OrderAppRagCustomAdvisorFactory.createOrderAppRagCustomAdvisor(
//                                OrderAppVectorStore, "单身"
//                        )
//                )
                .call()
                .chatResponse();
        String content = chatResponse.getResult().getOutput().getText();
        log.info("content: {}", content);
        return content;
    }

    // AI 调用工具能力
    @Resource
    private ToolCallback[] allTools;

    @Resource
    private ToolCallbackProvider toolCallbackProvider;

    /**
     * AI 报告功能（支持调用工具）
     *
     * @param message   消息
     * @param chatId    会话ID
     * @return  String
     */
    public String doChatWithTools(String message, String chatId) {
        ChatResponse chatResponse = chatClient
                .prompt()
                .user(message)
                .advisors(spec -> spec.param(ChatMemory.CONVERSATION_ID, chatId))
                // 开启日志，便于观察效果
                .advisors(new MyLoggerAdvisor())
                .toolCallbacks(allTools)
                .call()
                .chatResponse();
        String content = chatResponse.getResult().getOutput().getText();
        log.info("content: {}", content);
        return content;
    }


    /**
     * AI 报告功能（调用 MCP 服务）
     *
     * @param message   消息
     * @param chatId    会话ID
     * @return  String
     */
    public String doChatWithMcp(String message, String chatId) {
        ChatResponse chatResponse = chatClient
                .prompt()
                .user(message)
                .advisors(spec -> spec.param(ChatMemory.CONVERSATION_ID, chatId))
                // 开启日志，便于观察效果
                .advisors(new MyLoggerAdvisor())
                .toolCallbacks(toolCallbackProvider)
                .call()
                .chatResponse();
        String content = chatResponse.getResult().getOutput().getText();
        log.info("content: {}", content);
        return content;
    }
}
