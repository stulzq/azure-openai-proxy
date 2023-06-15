package util

import (
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
)

func GetCurrentAbsPath() string {
	// 最终方案-全兼容
	dir := getCurrentAbPathByExecutable()
	tmpDir, _ := filepath.EvalSymlinks(os.TempDir())
	// 如果是临时目录或者是goland运行的目录，那么就用caller的方式获取
	if strings.Contains(dir, tmpDir) || strings.Contains(dir, "GoLand") {
		return getCurrentAbPathByCaller()
	}
	return dir
}

func GetWorkdir() string {
	currentAbsPath := GetCurrentAbsPath()
	workDir := filepath.Join(currentAbsPath, "..")
	return workDir
}

// 获取当前执行文件绝对路径
func getCurrentAbPathByExecutable() string {
	exePath, err := os.Executable()
	if err != nil {
		panic(err)
	}
	res, _ := filepath.EvalSymlinks(filepath.Dir(exePath))
	return res
}

// 获取当前执行文件绝对路径（go run）
func getCurrentAbPathByCaller() string {
	var abPath string
	_, filename, _, ok := runtime.Caller(0)
	if ok {
		abPath = path.Dir(filename)
	}
	return abPath
}
