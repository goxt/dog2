package bootstrap

import (
	"config"
	"encoding/json"
	"github.com/goxt/dog2/env"
	"github.com/goxt/dog2/util"
)

/**
 * 加载项目的配置文件，和项目同名的json文件
 */
func loadConfigJson() {

	// 读取文件
	appName := config.AppName
	fileName := appName + ".json"
	b, ok := util.ReadFile(fileName)
	if !ok {
		panic("由于配置文件读取失败，服务无法正常启动，请确保项目名和json文件是否同名")
	}

	// 解析json
	err := json.Unmarshal(b, &env.Config)
	if err != nil {
		panic("由于配置文件格式错误，无法解析成标准配置数据：" + err.Error())
	}

	// 解析业务配置json
	err = json.Unmarshal(b, &config.Config)
	if err != nil {
		panic("由于配置文件格式错误，无法解析成业务自己的Config标准配置数据：" + err.Error())
	}

	// 处理特殊字段
	env.Config.App.AppName = appName
	env.Config.App.OpenSession = config.OpenSession
	env.Config.App.SessionLifeTime = config.SessionLifeTime
	env.Config.App.FileUploadMax = config.FileUploadMax
	env.Config.ConfigInitFlag = true

	// 校验Tls文件
	if env.Config.Tls.Crt != "" || env.Config.Tls.Key != "" {
		ok, err := util.IsExistFile(env.Config.Tls.Crt)
		if !ok || err != nil {
			panic("检测到启用tls，但无法读取crt文件")
		}

		ok, err = util.IsExistFile(env.Config.Tls.Key)
		if !ok || err != nil {
			panic("检测到启用tls，但无法读取key文件")
		}
	}

	// 校验File.Type
	var arrFile = []string{env.FileDriverLocal, env.FileDriverFileCamp}
	if !util.InArrayString(env.Config.File.Type, arrFile) {
		panic("配置文件错误，文件驱动类型不合法")
	}

	// 校验File.Host和pwd
	if env.Config.File.Type == env.FileDriverFileCamp {
		if env.Config.File.Host == "" || env.Config.File.Password == "" {
			panic("配置文件错误，fileCamp方式的服务器不和密码不能为空")
		}
	}

	// 校验MsgRouter
	if env.Config.MsgRouter == "" {
		env.Config.MsgRouter = "http://127.0.0.1:5223"
	}

	// 校验Log.Type
	var arrLog = []string{env.LogDriverConsole, env.LogDriverFile, env.LogDriverServer}
	if !util.InArrayString(env.Config.Log.Type, arrLog) {
		panic("配置文件错误，日志驱动类型不合法")
	}

	// 校验Log.MailAddr
	if env.Config.Log.ErrorEmailOn && env.Config.Log.ErrorMailAddr == "" {
		panic("配置文件错误，开启错误日志推送时，必须指定运维人员的邮箱地址")
	}
}
