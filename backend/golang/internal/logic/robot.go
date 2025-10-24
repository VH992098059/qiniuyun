package logic

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"runtime"
	"strings"
	"time"

	"github.com/go-vgo/robotgo"
	"github.com/vcaesar/imgo"
)

func RobotAutoApplication(ctx context.Context, text string) {
	log.Println("--- [Desktop Control Program Started] ---")
	log.Printf("--- [OS Detected: %s] ---", runtime.GOOS)
	launch, err := ModelService(context.Background(), text)
	if err != nil {
		log.Println(err)
		return
	}
	jsonStr, err := json.Marshal(launch.Data)
	if err != nil {
		log.Println(err)
		return
	}
	var appName string
	start := strings.Index(string(jsonStr), `"`)
	end := strings.Index(string(jsonStr[start+1:]), `"`) + start + 1
	if start != -1 && end != -1 {
		appName = string(jsonStr[start+1 : end])
	}

	// --- 任务1: 启动一个应用程序 ---
	// 你可以修改 "spotify" 为你想测试的应用, 比如 "notepad" (Windows) 或 "calculator" (macOS)
	log.Println("\n[TASK 1] Launching an application (" + appName + ")...")
	LaunchApplication(launch.Data, appName)
	log.Println("\nWaiting for the app to open")
	appNameAction := appName + ".exe"
	for {
		if isAppActive(appNameAction) {
			capture := robotgo.CaptureScreen()
			defer robotgo.FreeBitmap(capture)
			img := robotgo.ToImage(capture)
			robotgo.Sleep(3)
			imgo.Save("screenshot.png", img)
			log.Println("截图成功")
			log.Println("等待中")
			break
		}
		time.Sleep(800 * time.Millisecond)
		fmt.Println("暂时未启动")
	}

	/*// --- 等待几秒钟，给应用启动的时间，也让你能切换到音乐播放器 ---
	log.Println("\nWaiting for 8 seconds to allow the app to open and start playing...")
	time.Sleep(8 * time.Second)

	// --- 任务2: 模拟媒体按键控制 ---
	log.Println("\n[TASK 2] Sending 'Play/Pause' command to the OS.")
	login.ExecuteMediaKey("toggle_play_pause")

	log.Println("Waiting for 4 seconds...")
	time.Sleep(4 * time.Second)

	log.Println("\n[TASK 3] Sending 'Next Track' command.")
	login.ExecuteMediaKey("next")

	log.Println("Waiting for 4 seconds...")
	time.Sleep(4 * time.Second)

	log.Println("\n[TASK 4] Sending 'Previous Track' command.")
	login.ExecuteMediaKey("prev")

	log.Println("\n--- [All tasks completed. Program will now exit.] ---")*/

}
func isAppActive(targetProcess string) bool {
	pid := robotgo.GetPid()
	pname, _ := robotgo.FindName(pid)
	fmt.Println("应用名称：", pname)
	return strings.EqualFold(pname, targetProcess)
}
