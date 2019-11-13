package framework

import (
	"config/models"
	"github.com/goxt/dog2/util"
)

/**
 * 【接口说明】
 * 框架提供的通用登录接口
 * 通过账号、密码校验，成功后记录登录日志，并保存会话，返回用户信息、角色、权限、token
 */

type login struct {
	*util.Base
	Account string `json:"account" valid:"required,length(1|32)"`
	Pwd     string `json:"pwd" valid:"required,length(1|255)"`

	user models.User
	data map[string]interface{}
}

func Login(base *util.Base) util.ApiInterface {
	this := &login{Base: base}
	this.Init(this)
	return this
}

func (this *login) Handler() int {

	// 连接数据库，并开启事务
	this.OpenConnection()
	this.Begin()

	// 根据提供的账号，查询用户信息
	this.queryUserByAccount()

	// 校验密码和账号状态
	this.check()

	// 记录登录日志
	this.writeLoginLog()

	// 会话数据
	sessUser := this.sessionDataForUser()

	// 角色数据
	sessRoles := this.sessionDataForRoles()

	// 权限数据
	sessAuth := this.sessionDataForAuth(sessRoles)

	// 注册会话
	this.Session.SignIn("UserId", sessUser, sessAuth, sessRoles)

	// 提交数据
	this.Commit()

	// 响应数据
	var data = map[string]interface{}{
		"User":  sessUser,
		"Roles": sessRoles,
		"Auth":  sessAuth,
		"Token": this.Session.GetToken(),
	}
	return this.SuccessWithData(data)
}

func (this *login) queryUserByAccount() {
	accountExist := findUserByAccount(this.DB, &this.user, this.Account)
	if !accountExist {
		util.ThrowBiz("账号不存在")
	}
}

func (this *login) check() {
	// 校验密码
	if !util.CheckPwdStr(this.Pwd, this.user.Pwd) {
		util.ThrowBiz("密码错误")
	}

	// 校验账号状态
	if this.user.AccountStatus != 2 {
		if this.user.AccountStatus == 1 {
			util.ThrowBiz("您的账号处于初始化状态，请先修改密码", util.CodeToResetPwd)
		}
		util.ThrowBiz("您的账号已被冻结或注销，如需激活，请联系管理员")
	}

	// 校验分组
	if this.user.DeptId == 0 {
		util.ThrowBiz("您尚未归属于任何分组/部门，暂时无法登陆，请联系管理员")
	}
}

func (this *login) writeLoginLog() {
	var log = &models.LoginLog{
		UserId:    this.user.UserId,
		DeptId:    this.user.DeptId,
		Ip:        this.GetIp(),
		Device:    this.Ctx.Request().UserAgent(),
		LoginType: 1,
		CreatedAt: this.TimeStamp,
	}
	util.CheckRst(this.DB.Create(log))
}

func (this *login) sessionDataForUser() map[string]interface{} {
	filter := []string{
		"AccountStatus", "Pwd", "DynamicPwd", "DynamicPwdAt", "Base",
		"DeleteFlag", "CreatedId", "CreatedAt",
		"UpdatedId", "UpdatedAt", "DeletedId", "DeletedAt",
	}
	return util.ToMap(this.user, filter)
}

func (this *login) sessionDataForRoles() []string {
	var data []string
	rst := this.DB.Table("role_user").
		Where("user_id = ?", this.user.UserId).
		Pluck("role_key", &data)
	util.IsEmpty(rst)
	if len(data) == 0 {
		data = []string{}
	}
	return data
}

func (this *login) sessionDataForAuth(roles []string) []string {
	var data []string
	rst := this.DB.Table("role_auth").
		Where("role_key in (?)", roles).
		Pluck("auth_key", &data)
	util.IsEmpty(rst)
	if len(data) == 0 {
		data = []string{}
	}
	return data
}
