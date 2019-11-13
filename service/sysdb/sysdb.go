package sysdb

import (
	"config/models"
	"github.com/goxt/dog2/util"
	"github.com/jinzhu/gorm"
	"time"
)

type SysDB struct {
	DB      *gorm.DB
	isBegin bool

	TimeStamp  int64
	MyId       uint64
	MyName     string
	MyDeptId   uint64
	MyDeptName string

	runName string
}

/**
 * 建立系统级运行对象，支持数据库操作、系统用户信息、日志等操作，同时开启数据库事务
 * @param	runName		运行名称，用于发生错误时，记录在日志内容中
 */
func NewSuData() *SysDB {
	var this = &SysDB{
		DB:         util.OpenDbConnection(),
		isBegin:    false,
		TimeStamp:  time.Now().Unix(),
		MyId:       models.SystemUserId,
		MyName:     models.SystemUserName,
		MyDeptId:   models.SystemDeptId,
		MyDeptName: models.SystemDeptName,
	}

	return this
}

/**
 * 开启事务
 */
func (this *SysDB) Begin() {
	this.DB = this.DB.Begin()
	this.isBegin = true
}

/**
 * 关闭系统级用户的数据库连接，配合Begin的标准用法：
 * 		su := sysdb.NewSuData()
 * 		defer sysDB.Close()
 */
func (this *SysDB) Close() {
	// 自动回滚并关闭连接，同时过滤BizException异常，其他异常都会记录日志
	util.CloseDbConnection(this.DB, &this.isBegin)
	if e := recover(); e != nil {
		switch v := e.(type) {
		case util.BizException:
		case util.SysException: util.LogException(v.Msg)
		case error: util.LogException(v.Error())
		case string: util.LogException(v)
		default:
			util.LogException("未知错误")
		}
	}
}

/**
 * 提交事务
 */
func (this *SysDB) Commit() {
	this.DB.Commit()
	this.isBegin = false
}
