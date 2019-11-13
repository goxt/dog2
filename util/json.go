package util

import (
	"github.com/kataras/iris"
)

type JsonData struct {
	code int
	msg  string
	data interface{}
}

func getDefaultErrorCode(code []int) int {
	if len(code) > 0 {
		return code[0]
	}
	return CodeBizFailure
}

func getDefaultMsg(msg []string) string {
	if len(msg) > 0 {
		return msg[0]
	}
	return "成功"
}

func JsonCustom(msg string, data interface{}, code int) *JsonData {
	if data == nil {
		data = []int{}
	}
	return &JsonData{
		code, msg, data,
	}
}

func JsonError(msg string, code int) *JsonData {
	return &JsonData{
		code, msg, nil,
	}
}

func JsonSuccess(msg string) *JsonData {
	return &JsonData{
		0, msg, []int{},
	}
}

func JsonSuccessWithData(data interface{}, msg string) *JsonData {
	if data == nil {
		data = []int{}
	}
	return &JsonData{
		0, msg, data,
	}
}

func JsonSuccessWithList(list interface{}, total int) *JsonData {
	if list == nil {
		list = []int{}
	}
	return &JsonData{
		0, "查询成功", iris.Map{
			"total": total,
			"list":  list,
		},
	}
}

func JsonResponse(j *JsonData, ctx iris.Context) {
	_, _ = ctx.JSON(iris.Map{
		"code": j.code,
		"msg":  j.msg,
		"data": j.data,
	})
}
