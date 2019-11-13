package util

import (
	"errors"
	"os/exec"
	"runtime"
	"strings"
)

// 压缩目录（命令行）
func ZipByCommand(fileName string, files string) error {

	// 校验文件名
	if strings.ToLower(StrRR(&fileName, ".")) != "zip" {
		return errors.New("目前只能压缩成zip格式")
	}

	// 校验操作系统
	if runtime.GOOS != "linux" {
		return errors.New("目前只支持linux系统使用压缩指令")
	}

	// 检查是否安装了zip
	cmd1 := exec.Command("zip", "-v")
	out1, e1 := cmd1.CombinedOutput()

	if e1 != nil || strings.Index(string(out1), "Info-ZIP") == -1 {
		return errors.New("系统没有安装zip指令")
	}

	// 调用脚本开始压缩
	cmd2 := exec.Command("./zipFileCamp.sh", fileName, files)
	out2, e2 := cmd2.CombinedOutput()
	if e2 != nil {
		return errors.New("压缩失败：" + string(out2))
	}

	return nil
}
