package logic

import (
	"errors"
	"log"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/go-vgo/robotgo"
	"golang.org/x/sys/windows/registry"
)

func FindExecutablePathInRegistry(appDisplayName, exeName string) (string, error) {
	log.Printf("信息：正在搜索注册表以查找 '%s'", appDisplayName)
	keys := []registry.Key{
		registry.LOCAL_MACHINE,
		registry.CURRENT_USER,
	}
	paths := []string{
		`SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall`,
		`SOFTWARE\WOW6432Node\Microsoft\Windows\CurrentVersion\Uninstall`,
	}
	for _, key := range keys {
		for _, path := range paths {
			k, err := registry.OpenKey(key, path, registry.READ)
			if err != nil {
				continue
			}
			defer k.Close()
			//读取所有子健
			subKeyNames, err := k.ReadSubKeyNames(-1)
			if err != nil {
				continue
			}
			//遍历所有项目
			for _, subKeyName := range subKeyNames {
				subKey, err := registry.OpenKey(k, subKeyName, registry.READ)
				if err != nil {
					continue
				}
				defer subKey.Close()

				//读取 "DisplayName"值
				displayName, _, err := subKey.GetStringValue("DisplayName")
				if err != nil || displayName == "" {
					continue
				}

				//如果找到想要的程序
				if strings.Contains(displayName, appDisplayName) {
					//优先读取“InstallLocation”
					installLocation, _, err := subKey.GetStringValue("InstallLocation")
					if err == nil && installLocation != "" {
						fullPath := filepath.Join(installLocation, exeName)
						log.Printf("成功：通过 InstallLocation 找到了 %s 的路径：%s", appDisplayName, fullPath)
						return fullPath, nil
					}
					//如果没有“InstallLocation”,尝试从“DisplayIcon”中提取路径
					displayIcon, _, err := subKey.GetStringValue("DisplayIcon")
					if err == nil && displayIcon != "" {
						//移除末尾的逗号和图片索引
						if commIndex := strings.LastIndex(displayIcon, ","); commIndex != -1 {
							displayIcon = displayIcon[:commIndex]
						}
						//移除可能存在的引号
						displayIcon = strings.Trim(displayIcon, `"`)
						log.Printf("成功：通过 DisplayIcon 找到了 %s 的路径：%s", appDisplayName, displayIcon)
						return displayIcon, nil
					}
				}
			}
		}
	}
	return "", errors.New("注册表中未找到应用程序")
}

// LaunchApplication 启动应用程序
func LaunchApplication(appName string) error {
	log.Printf("INFO: Attempting to launch application: %s", appName)
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "windows":
		appMap := map[string][]string{
			"QQMusic": {"QQ音乐", "QQMusic.exe"},
		}
		if appInfo, ok := appMap[appName]; ok {
			var fullPath string
			var err error

			// 如果有注册表显示名称，就去搜索
			if appInfo[0] != "" {
				fullPath, err = FindExecutablePathInRegistry(appInfo[0], appInfo[1])
			} else {
				// 对于notepad这种，直接使用可执行文件名
				fullPath = appInfo[1]
			}

			if err != nil {
				log.Printf("WARN: Could not find '%s' in registry, will try to start directly. Error: %v", appName, err)
				// 即使没找到，也尝试用 start 命令，万一它碰巧在PATH里呢
				fullPath = appInfo[1]
			}
			cmd = exec.Command(fullPath)
			log.Println(cmd)
		} else {
			log.Printf("WARN: Unsupported app on Windows: %s", appName)
			return nil
		}

	/*case "darwin": // macOS
		appMap := map[string]string{
			"spotify":       "Spotify",
			"netease_music": "NeteaseMusic",
			"apple_music":   "Music",
			"calculator":    "Calculator", // 增加一个常用应用
		}
		if appFullName, ok := appMap[appName]; ok {
			cmd = exec.Command("open", "-a", appFullName)
		} else {
			log.Printf("WARN: Unsupported app on macOS: %s", appName)
			return nil
		}

	case "linux":
		appMap := map[string]string{
			"spotify":       "spotify",
			"netease_music": "netease-cloud-music",
		}
		if binName, ok := appMap[appName]; ok {
			cmd = exec.Command(binName)
		} else {
			log.Printf("WARN: Unsupported app on Linux: %s", appName)
			return nil
		}*/

	default:
		log.Printf("ERROR: Unsupported operating system: %s", runtime.GOOS)
		return nil
	}

	err := cmd.Start()
	if err != nil {
		log.Printf("ERROR: Failed to launch %s: %v", appName, err)
		return err
	}

	log.Printf("SUCCESS: Launch command for '%s' sent.", appName)
	return nil
}

// ExecuteMediaKey 模拟媒体按键
func ExecuteMediaKey(key string) error {
	log.Printf("INFO: Executing media key: %s", key)
	var keyToTap string
	switch key {
	case "toggle_play_pause":
		keyToTap = "audio_play" // 这个键通常兼具播放和暂停功能
	case "next":
		keyToTap = "audio_next"
	case "prev":
		keyToTap = "audio_prev"
	default:
		log.Printf("WARN: Unknown media key: %s", key)
		return nil // 或者返回一个错误
	}

	// 执行按键模拟
	err := robotgo.KeyTap(keyToTap)
	if err != nil {
		log.Printf("ERROR: Failed to tap key %s: %v", keyToTap, err)
		return err
	}

	log.Printf("SUCCESS: Tapped media key: %s", keyToTap)
	return nil
}
