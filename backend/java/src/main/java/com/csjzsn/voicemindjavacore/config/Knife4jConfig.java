package com.csjzsn.voicemindjavacore.config;


import io.swagger.v3.oas.models.OpenAPI;
import io.swagger.v3.oas.models.info.Contact;
import io.swagger.v3.oas.models.info.Info;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;

@Configuration
public class Knife4jConfig {
    @Bean
    public OpenAPI customOpenAPI() {
        return new OpenAPI()
                .info(new Info()
                        .title("API")
                        .version("1.0")
                        .description("PersonalBlog-API")
                        .termsOfService("https://test.com")
                        .contact(new Contact().name("CSJZSN").url("https://test.com").email("CSJZSN@163.com"))
                );
    }
}