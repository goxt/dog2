package util

import (
	"encoding/json"
	"github.com/goxt/dog2/env"
	"strings"
)

/**
 * 发送短信 - 通过msgRouter服务转发
 * @param	mobiles			手机号数组
 * @param	code			短信模板名称
 * @param	templateData	模板参数
 */
func SendSms(mobiles []string, code string, templateData map[string]interface{}) bool {

	var mobileArr []string

	// 过滤空手机号
	for _, v := range mobiles {
		if v != "" {
			mobileArr = append(mobileArr, v)
		}
	}
	if len(mobileArr) == 0 {
		return true
	}

	// 请求对象
	var srv = NewService(env.Config.MsgRouter, "api")

	// 请求参数
	formData := map[string]interface{}{
		"app":          env.Config.App.AppName,
		"mobiles":      strings.Join(mobileArr, ","),
		"templateData": templateData,
		"templateCode": code,
	}
	formBytes, err := json.Marshal(formData)
	if err != nil {
		ThrowSys("msgRouter发送短信前，组装参数失败："+err.Error(), formData)
	}
	srv.SetFormData(&map[string]string{
		"method": "Router.Sms.relay",
		"data":   string(formBytes),
	})

	// 开始请求
	err = srv.Post()
	if err != nil {
		return false
	}

	// 解析请求
	var data = struct {
		Code    int         `json:"code"`
		Msg     string      `json:"msg"`
		Data    interface{} `json:"data"`
		Success bool        `json:"success"`
	}{}
	ok := srv.ToJson(&data)
	if !ok {
		return false
	}

	// 处理业务结果
	if !data.Success {
		srv.LogError("发送短信，下一级路由节点返回业务失败：" + data.Msg)
	}

	return true
}
