package utils

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

// GetPath 根据环境参数返回相应的路径
// 当 env 为 "dev" 时，返回项目根目录路径（包含 go.mod 文件的目录）
// 当 env 为 "prod" 时，返回程序的实际运行目录
func GetPath(env string) (string, error) {
	switch env {
	case "dev":
		// 获取当前工作目录
		wd, err := os.Getwd()
		if err != nil {
			return "", fmt.Errorf("failed to get working directory: %w", err)
		}

		// 向上查找直到找到 go.mod 文件
		dir := wd
		for {
			// 检查当前目录是否有 go.mod 文件
			goModPath := filepath.Join(dir, "go.mod")
			if _, err := os.Stat(goModPath); err == nil {
				// 找到了 go.mod 文件，返回该目录
				return dir, nil
			}

			// 向上移动一级目录
			parentDir := filepath.Dir(dir)

			// 如果已经到达根目录且没有找到 go.mod，则返回错误
			if parentDir == dir {
				return "", fmt.Errorf("unable to find project root directory containing go.mod")
			}

			dir = parentDir
		}

	case "prod":
		// 获取可执行文件的路径
		executable, err := os.Executable()
		if err != nil {
			return "", fmt.Errorf("failed to get executable path: %w", err)
		}

		// 返回可执行文件所在的目录
		return filepath.Dir(executable), nil

	default:
		return "", fmt.Errorf("invalid environment parameter: %s, expected 'dev' or 'prod'", env)
	}
}

// GetDefaultCacheDirectory 获取系统默认缓存目录
func GetDefaultCacheDirectory() (string, error) {
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("could not get user home directory: %w", err)
	}

	appName := getAppName()
	switch runtime.GOOS {
	case "windows":
		return filepath.Join(userHomeDir, "AppData", "Local", appName), nil
	case "darwin":
		return filepath.Join(userHomeDir, "Library", "Caches", appName), nil
	case "linux":
		return filepath.Join(userHomeDir, ".cache", appName), nil
	}
	return "", fmt.Errorf("could not determine cache directory")
}

// getAppName 获取应用名称的可靠方法
func getAppName() string {
	// 方案1: 从环境变量获取
	if appName := os.Getenv("APP_NAME"); appName != "" {
		return appName
	}

	// 方案2: 从可执行文件名获取
	if executable, err := os.Executable(); err == nil {
		executableName := filepath.Base(executable)
		// 移除文件扩展名（特别是在Windows上）
		ext := filepath.Ext(executableName)
		if ext != "" {
			executableName = executableName[:len(executableName)-len(ext)]
		}
		if executableName != "" {
			return executableName
		}
	}

	// 方案3: 从命令行参数获取
	if len(os.Args) > 0 {
		arg0 := filepath.Base(os.Args[0])
		// 移除文件扩展名
		ext := filepath.Ext(arg0)
		if ext != "" {
			arg0 = arg0[:len(arg0)-len(ext)]
		}
		if arg0 != "" {
			return arg0
		}
	}

	// 方案4: 使用默认名称
	return "xiaohongshu"
}

func ReadEmbeddedFile(fs embed.FS, filePath string) (string, error) {
	data, err := fs.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
