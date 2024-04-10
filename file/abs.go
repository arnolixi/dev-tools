package file

import (
	"fmt"
	"os"
	"path/filepath"
)

func GetAbsPath() string {
	ex, err := os.Executable()
	if err != nil {
		fmt.Println("Error:", err)
		return ""
	}
	// 获取当前执行文件所在的目录路径
	return filepath.Dir(ex)
}
