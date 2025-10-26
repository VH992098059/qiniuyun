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
	appNameAction, err := LaunchApplication(launch.Data, appName)
	if err != nil {
		log.Println(err)
		return
	} else {
		log.Println("\nWaiting for the app to open")
		absPhotoSave, _ := filepath.Abs("files/photos/screenshot.png")
		temp := 0
		for {
			if isAppActive(appNameAction) && temp != 10 {
				screen(absPhotoSave)
				break
			} else if temp > 10 {
				log.Println("启动失败")
				return
			}
			time.Sleep(800 * time.Millisecond)
			fmt.Println("暂时未启动")
			temp++
		}
		//如果需要操作
		if launch.SafetyCheck.Intention == "action" {
			var status string
			for {
				log.Println("等待中")
				time.Sleep(500 * time.Millisecond)
				//OCR处理
				resultJsonPic := OCR(ctx, absPhotoSave)
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
				log.Println("正在分析意图")
				model, err := ActionModel(ctx, actionModelReq)
				if err != nil {
					log.Println("模型处理错误：", err)
					return
				}
				status = model.Status
				controlApplication(model.Actions)
				if status == "completed" {
					log.Println("已完成")
					break
				}
				screen(absPhotoSave)
			}
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
func controlApplication(action []taskDataActions) {
	for i, act := range action {
		switch strings.ToLower(act.ActionName) {
		case "click":
			var p struct {
				X           int    `json:"x"`
				Y           int    `json:"y"`
				DoubleClick bool   `json:"double_click"`
				Button      string `json:"button"`
				Keyword     string `json:"keyword"`
				Comment     string `json:"comment"`
			}
			if err := decodeParamsTo(act.Parameters, &p); err != nil {
				log.Println("解析 Click 失败：", err)
				continue
			}
			btn := strings.ToLower(strings.TrimSpace(p.Button))
			if btn == "" {
				btn = "left"
			}
			robotgo.MoveSmooth(p.X, p.Y, 0.9, 0.9)
			time.Sleep(500 * time.Millisecond)
			robotgo.Click(btn, p.DoubleClick)

		case "typestr":
			var p struct {
				Text string `json:"text"`
			}
			if err := decodeParamsTo(act.Parameters, &p); err != nil {
				log.Println("解析 TypeStr 失败：", err)
				continue
			}
			time.Sleep(500 * time.Millisecond)
			robotgo.TypeStr(p.Text)

		case "keytap":
			var p struct {
				Key       string   `json:"key"`
				Modifiers []string `json:"modifiers"`
			}
			if err := decodeParamsTo(act.Parameters, &p); err != nil {
				log.Println("解析 KeyTap 失败：", err)
				continue
			}
			if len(p.Modifiers) > 0 {
				mods := make([]interface{}, len(p.Modifiers))
				for i := range p.Modifiers {
					mods[i] = p.Modifiers[i]
				}
				time.Sleep(500 * time.Millisecond)
				robotgo.KeyTap(p.Key, mods...)
			} else {
				time.Sleep(500 * time.Millisecond)
				robotgo.KeyTap(p.Key)
			}

		case "finish":
			var p struct {
				Comment string `json:"comment"`
			}
			time.Sleep(500 * time.Millisecond)
			_ = decodeParamsTo(act.Parameters, &p)
			log.Printf("任务完成：%s", p.Comment)

		default:
			log.Printf("未识别的动作 '%s'（第 %d 个）", act.ActionName, i)
		}
		time.Sleep(180 * time.Millisecond)
	}
}

// decodeParamsTo 将通用 map 解析为目标结构体
func decodeParamsTo(params map[string]any, out interface{}) error {
	b, err := json.Marshal(params)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, out)
}
func OCR(ctx context.Context, absPhotoSave string) (resultJsonPic string) {
	//OCR解析
	jsonPath := common.OcrLogic(ctx, absPhotoSave)
	//json分解
	jsonResultPath := common.NewJsonFile(jsonPath)
	abs, _ := filepath.Abs("files/output/new_orc.png")
	//坐标嵌入当前截图
	resultJsonPic = common.SynthesisPhoto(jsonResultPath, abs)
	return
}
func screen(absPhotoSave string) {
	time.Sleep(1 * time.Second)
	capture := robotgo.CaptureScreen()
	defer robotgo.FreeBitmap(capture)
	img := robotgo.ToImage(capture)
	log.Println("等待中")
	imgo.Save(absPhotoSave, img)
	log.Println("截图成功:", absPhotoSave)
	log.Println("准备进行下一步")
}
