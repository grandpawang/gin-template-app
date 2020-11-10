package utils

import (
	"runtime"
	"strings"
)

// GetSubPath get parent directory
func GetSubPath(dirctory string, level int) string {
	for i := 0; i < level; i++ {
		dirctory = Substring(dirctory, 0, strings.LastIndex(dirctory, "/"))
	}
	return dirctory
}

// GetRealFilePath get file path in runtime
func GetRealFilePath(back2root int) string {
	_, filename, _, _ := runtime.Caller(0)
	filename = GetSubPath(filename, back2root)
	return filename
}
