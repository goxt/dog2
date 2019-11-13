package util

import (
	"github.com/goxt/dog2/env"
)

// 发送邮件
func SendMail(viewData map[string]string, view, title, to, toName string, e bool) bool {

	if env.Config.MsgRouter == "" {
		ConsoleError("server.SendMail通过msgRouter发送邮件，但由于配置不完整导致无法发送！")
	}

	// 请求对象
	var srv = NewService(env.Config.MsgRouter, "api")

	// 请求参数
	srv.AddFormData("method", "Router.Mail.relay")
	srv.AddFormDataForJson("data", map[string]interface{}{
		"app":       env.Config.App.AppName,
		"title":     title,
		"to":        to,
		"toName":    toName,
		"view":      view,
		"viewData":  viewData,
		"exception": e,
	})

	// 开始请求
	err := srv.Post()
	if err != nil {
		return false
	}

	// 解析请求
	var data = struct {
		Code    int
		Msg     string
		Data    interface{}
		Success bool
	}{}
	ok := srv.ToJson(&data)
	if !ok {
		return false
	}

	// 处理业务结果
	if !data.Success {
		srv.LogError("发送邮件，下一级路由节点返回业务失败：" + data.Msg)
	}

	return true
}
