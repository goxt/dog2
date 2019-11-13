package framework

import (
	"config/models"
	"github.com/goxt/dog2/util"
)

/**
 * 【接口说明】
 * 框架提供的重置密码接口
 */

type resetPwd struct {
	*util.Base
	Account string `json:"account" valid:"required,length(1|20)"`
	OldPwd  string `json:"oldPwd" valid:"required,length(1|255)"`
	NewPwd  string `json:"newPwd" valid:"required,length(1|255)"`

	user models.User
}

func ResetPwd(base *util.Base) util.ApiInterface {
	this := &resetPwd{Base: base}
	this.Init(this)
	return this
}

func (this resetPwd) Handler() int {
	// 连接数据库，并开启事务
	this.OpenConnection()
	this.Begin()

	// 根据提供的账号，查询用户信息
	this.queryUserByAccount()

	// 校验密码和账号状态
	this.check()

	// 更新密码
	this.resetPwd()

	// 清空会话并退出
	signOut(this.Base)

	// 提交事务
	this.Commit()

	// 响应
	return this.Success("密码修改成功，请重新登录")
}

func (this *resetPwd) queryUserByAccount() {
	accountExist := findUserByAccount(this.DB, &this.user, this.Account)
	if !accountExist {
		util.ThrowBiz("账号不存在")
	}
}

func (this *resetPwd) check() {
	// 校验密码
	if !util.CheckPwdStr(this.OldPwd, this.user.Pwd) {
		util.ThrowBiz("旧密码错误")
	}

	//  校验部门
	if this.user.DeptId == 0 {
		util.ThrowBiz("您尚未归属于任何部门，暂时无法使用该账号，请联系管理员")
	}
}

func (this *resetPwd) resetPwd() {
	rst := this.DB.Table("user").Where("user_id = ?", this.user.UserId).
		Updates(map[string]interface{}{
			"account_status": 2,
			"pwd":            util.EncryptPwdStr(this.NewPwd),
		})
	util.CheckRst(rst)
}
