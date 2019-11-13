package util

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"github.com/kataras/iris/core/errors"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
)

/**
 * 发送Http/Https的Server组件
 *
 * host			主机（协议://IP:端口号）
 * path			请求路由
 * formData		请求参数
 * bodyData		请求body
 * ResponseData	响应数据
 * StatusCode	响应状态码
 * validTls		是否启用tls证书
 * tls			tls证书
 * cookie		会话数据
 */
type Service struct {
	host         string
	path         string
	formData     map[string]string
	bodyData     string
	ResponseData string
	StatusCode   int

	validTls bool
	tls      *http.Transport

	cookie *cookiejar.Jar
}

/**
 * 创建Server组件对象
 * @param	host	主机（协议://IP:端口号）
 * @param	path	路由
 * @return			Server组件对象
 */
func NewService(host string, path string) *Service {
	var service = &Service{
		host:     host,
		path:     path,
		formData: map[string]string{},
		validTls: false,
		tls:      nil,
		cookie:   nil,
	}
	return service
}

/**
 * 设置服务地址
 * @param	host	主机（协议://IP:端口号）
 */
func (this *Service) SetHost(host string) {
	this.host = host
}

/**
 * 设置请求路由
 * @param	path	请求路由
 */
func (this *Service) SetPath(path string) {
	this.path = path
}

/**
 * 设置请求参数(formData)
 * @param	formData	请求参数
 */
func (this *Service) SetFormData(formData *map[string]string) {
	this.ClearParam()
	for k, v := range *formData {
		this.formData[k] = v
	}
}

/**
 * 添加请求参数(formData)
 * @param	key		键
 * @param	value	值
 */
func (this *Service) AddFormData(key string, value string) {
	this.formData[key] = value
}

/**
 * 将指定参数转成JSON，并作为一个参数
 */
func (this *Service) AddFormDataForJson(key string, data interface{}) {
	b, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}
	this.formData[key] = string(b)
}

/**
 * 设置请求参数(bodyData)
 * @param	bodyData	请求参数
 */
func (this *Service) SetBodyData(bodyData string) {
	this.ClearParam()
	this.bodyData = bodyData
}

/**
 * 设置请求参数，自动转成JSON格式(bodyData)
 * @param	bodyData	请求参数
 */
func (this *Service) SetBodyDataForJson(bodyData interface{}) {
	this.ClearParam()
	b, err := json.Marshal(bodyData)
	if err != nil {
		panic(err)
	}
	this.bodyData = string(b)
}

/**
 * 清空请求参数
 */
func (this *Service) ClearParam() {
	this.formData = map[string]string{}
	this.bodyData = ""
}

/**
 * 设置TLS对象，并启用证书请求服务
 * @param	tls		tls对象
 */
func (this *Service) SetTls(tls *http.Transport) {
	this.tls = tls
	this.validTls = true
}

/**
 * 设置会话对象
 * @param	jar		会话对象
 */
func (this *Service) SetCookie(jar *cookiejar.Jar) {
	this.cookie = jar
}

/**
 * 通过组件的错误日志方式，记录请求和响应对象
 */
func (this *Service) LogError(msg string) {
	host := this.host + "/" + this.path
	LogError("Server组件捕捉到错误信息，地址："+host+"，错误信息："+msg, map[string]interface{}{
		"request":  this.formData,
		"body":     this.bodyData,
		"response": this.ResponseData,
	})
}

/**
 * 发起Post请求
 * @param	code200		请求码200才算正确请求，其他状态码都一律以错误处理，默认为false
 * @return				如果请求发生错误，返回error对象
 */
