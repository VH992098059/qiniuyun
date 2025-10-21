package com.csjzsn.voicemindjavacore.common;

import lombok.Getter;
import lombok.Setter;
import org.springframework.stereotype.Component;

import java.util.List;

@Component
@Setter
@Getter
public class OrderReport {
    // Getters and Setters
    private String apiVersion;
    private String status;
    private CommandSequence commandSequence;
    private List<Command> commands;
    private GlobalSafetyCheck globalSafetyCheck;

    @Setter
    @Getter
    public static class CommandSequence {
        // Getters and Setters
        private int total;
        private int current;
        private boolean canParallel;

    }

    @Setter
    @Getter
    public static class Command {
        // Getters and Setters
        private Data data;
        private Execution execution;

    }

    @Setter
    @Getter
    public static class Data {
        // Getters and Setters
        private String actionType;
        private Parameters parameters;
        private SafetyCheck safetyCheck;

    }

    @Setter
    @Getter
    public static class Parameters {
        // Getters and Setters
        private String path;
        private String content;
        private String appName;
        private String filePath;

    }

    @Setter
    @Getter
    public static class SafetyCheck {
        // Getters and Setters
        private String level;

    }

    @Setter
    @Getter
    public static class Execution {
        // Getters and Setters
        private Backend backend;
        private Frontend frontend;

    }

    @Setter
    @Getter
    public static class Backend {
        // Getters and Setters
        private String golang;

    }

    @Setter
    @Getter
    public static class Frontend {
        // Getters and Setters
        private String displayText;
        private String voiceResponse;

    }

    @Setter
    @Getter
    public static class GlobalSafetyCheck {
        // Getters and Setters
        private String level;
        private List<Dependency> dependencies;

    }

    @Setter
    @Getter
    public static class Dependency {
        // Getters and Setters
        private int commandIndex;
        private String requiredStatus;

    }
}