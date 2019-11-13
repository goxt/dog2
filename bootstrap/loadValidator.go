package bootstrap

import (
	"github.com/asaskevich/govalidator"
	"github.com/goxt/dog2/util"
	"strings"
	"time"
)

func loadValidator() {
	validatorTimestamp()
	validatorIdArray()
	validatorUploadFile()
	validatorUploadFileCamp()
}

func validatorFileObject(day string, v util.UploadFile, check bool) {
	if v.Name == "" || v.Type == "" || v.Ext == "" || v.Id == "" {
		util.ThrowBiz("上传文件的回调参数，格式错误")
	}

	var ext = strings.ToLower(util.StrRR(&v.Name, "."))
	if ext != v.Ext {
		util.ThrowBiz("上传文件的回调参数，后缀名被篡改")
	}

	if !check {
		return
	}

	var file = "./tmpFiles/" + day + "/" + v.Id + "." + v.Ext
	var fileSize, flag = util.GetFileSize(file)
	if !flag {
		util.ThrowBiz("上传的文件已失效，请重新上传")
	}
	if fileSize != v.Size {
		util.ThrowBiz("上传文件的回调参数，文件大小被篡改")
	}
}

func validatorTimestamp() {
	var f = govalidator.CustomTypeValidator(func(i interface{}, o interface{}) bool {
		switch v := i.(type) {
		case int64:
			if v < 0 || v > 4294967295 {
				util.ThrowBiz("时间戳字段值目前支持0~4294967295")
			}
			return true
		}
		return false
	})
	govalidator.CustomTypeTagMap.Set("timestamp", f)
}

func validatorIdArray() {
	var f = govalidator.CustomTypeValidator(func(i interface{}, o interface{}) bool {
		switch v := i.(type) {
		case []string:
			for _, vv := range v {
				if len([]byte(vv)) != 32 {
					util.ThrowBiz("uuid长度必须是32位")
				}
			}
			return true
		case []uint64:
			for _, vv := range v {
				if vv < 0 || vv > 18446744073709551615 {
					util.ThrowBiz("大整型ID值必须在0~18446744073709551615")
				}
			}
			return true
		case []uint32:
			for _, vv := range v {
				if vv < 0 || vv > 4294967295 {
					util.ThrowBiz("小整型ID值必须在0~4294967295")
				}
			}
			return true
		}
		return false
	})
	govalidator.CustomTypeTagMap.Set("idArray", f)
}

func validatorUploadFile() {
	var f = govalidator.CustomTypeValidator(func(i interface{}, o interface{}) bool {
		var day = time.Now().Format("2006_01_02")
		switch v := i.(type) {
		case util.UploadFile:
			validatorFileObject(day, v, true)
			return true
		case []util.UploadFile:
			for _, vv := range v {
				validatorFileObject(day, vv, true)
			}
			return true
		}
		return false
	})
	govalidator.CustomTypeTagMap.Set("uploadFile", f)
}

func validatorUploadFileCamp() {
	var f = govalidator.CustomTypeValidator(func(i interface{}, o interface{}) bool {
		var day = time.Now().Format("2006_01_02")
		switch v := i.(type) {
		case util.UploadFile:
			validatorFileObject(day, v, false)
			return true
		case []util.UploadFile:
			for _, vv := range v {
				validatorFileObject(day, vv, false)
			}
			return true
		}
		return false
	})
	govalidator.CustomTypeTagMap.Set("uploadFileCamp", f)
}