func (this *Service) Post(code200 ...bool) error {

	path, code200IsSuccess, client := this.newClient(code200)

	// 请求参数
	var bodyBuffer = &bytes.Buffer{}
	var bodyWriter = multipart.NewWriter(bodyBuffer)
	for key, value := range this.formData {
		err := bodyWriter.WriteField(key, value)
		if err != nil {
			this.LogError("post请求Http服务，参数body写入失败:" + err.Error())
			return err
		}
	}

	// 参数类型
	var contentType = bodyWriter.FormDataContentType()
	_ = bodyWriter.Close()

	// 开始请求，并关闭连接
	rep, err := client.Post(path, contentType, bodyBuffer)
	defer client.CloseIdleConnections()
	if err != nil {
		this.LogError("post请求Http服务，请求失败：" + err.Error())
		return err
	}

	// 解析响应
	defer func() {
		_ = rep.Body.Close()
	}()
	responseData, err := ioutil.ReadAll(rep.Body)
	if err != nil {
		this.LogError("post请求Http服务，响应数据解析失败：" + err.Error())
		return err
	}

	this.ResponseData = string(responseData)
	this.StatusCode = rep.StatusCode
	if code200IsSuccess && this.StatusCode != 200 {
		this.LogError("post请求Http服务，响应码不是200")
		return errors.New("非200的响应状态码")
	}

	return nil
}

/**
 * 发起Get请求
 * @param	code200		请求码200才算正确请求，其他状态码都一律以错误处理，默认为false
 * @return				如果请求发生错误，返回error对象
 */
func (this *Service) Get(code200 ...bool) error {
	getUrl, code200IsSuccess, client := this.newClient(code200)

	// 请求参数
	var urlObj = url.Values{}
	for key, value := range this.formData {
		urlObj.Add(key, value)
	}
	getUrl += "?" + urlObj.Encode()

	// 开始请求，并关闭连接
	rep, err := client.Get(getUrl)
	defer client.CloseIdleConnections()
	if err != nil {
		this.LogError("get请求Http服务，请求失败：" + err.Error())
		return err
	}

	// 解析响应
	defer func() {
		_ = rep.Body.Close()
	}()
	responseData, err := ioutil.ReadAll(rep.Body)
	if err != nil {
		this.LogError("get请求Http服务，响应数据解析失败：" + err.Error())
		return err
	}

	this.ResponseData = string(responseData)
	this.StatusCode = rep.StatusCode
	if code200IsSuccess && this.StatusCode != 200 {
		this.LogError("get请求Http服务，响应码不是200")
		return errors.New("非200的响应状态码")
	}

	return nil
}

/**
 * 发起BodyJson请求
 * @param	code200		请求码200才算正确请求，其他状态码都一律以错误处理，默认为false
 * @return				如果请求发生错误，返回error对象
 */
func (this *Service) Body(code200 ...bool) error {

	path, code200IsSuccess, client := this.newClient(code200)

	// 开始请求，并关闭连接
	reader := strings.NewReader(this.bodyData)
	rep, err := client.Post(path, "application/json", reader)
	defer client.CloseIdleConnections()
	if err != nil {
		this.LogError("post请求Http服务，请求失败：" + err.Error())
		return err
	}

	// 解析响应
	defer func() {
		_ = rep.Body.Close()
	}()
	responseData, err := ioutil.ReadAll(rep.Body)
	if err != nil {
		this.LogError("post请求Http服务，响应数据解析失败：" + err.Error())
		return err
	}

	this.ResponseData = string(responseData)
	this.StatusCode = rep.StatusCode
	if code200IsSuccess && this.StatusCode != 200 {
		this.LogError("post请求Http服务，响应码不是200")
		return errors.New("非200的响应状态码")
	}

	return nil
}

/**
 * 将响应数据转成json
 */
func (this *Service) ToJson(data interface{}) bool {
	err := json.Unmarshal([]byte(this.ResponseData), data)
	if err != nil {
		this.LogError("响应数据无法解析成预期的格式！")
		return false
	}
	return true
}

func (this *Service) newClient(code200 []bool) (string, bool, *http.Client) {
	// 响应200才算成功，默认：是
	var code200IsSuccess = true
	if len(code200) > 0 {
		code200IsSuccess = code200[0]
	}

	// 请求地址
	var path = this.host
	if this.path != "" {
		path += "/" + this.path
	}

	// 客户端
	var client = &http.Client{}

	// 处理https证书
	if strings.Index(this.host, "https") == 0 {
		if !this.validTls {
			client.Transport = &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			}
		} else {
			client.Transport = this.tls
		}
	}

	// 处理cookie
	if this.cookie != nil {
		client.Jar = this.cookie
	}

	return path, code200IsSuccess, client
}
