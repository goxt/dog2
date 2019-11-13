package framework

import (
	"config/models"
	"github.com/goxt/dog2/util"
	"github.com/jinzhu/gorm"
	"strconv"
)

func findUserByAccount(db *gorm.DB, user *models.User, account string) bool {
	rst := db.Table(user.TableName()).Where("account = ?", account).Find(user)
	return !util.IsEmpty(rst)
}

func signOut(base *util.Base) {
	// 移除登录设备
	uIdStr := strconv.FormatUint(base.MyId, 10)
	cacheKey := base.Session.GetKey()
	base.Session.RemoveDeviceList(uIdStr, cacheKey)

	// 清空会话
	base.Session.SignOut()
}
