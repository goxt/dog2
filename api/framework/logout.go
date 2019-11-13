package framework

import (
	"github.com/goxt/dog2/util"
)

/**
 * 【接口说明】
 * 框架提供的通用退出接口
 */

type logout struct {
	*util.Base
}

func Logout(base *util.Base) util.ApiInterface {
	this := &logout{Base: base}
	this.Init(this)
	return this
}

func (this *logout) Handler() int {
	signOut(this.Base)
	return this.Success("退出成功")
}
