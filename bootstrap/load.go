package bootstrap

import (
	"github.com/goxt/dog2/env"
	"github.com/kataras/iris"
	"strconv"
)

func Load() {
	loadConfigJson()
	loadConnections()
	loadAppSet()
	loadRouters()
	loadValidator()
}

func Run() {
	var host = ":" + strconv.Itoa(env.Config.App.Port)

	// 监听HTTP
	if env.Config.Tls.Crt == "" {
		_ = env.Application.Run(iris.Addr(host))
	} else {
		_ = env.Application.Run(iris.TLS(host, env.Config.Tls.Crt, env.Config.Tls.Key))
	}
}
