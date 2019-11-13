package util

import (
	"config/models"
	"github.com/jinzhu/gorm"
)

func VerifyCodeCheck(db *gorm.DB, userId uint64, code string, currTime int64) {
	var user models.User
	var rst = db.Table(models.User{}.TableName()).Where("user_id = ?", userId).Scan(&user)
	if IsEmpty(rst) {
		ThrowBiz("您的账户已被删除，请联系管理员")
	}
	if user.DynamicPwd == "" {
		ThrowBiz("请先获取验证码")
	}
	if currTime > int64(user.DynamicPwdAt)+300 {
		ThrowBiz("验证码有效时间为5分钟，已失效，请重新获取")
	}
	if !CheckPwdStr(code, user.DynamicPwd) {
		ThrowBiz("验证码错误，请重新输入")
	}
}

func VerifyCodeDeal(db *gorm.DB, userId string) {
	var rst = db.Table("user").Where("user_id = ?", userId).Update(map[string]interface{}{
		"dynamic_pwd":    nil,
		"dynamic_pwd_at": nil,
	})
	CheckRst(rst)
}
