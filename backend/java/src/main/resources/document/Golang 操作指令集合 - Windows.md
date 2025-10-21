# Golang 操作指令集合 - Windows 系统

#### 在桌面创建周报文件

{"apiVersion":"v1","status":"success","data":{"actionType":"file","parameters":{"path":"~/Desktop/周报.txt","action":"create"},"safetyCheck":{"level":"yellow","message":"将在桌面创建新文件"}},"execution":{"backend":{"golang":"package main\n\nimport (\n\t\"os\"\n\t\"path/filepath\"\n)\n\nfunc main() {\n\tpath := filepath.Join(os.Getenv(\"HOME\"), \"Desktop\", \"周报.txt\")\n\tfile, _ := os.Create(path)\n\tdefer file.Close()\n}"},"frontend":{"displayText":"确认在桌面创建 周报.txt 文件？","voiceResponse":"您确认要在桌面创建周报文件吗？请说确认或取消"}}}

#### 播放周杰伦的七里香

{"apiVersion":"v1","status":"success","data":{"actionType":"music","parameters":{"artist":"周杰伦","track":"七里香"},"safetyCheck":{"level":"green","message":""}},"execution":{"backend":{"golang":"package main\n\nimport (\n\t\"os/exec\"\n)\n\nfunc main() {\n\tcmd := exec.Command(\"open\", \"spotify:search:七里香\")\n\tcmd.Run()\n}"},"frontend":{"displayText":"即将播放：周杰伦 - 七里香","voiceResponse":"正在为您播放周杰伦的七里香"}}}