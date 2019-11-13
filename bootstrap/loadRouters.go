package bootstrap

import (
	"github.com/goxt/dog2/env"
	"github.com/goxt/dog2/util"
	"github.com/kataras/iris"
	"router"
)

func loadRouters() {
	// 自定义路由
	for routeName := range router.Router {
		env.Application.Handle("ALL", routeName, router.Handler)
	}

	// 首页
	env.Application.Handle("ALL", "/", router.Index)

	// 加载视图目录
	env.Application.RegisterView(iris.HTML("./views/"+env.Config.App.AppName, ".html"))

	ok, err := util.IsExistFile("views/" + env.Config.App.AppName + "500.html")
	if err == nil && ok {
		util.Exists500File = true
	}
}
