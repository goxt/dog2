package util

import (
	"encoding/json"
	"github.com/goxt/dog2/env"
	"net/url"
)

/**
 * 获取文件下载路径
 * @param	fileName	原始文件名称（一般都是中文）
 * @param	filePath	文件路径（相对项目，比如：2019_10_02/uuid.txt）
 * @return				下载地址（相对路径，需要前端拼上服务器地址才完整）
 */
func GetFileDownloadPath(fileName, filePath string) string {
	var urlObj = url.Values{}
	urlObj.Add("filename", fileName)
	fileName = urlObj.Encode()
	return env.Config.App.AppName + "/" + filePath + "?" + fileName
}

/**
 * 获取文件预览路径（通过调预览服务器获得）
 * @param	filePath	文件路径（相对项目，比如：2019_10_02/uuid.txt）
 * @return				预览地址（相对路径，需要前端拼上服务器地址才完整）
 */
func GetFilePreviewPath(filePath string) string {

	// 请求对象
	var srv = NewService(env.Config.Preview.Host, "api")

	// 表单数据
	var path = env.Config.App.AppName + "/" + filePath
	var formData = map[string]interface{}{
		"clientName": env.Config.App.AppName,
		"clientPwd":  env.Config.Preview.Password,
		"filePath":   path,
	}
	formBytes, err := json.Marshal(formData)
	if err != nil {
		ThrowSys("在线预览失败，请求服务器之前，组装表单数据，json转换失败！", formData)
	}

	var param = map[string]string{
		"method": "preview.office.file",
		"data":   string(formBytes),
	}
	srv.SetFormData(&param)

	// 开始请求
	if err := srv.Post(); err != nil {
		ThrowBiz("文件预览服务器繁忙，请稍后重试")
	}

	// 解析响应
	var data1 = struct {
		Msg     string `json:"msg"`
		Code    int    `json:"code"`
		Success bool   `json:"success"`
	}{}
	if err := json.Unmarshal([]byte(srv.ResponseData), &data1); err != nil {
		srv.LogError("文件预览服务端返回的数据无法解析成业务JSON")
		ThrowBiz("文件预览服务器暂时无法正常访问，请稍后重试")
	}

	if !data1.Success {
		ThrowBiz(data1.Msg)
	}

	var data2 = struct {
		Data string `json:"data"`
	}{}
	if err := json.Unmarshal([]byte(srv.ResponseData), &data2); err != nil {
		srv.LogError("业务成功，但是，文件预览服务端返回的data不是字符串")
		ThrowBiz("文件预览服务器暂时无法正常访问，请稍后重试")
	}

	return data2.Data
}

/**
 * 从FileCamp服务器上指定的文件夹下打包
 * @param	filePath	文件路径（相对项目，比如：2019_10_02/uuid.txt）
 * @return				预览地址（相对路径，需要前端拼上服务器地址才完整）
 */
func DownloadZipFileFromCamp(zipName string, dir []string, files []map[string]string) string {

	var srv = NewService(env.Config.File.Host, "api")

	// 1. 请求参数
	formData := map[string]interface{}{
		"appName": env.Config.App.AppName,
		"pwd":     env.Config.File.Password,
		"zipName": zipName,
		"dir":     dir,
		"files":   files,
	}
	formBytes, err := json.Marshal(formData)
	if err != nil {
		ThrowSys("fileCamp打包失败，请求服务器之前，组装表单数据，json转换失败！", formData)
	}

	srv.SetFormData(&map[string]string{
		"method": "common.file.createZip",
		"data":   string(formBytes),
	})

	// 2. 开始请求
	err = srv.Post()
	if err != nil {
		ThrowBiz("文件服务器繁忙，不要关闭页面，请稍后再重试")
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

	// 处理业务结果
	if !data.Success {
		srv.LogError("fileCamp打包文件失败，返回的业务错误信息" + data.Msg)
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
