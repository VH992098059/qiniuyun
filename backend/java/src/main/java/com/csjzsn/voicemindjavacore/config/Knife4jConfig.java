package com.csjzsn.voicemindjavacore.config;

import io.swagger.v3.oas.models.OpenAPI;
import io.swagger.v3.oas.models.info.Contact;
import io.swagger.v3.oas.models.info.Info;
import io.swagger.v3.oas.models.info.License;
import io.swagger.v3.oas.models.servers.Server;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;

import java.util.List;

@Configuration
public class Knife4jConfig {
    
    @Bean
    public OpenAPI customOpenAPI() {
        return new OpenAPI()
                .info(new Info()
                        .title("VoiceMind Java Core API")
                        .version("1.0.0")
                        .description("基于Spring AI的智能语音助手后端API，提供多种AI聊天模式、图像分析、工具集成等功能")
                        .termsOfService("https://github.com/VH992098059/qiniuyun")
                        .contact(new Contact()
                                .name("CSJZSN")
                                .url("https://github.com/csjzsn")
                                .email("CSJZSN@163.com"))
                        .license(new License()
                                .name("MIT License")
                                .url("https://opensource.org/licenses/MIT")));
    }
}