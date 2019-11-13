package util

import (
	"github.com/kataras/iris"
	"strings"
)

const (
	CodeSysError       = 1000
	CodeSysBusy        = 1001
	CodeParamError     = 2000
	CodeBizFailure     = 2001
	CodeIllegalRequest = 2002
	CodeNoLogin        = 2003
	CodeNoAuth         = 2004
	CodeToHome         = 3000
	CodeToResetPwd     = 3001

	MsgSysError   = "抱歉，此功能故障，请联系技术人员"
	MsgSysBusy    = "抱歉，服务器繁忙，请稍后重试"
	MsgParamError = "参数不符合接口要求，请检查字段类型、嵌套结构的格式等"
	MsgNoLogin    = "请先登录"
	MsgNoAuth     = "抱歉，您的角色没有此功能的操作权限，请联系技术人员"
)

var (
	Exists500File = false
)

type SysException struct {
	Msg  string
	Code int
}

type BizException struct {
	Msg  string
	Code int
}

/**
 * 抛出业务失败的异常
 * @param	msg		错误提示信息
 * @param	code	非必传，错误码
 */
func ThrowBiz(msg string, code ...int) {
	var c = CodeBizFailure
	if len(code) >= 1 {
		c = code[0]
	}
	panic(BizException{
		Msg:  msg,
		Code: c,
	})
}

/**
 * 抛出系统异常提示，并记录日志
 * @param	msg		日志的标题
 * @param	data	非必传，日志内容
 */
func ThrowSys(msg string, data ...interface{}) {
	LogError(msg, data)
	panic(SysException{
		Msg:  msg,
		Code: CodeSysError,
	})
}

func SysExceptionHandler(ctx iris.Context) {
	JsonResponse(JsonError(MsgSysError, CodeSysError), ctx)
}

func BizExceptionHandler(biz BizException, ctx iris.Context) {
	JsonResponse(JsonError(biz.Msg, biz.Code), ctx)
}

func SysExceptionViewHandler(ctx iris.Context) {
	if Exists500File {
		ctx.ViewData("ErrorMsg", MsgSysError)
		ctx.View("500.html")
	} else {
		ctx.HTML(strings.ReplaceAll(ViewErrorMsg, "{{.ErrorMsg}}", MsgSysError))
	}
}

func BizExceptionViewHandler(biz BizException, ctx iris.Context) {
	if Exists500File {
		ctx.ViewData("ErrorMsg", biz.Msg)
		ctx.View("500.html")
	} else {
		ctx.HTML(strings.ReplaceAll(ViewErrorMsg, "{{.ErrorMsg}}", biz.Msg))
	}
}

func NotFoundHandler(ctx iris.Context) {
	JsonResponse(JsonError("接口不存在", CodeBizFailure), ctx)
}
