package bootstrap

import (
	"github.com/goxt/dog2/env"
	"github.com/goxt/dog2/util"
)

func loadAppSet() {
	// 禁用请求日志（关闭后并发可突破1万，且不产生错误）
	loadSetCloseRequestLog()

	// 500响应处理
	loadSet500()

	// 404响应处理
	loadSet404()
}

func loadSetCloseRequestLog() {
	// 不启用请求日志即可
	// env.Application.Use(logger.New())
}

func loadSet500() {
	env.Application.OnErrorCode(500, util.SysExceptionHandler)
}

func loadSet404() {
	env.Application.OnErrorCode(404, util.NotFoundHandler)
}
