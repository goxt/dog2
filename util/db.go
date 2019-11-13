package util

import (
	"github.com/goxt/dog2/env"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"reflect"
)

func OpenDbConnection(name ...string) *gorm.DB {
	var connName = GetDbConnectionName(name)

	var conf = env.Config.DbConfigs[connName]
	if conf == nil {
		panic("数据库链接`" + connName + "`未定义：config/database.go")
	}

	db, err := gorm.Open(conf.Driver, conf.ConnStr)
	if db == nil {
		panic("数据库链接`" + connName + "`失败")
	}
	if err != nil {
		_ = db.Close()
		panic(err)
	}

	db.LogMode(false)
	return db
}

func CloseDbConnection(db *gorm.DB, isBegin *bool) {
	if db == nil {
		return
	}
	if *isBegin {
		db.Rollback()
	}
	_ = db.Close()
	db = nil
}

func GetDbConnectionName(name []string) string {
	var connName = "default"
	if len(name) != 0 {
		connName = name[0]
	}
	return connName
}

func IsEmpty(p interface{}) bool {

	// 数据库连接
	switch v := p.(type) {
	case gorm.DB:
		if v.RecordNotFound() {
			return true
		}
		if v.Error != nil {
			panic(v.Error)
		}
		return false
	case *gorm.DB:
		if v.RecordNotFound() {
			return true
		}
		if v.Error != nil {
			panic(v.Error)
		}
		return false
	}

	value := reflect.ValueOf(p)
	switch value.Kind() {
	case reflect.String:
		return value.Len() == 0
	case reflect.Bool:
		return !value.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return value.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return value.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return value.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return value.IsNil()
	}
	return reflect.DeepEqual(value.Interface(), reflect.Zero(value.Type()).Interface())

}

func CheckRst(p *gorm.DB) {
	if p.Error != nil {
		panic(p.Error)
	}
}

func DeleteHard(db *gorm.DB, table, where string, value ...interface{}) {
	var sql = "delete from `" + table + "` where " + where
	
	var rst *gorm.DB
	switch len(value) {
	case 0:
		rst = db.Exec(sql)
	case 1:
		rst = db.Exec(sql, value[0])
	case 2:
		rst = db.Exec(sql, value[0], value[1])
	case 3:
		rst = db.Exec(sql, value[0], value[1], value[2])
	case 4:
		rst = db.Exec(sql, value[0], value[1], value[2], value[3])
	case 5:
		rst = db.Exec(sql, value[0], value[1], value[2], value[3], value[4])
	case 6:
		rst = db.Exec(sql, value[0], value[1], value[2], value[3], value[4], value[5])
	case 7:
		rst = db.Exec(sql, value[0], value[1], value[2], value[3], value[4], value[5], value[6])
	case 8:
		rst = db.Exec(sql, value[0], value[1], value[2], value[3], value[4], value[5], value[6], value[7])
	default:
		panic("dog2框架的硬删方法，最多支持8个占位符参数")
	}
	
	CheckRst(rst)
}

func GetTotal(db *gorm.DB) int {
	var total = struct {
		Cnt int
	}{}
	var rst = db.Select("count(1) as cnt").Scan(&total)
	IsEmpty(rst)
	return total.Cnt
}

func SetPage(db *gorm.DB, skip int, take int) *gorm.DB {
	db = db.Limit(take).Offset(skip)
	return db
}
