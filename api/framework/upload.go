package framework

import (
	"config"
	"github.com/goxt/dog2/env"
	"github.com/goxt/dog2/util"
	"io"
	"mime/multipart"
	"os"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

/**
 * 【接口说明】
 * 框架提供的文件上传接口
 */

type upload struct {
	*util.Base

	file     multipart.File
	info     *multipart.FileHeader
	fileSize int64
	fileName string
	fileExt  string
	fileType string
	fileId   string
}

func Upload(base *util.Base) util.ApiInterface {
	this := &upload{Base: base}
	this.Init(this)
	return this
}

func (this *upload) Handler() int {

	this.getFile()
	defer this.closeFile()

	// 获取文件信息
	this.getAndCheckFileSize()
	this.getAndCheckFileName()
	this.getAndCheckFileType()

	// 保存至临时目录
	this.saveFile()

	// 成功
	var data = map[string]interface{}{
		"name": this.fileName,
		"ext":  this.fileExt,
		"type": this.fileType,
		"size": this.fileSize,
		"id":   this.fileId,
	}
	return this.SuccessWithData(data, "上传成功，已保存至临时存储空间，请调用相关业务接口对临时文件进行处理")
}

func (this *upload) getFile() {

	if env.Config.File.Type == env.FileDriverFileCamp {
		util.ThrowSys("服务文件驱动为fileCamp，不应该存储在本地磁盘")
	}

	if this.Ctx.Method() != "POST" {
		util.ThrowBiz("此接口仅支持POST", util.CodeParamError)
	}

	var err error
	this.file, this.info, err = this.Ctx.FormFile("file")
	if err != nil {
		this.closeFile()
		if this.Ctx.FormValue("file") == "" {
			util.ThrowBiz("请先选择文件", util.CodeParamError)
		}
		util.LogError("文件上传失败", err.Error())
		util.ThrowBiz("服务器繁忙，文件上传失败，请稍后重试！", util.CodeSysBusy)
	}
}

func (this *upload) closeFile() {
	if this.file != nil {
		_ = this.file.Close()
	}
}

func (this *upload) getAndCheckFileSize() {
	// 文件大小
	this.fileSize = this.info.Size

	// 文件转移到临时存储空间
	this.Ctx.SetMaxRequestBodySize((config.FileUploadMax + 1) * 1024 * 1024)
	if this.fileSize > config.FileUploadMax*1024*1024 {
		util.ThrowBiz("上传的文件太大，目前最大支持：" + strconv.Itoa(config.FileUploadMax) + "MB")
	}
}

func (this *upload) getAndCheckFileName() {
	// 文件名
	this.fileName = this.info.Filename
	if utf8.RuneCountInString(this.fileName) > env.FileNameMax {
		util.ThrowBiz("文件名超过" + strconv.Itoa(env.FileNameMax) + "个字，请修改后再重新上传！")
	}
}

func (this *upload) getAndCheckFileType() {
	this.fileExt = strings.ToLower(util.StrRR(&this.fileName, "."))

	var inArray = false
	for ft, arr := range env.FileUploadExt {
		if util.InArrayString(this.fileExt, arr) {
			inArray = true
			this.fileType = ft
			break
		}
	}
	if !inArray {
		util.ThrowBiz("目前不支持上传：" + this.fileExt + "格式的文件")
	}
}

func (this *upload) saveFile() {

	var fileId = util.Uuid()
	var tmpDir = "./tmpFiles/" + time.Now().Format("2006_01_02")
	var tmpFile = tmpDir + "/" + fileId + "." + this.fileExt

	err := util.MkDir(tmpDir)
	if err != nil {
		panic(err.Error())
	}

	file, err := os.OpenFile(tmpFile, os.O_WRONLY|os.O_CREATE, 0666)
	defer func() {
		if file != nil {
			_ = file.Close()
		}
	}()
	if err != nil {
		util.ThrowSys("上传文件出错", err.Error())
	}

	_, _ = io.Copy(file, this.file)
}
