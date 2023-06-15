package util

import (
	"log"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
)

func GetCurrentAbsPath() (string, bool) {
	// Work in dev and prod runtime
	isDebug := false
	absDir := getCurrentAbsPathByExecutable()
	tmpDir, _ := filepath.EvalSymlinks(os.TempDir())
	// If it is a temporary directory or the directory contains 'GoLand' keyword,
	// then retrieve it using the `caller` method.
	if strings.Contains(absDir, tmpDir) || strings.Contains(absDir, "GoLand") {
		isDebug = true
		return getCurrentAbsPathByCaller(), isDebug
	}
	return absDir, isDebug
}

func GetWorkdir() string {
	workDir, isDebug := GetCurrentAbsPath()
	if isDebug {
		workDir = filepath.Join(workDir, "..")
	}
	return workDir
}

// getCurrentAbsPathByExecutable Retrieve the absolute path to the currently executing file.
func getCurrentAbsPathByExecutable() string {
	exePath, err := os.Executable()
	log.Printf("executable path: %s", exePath)
	if err != nil {
		panic(err)
	}
	res, _ := filepath.EvalSymlinks(exePath)
	return path.Dir(res)
}

// getCurrentAbsPathByCaller Retrieve the absolute path to the currently executing file.（by `go run`）
func getCurrentAbsPathByCaller() string {
	var abPath string
	_, filename, _, ok := runtime.Caller(0)
	if ok {
		abPath = path.Dir(filename)
	}
	return abPath
}
