package util

import (
	"encoding/json"
	"github.com/goxt/dog2/env"
	goLog "log"
	"os"
	"runtime/debug"
	"strings"
	"time"
)

var filePath string
var fileObj *os.File

/**
 * 写日志 - Info类
 * @param	msg		日志标题
 * @param	data	日志内容（复杂数据类型转换成JSON格式输出）
 */
func LogInfo(msg string, data interface{}) {
	go log(env.LogLvInfo, msg, &data)
}

/**
 * 写日志 - Error类
 * @param	msg		日志标题
 * @param	data	日志内容（复杂数据类型转换成JSON格式输出）
 */
func LogError(msg string, data interface{}) {
	go log(env.LogLvError, msg, &data)
}

/**
 * 写日志 - 异常类（Error类型，堆栈信息作为日志内容）
 * @param	msg		日志标题
 */
func LogException(msg string) {
	var debugArr []string
	for i, v := range strings.Split(string(debug.Stack()), "\n") {
		if i >= 2 && i <= 8 {
			continue
		}
		debugArr = append(debugArr, v)
	}
	var v interface{} = strings.Join(debugArr, "\n")
	go log(env.LogLvError, msg, &v)

	if env.Config.Log.ErrorEmailOn {
		go logErrorToAdmin(msg, &v)
	}
}

/**
 * 打印错误日志到控制台，支持颜色
 * @param	content		日志内容
 */
func ConsoleError(content string) {
	goLog.Printf("%c[%d;%d;%dm%s%c[0m\n", 0x1B, 0, 40, 31, content, 0x1B)
}

func log(level string, msg string, data *interface{}) {

	// 1.配置文件未加载，以控制台方式输出日志
	if !env.Config.ConfigInitFlag {
		logConsole(level, msg, data)
		return
	}

	// 2.配置文件已正确加载，以配置的驱动方式输出日志
	logConsole(level, msg, data)
	switch env.Config.Log.Type {
	case env.LogDriverConsole:
	case env.LogDriverFile:
		logFile(level, msg, data)
	case env.LogDriverServer:
		logServer(level, msg, data)
	}
}

func logConsole(level string, msg string, content *interface{}) {
	goLog.Printf("%s\r\n%v\r\n", "["+level+"]"+" "+msg, *content)
}

func logFile(level string, msg string, content *interface{}) {

	defer func() {
		if err := recover(); err != nil {
			goLog.Println("日志写入失败：" + err.(error).Error())
		}
	}()

	day := time.Now().Format("2006_01_02")
	target := "./logs/" + env.Config.App.AppName
	f := strings.TrimRight(target, "/") + "/" + day + ".log"

	if filePath != f {

		filePath = f

		// 关闭文件
		if fileObj != nil {
			_ = fileObj.Close()
		}

		// 建立目录
		if err := MkDir(target); err != nil {
			panic(err.Error())
		}

		// 创建文件
		if err := MkFile(filePath); err != nil {
			panic(err.Error())
		}

		// 追写文件
		ff, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
		if err != nil {
			panic(err)
		}
		fileObj = ff
	}

	prefix := time.Now().Format("[2006-01-02 15:04:05]") + " [" + level + "] "

	// 数据转成JSON字符串
	switch v := (*content).(type) {
	case string:
		buf := []byte(prefix + msg + "\n" + v + "\n\n")
		_, _ = fileObj.Write(buf)
	default:
		j, err := json.Marshal(content)
		if err != nil {
			panic(err.Error())
		}
		buf := []byte(prefix + msg + "\n" + string(j) + "\n\n")
		_, _ = fileObj.Write(buf)
	}
}

func logServer(level string, msg string, content *interface{}) {
	// todo:
	goLog.Println("暂不支持文件服务器方式：", level, msg, content)
}

func logErrorToAdmin(msg string, data *interface{}) {
	var title = "系统发生异常，请及时处理"
	var toName = "项目负责人"
	var content = "<b style='color:red;'>" + msg + "</b>" + "\n" + (*data).(string)
	content = strings.Replace(content, "\n", "<br/>", -1)
	SendMail(map[string]string{
		"content": content,
		"toName":  toName,
	}, "common", title, env.Config.Log.ErrorMailAddr, toName, true)
}
