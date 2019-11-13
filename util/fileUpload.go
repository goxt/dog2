package util

import (
	"encoding/json"
	"github.com/goxt/dog2/env"
	"io"
	"os"
	"time"
)

type UploadFile struct {
	Name string
	Ext  string
	Type string
	Size int64
	Id   string
}

/**
 * 上传文件后，将本地临时目录文件迁移至本地正式目录
 * @param	fileInfo	文件上传回调对象
 * @return				正式目录文件路径（可直接用于保存在数据库中）
 */
func UploadFileMoveToLocal(fileInfo UploadFile) string {

	// 文件路径
	var suffix = time.Now().Format("2006_01_02") + "/" + fileInfo.Id + "." + fileInfo.Ext
	var tmpFile = "./tmpFiles/" + suffix
	var finalFile = "./files/" + env.Config.App.AppName + "/" + suffix

	// 打开文件
	file, err := os.Open(tmpFile)
	defer func() {
		if file != nil {
			_ = file.Close()
		}
	}()
	if err != nil {
		ThrowSys("文件迁移至正式目录出错，文件打开失败："+err.Error(), fileInfo)
	}

	// 迁移文件
	if err := MkDir(StrRL(&finalFile, "/")); err != nil {
		ThrowSys("文件迁移至正式目录出错，文件夹创建失败："+err.Error(), fileInfo)
	}
	out, err := os.OpenFile(finalFile, os.O_WRONLY|os.O_CREATE, 0666)
	defer func() {
		if out != nil {
			_ = out.Close()
		}
	}()
	if err != nil {
		ThrowSys("文件迁移至正式目录失败，文件copy出错"+err.Error(), fileInfo)
	}

	_, _ = io.Copy(out, file)

	return suffix
}

/**
 * 将临时存储的文件迁移到正式目录（支持FTP方式）
 * @return	fileCamp服务器上目录存储的日期，将【日期值+'/'+uuid.ext】这串保存到数据库中即可
 */
func UploadFileMoveWithFileCamp(fileInfo []UploadFile, token string, ip string) string {

	if len(fileInfo) == 0 {
		return ""
	}

	var srv = NewService(env.Config.File.Host, "api")

	// 1. 请求参数
	var files []string
	for _, v := range fileInfo {
		files = append(files, v.Id+"."+v.Ext)
	}
	formData := map[string]interface{}{
		"appName": env.Config.App.AppName,
		"pwd":     env.Config.File.Password,
		"files":   files,
		"token":   token,
		"ip":      ip,
	}
	formBytes, err := json.Marshal(formData)
	if err != nil {
		ThrowSys("请求fileCamp文件迁移前，参数组装失败："+err.Error(), formData)
	}

	srv.SetFormData(&map[string]string{
		"method": "common.file.move",
		"data":   string(formBytes),
	})

	// 2. 开始请求
	err = srv.Post()
	if err != nil {
		ThrowBiz("文件服务器繁忙，请不要关闭页面，稍后再重试")
	}

	// 解析请求
	var data = struct {
		Code    int    `json:"code"`
		Msg     string `json:"msg"`
		Success bool   `json:"success"`
	}{}
	ok := srv.ToJson(&data)
	if !ok {
		ThrowBiz("文件服务器繁忙，请不要关闭页面，稍后再重试")
	}

	if !data.Success {
		if data.Code == 1000 {
			srv.LogError("fileCamp迁移文件失败，返回的业务错误信息" + data.Msg)
		}
		ThrowBiz(data.Msg)
	}

	var data2 = struct {
		Data string `json:"data"`
	}{}
	if err := json.Unmarshal([]byte(srv.ResponseData), &data2); err != nil {
		srv.LogError("业务成功，但是，fileCamp服务端返回的data不是字符串")
		ThrowBiz("fileCamp文件服务器繁忙，请稍后重试")
	}

	return data2.Data
}

/**
 * 拼接文件名
 * @return	从fileCamp服务器返回的日期+文件uuid和后缀，可将此值直接保存到数据库
 */
func UploadFilePathForCamp(fileInfo UploadFile, date string) string {
	return date + "/" + fileInfo.Id + "." + fileInfo.Ext
}
