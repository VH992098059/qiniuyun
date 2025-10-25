package logic

import (
	"context"
	"encoding/json"
	"fmt"
	common "golang/utility"
	"log"
	"os"
	"path/filepath"
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
	log.Println(launch.SafetyCheck.Intention)
	var appName string
	start := strings.Index(string(jsonStr), `"`)
	end := strings.Index(string(jsonStr[start+1:]), `"`) + start + 1
	if start != -1 && end != -1 {
		appName = string(jsonStr[start+1 : end])
	}

	// --- 任务1: 启动一个应用程序 ---
	// 你可以修改 "spotify" 为你想测试的应用, 比如 "notepad" (Windows) 或 "calculator" (macOS)
	log.Println("\n[TASK 1] Launching an application (" + appName + ")...")
	err = LaunchApplication(launch.Data, appName)
	if err != nil {
		log.Println(err)
		return
	} else {
		log.Println("\nWaiting for the app to open")
		appNameAction := appName + ".exe"
		absPhotoSave, _ := filepath.Abs("files/photos/screenshot.png")
		for {
			if isAppActive(appNameAction) {
				time.Sleep(900 * time.Millisecond)
				capture := robotgo.CaptureScreen()
				defer robotgo.FreeBitmap(capture)
				img := robotgo.ToImage(capture)

				log.Println("等待中")
				imgo.Save(absPhotoSave, img)
				log.Println("截图成功")
				break
			}
			time.Sleep(800 * time.Millisecond)
			fmt.Println("暂时未启动")
		}
		//如果需要操作
		if launch.SafetyCheck.Intention == "action" {
			log.Println("等待中")
			time.Sleep(500 * time.Millisecond)
			//OCR解析
			jsonPath := common.OcrLogic(ctx, absPhotoSave)
			//json分解
			jsonResultPath := common.NewJsonFile(jsonPath)
			abs, _ := filepath.Abs("files/output/new_orc.png")
			//坐标嵌入当前截图
			resultJsonPic := common.SynthesisPhoto(jsonResultPath, abs)
			open, err := os.Open(resultJsonPic)
			if err != nil {
				log.Println("文件打开失败", err)
				return
			}
			actionModelReq := &ActionModelReq{
				FileImage: open,
				Text:      text,
			}
			log.Println("等待中")
			time.Sleep(500 * time.Millisecond)
			model, err := ActionModel(ctx, actionModelReq)
			if err != nil {
				log.Println("模型处理错误：", err)
				return
			}
			log.Println(model.Actions)
		}
	}

}
func isAppActive(targetProcess string) bool {
	// 获取当前所有进程 ID
	pids, _ := robotgo.Pids()
	for _, pid := range pids {
		pname, err := robotgo.FindName(pid)
		if err != nil {
			continue
		}
		if strings.EqualFold(pname, targetProcess) {
			fmt.Println("找到目标应用：", pname, "PID:", pid)
			time.Sleep(800 * time.Millisecond)
			// 激活窗口（置顶）
			success := robotgo.ActivePid(pid)
			if success == nil {
				fmt.Println("已将应用置于最前面")
			} else {
				fmt.Println("激活失败")
			}
			return true
		}
	}
	fmt.Println("未找到目标应用：", targetProcess)
	return false

}
func controlApplication(action ActionModelInTaskData) {

}
