package util

import (
	"encoding/json"
	"github.com/asaskevich/govalidator"
	"github.com/jinzhu/gorm"
	"github.com/kataras/iris"
	"reflect"
	"strings"
	"time"
)

type ApiInterface interface {
	Handler() int
}

type Api struct {
	Login      bool
	AuthKey    string
	Controller func(base *Base) ApiInterface
}

/**
 * 接口结构体中包含的Base子结构体，提供了各种业务结构需要的数据和方法
 */
type Base struct {
	// 是否已经完成接口结构体的初始化过程
	initFlag bool

	// http对象和常用属性
	Ctx    iris.Context
	Method string
	Path   string
	token  string

	// 请求时间
	Time      time.Time
	TimeStamp int64

	// 分页参数
	Page     int
	PageSize int
	Skip     int

	// 默认数据库连接
	DB *gorm.DB

	// 默认数据库，自动回滚。是否开启了事务
	dbBegin bool

	// 会话数据
	Session  *Session
	MyId     uint64
	MyDeptId uint64

	// Json响应
	ResData *JsonData
	IsView  bool
}

/**
 * 控制器中创建Base对象，然后传入业务接口中
 */
func NewBaseApi(needLogin bool, ctx iris.Context) *Base {
	var this = &Base{}

	this.Ctx = ctx

	// 请求属性
	this.initRequestAttr()

	// 获取token
	this.initToken()

	// 生成Session
	this.Session = NewSession(needLogin, this.token)

	// 处理个人信息
	this.initMyInfo()

	return this
}

/**
 * 业务接口初始化Base对象的方法，并处理请求参数
 * 规范：在业务接口的New方法中调用，且每次请求只能调用一次，调用格式如下：
 *		this.Base = this.Init(ctx, this)
 */
func (this *Base) Init(data interface{}) {
	if this.initFlag {
		ThrowSys("接口的初始化过程仅仅只需做一次即可")
	}

	// 分页参数
	this.initPage()

	// 请求参数
	switch this.Method {
	case "POST":
		this.initParamPost(data)
	case "GET":
		this.initParamGet(data)
	default:
		ThrowBiz("业务接口仅支持GET和POST请求", CodeBizFailure)
	}

	// 参数校验
	_, err := govalidator.ValidateStruct(data)
	if err != nil {
		ThrowBiz("参数不符合接口要求："+err.Error(), CodeParamError)
	}

	this.initFlag = true
}

/**
 * 设置本次响应类型为视图，错误处理结果也相应的自动转换成HTML
 */
func (this *Base) SetViewResponse() {
	this.IsView = true
}

/**
 * 校验当前登录人是否拥有某个权限
 * @param	authKey		被校验的权限KEY
 * @return				是否有无此权限
 */
func (this *Base) HasAuth(authKey string) bool {
	for _, v := range this.Session.GetUserAuth() {
		if authKey == v {
			return true
		}
	}
	return false
}

/**
 * Json响应 - 业务失败
 * @param	msg		业务提示信息
 * @param	code	非必传参数，业务CODE，默认为 CodeBizFailure
 * @return			框架响应类型：JSON
 */
func (this *Base) Error(msg string, code ...int) int {
	this.ResData = JsonError(msg, getDefaultErrorCode(code))
	return 1
}

/**
 * Json响应 - 业务成功
 * @param	msg		异常日志的标题
 * @return			框架响应类型：JSON
 */
func (this *Base) Success(msg string) int {
	this.ResData = JsonSuccess(msg)
	return 1
}

/**
 * Json响应 - 业务成功，并附带数据体
 * @param	data	数据体，任意类型的数据（公有属性才可以转成JSON数据）
 * @param	msg		提示信息，非必传，默认为 "成功"
 * @return			框架响应类型：JSON
 */
func (this *Base) SuccessWithData(data interface{}, msg ...string) int {
	this.ResData = JsonSuccessWithData(data, getDefaultMsg(msg))
	return 1
}

/**
 * Json响应 - 业务成功，并返回列表数据
 * @param	list	列表数据，任意类型的数据（列表中的公有属性才可以转成JSON数据）
 * @param	total	总数
 * @return			框架响应类型：JSON
 */
func (this *Base) SuccessWithList(list interface{}, total int) int {
	this.ResData = JsonSuccessWithList(list, total)
	return 1
}

/**
 * 视图响应 - 业务成功，并返回指定的视图文件，相对于 ./view目录
 * @param	viewName	视图名称
 * @return			框架响应类型：视图
 */
func (this *Base) View(viewName string) int {
	this.IsView = true
	err := this.Ctx.View(viewName + ".html")
	if err != nil {
		panic(err)
	}
	return 2
}

/**
 * 开启默认数据库连接
 */
func (this *Base) OpenConnection() {
	var connName = GetDbConnectionName([]string{})
	CloseDbConnection(this.DB, &this.dbBegin)
	this.DB = OpenDbConnection(connName)
}

/**
 * 关闭默认数据库连接
 */
func (this *Base) CloseConnection() {
	CloseDbConnection(this.DB, &this.dbBegin)
}

/**
 * 默认数据库连接，开启事务
 */
func (this *Base) Begin() {
	this.dbBegin = true
	this.DB = this.DB.Begin()
}

/**
 * 默认数据库连接，提交事务
 */
func (this *Base) Commit() {
	this.dbBegin = false
	this.DB = this.DB.Commit()
	if this.DB.Error != nil {
		panic("事务提交失败，" + this.DB.Error.Error())
	}
}

/**
 * 默认数据库连接，回滚事务
 */
func (this *Base) RollBack() {
	this.dbBegin = false
	this.DB = this.DB.Rollback()
}

/**
 * 硬删数据
 * 当模型中有DeleteAt字段时，调DB.Delete是软删，而非硬删，因此业务需要硬删时，需要调此方法来物理删除
 * 该方法已经校验了SQL执行过程中出现的异常
 *
 * @param	table	表名
 * @param	where	筛选子句，可以使用占位符
 * @param	value	筛选子句，代替占位符的值
 */
func (this *Base) DeleteHard(table, where string, value ...interface{}) {
	DeleteHard(this.DB, table, where, value)
}

/**
 * 查询当前DB对象的列表总数
 * @param	obj		DB对象:一般已经指定了表名、过滤条件；没有指定分页、排序、查询列
 */
func (this *Base) GetTotal(obj *gorm.DB) int {
	return GetTotal(obj)
}

/**
 * 查询当前分页数据
 * @param	obj				DB对象，在total的DB对象基础上，已经处理了查询列、排序，无需指定分页
 * @data	interface{}		任意类型的数组或切片结构，外边需要传入指针
 */
func (this *Base) QueryPageList(obj *gorm.DB, data interface{}) {
	obj = SetPage(obj, this.Skip, this.PageSize)
	rst := obj.Scan(data)
	IsEmpty(rst)
}

/**
 * 获取当前访问用户的IP地址
 */
func (this *Base) GetIp() string {
	return this.Ctx.RemoteAddr()
}

func (this *Base) initRequestAttr() {
	this.Path = this.Ctx.Path()
	this.Method = this.Ctx.Method()
	this.Time = time.Now()
	this.TimeStamp = this.Time.Unix()
}

func (this *Base) initPage() {
	this.Page = this.Ctx.URLParamIntDefault("page", 1)
	this.PageSize = this.Ctx.URLParamIntDefault("pagesize", 10)

	if this.Page <= 0 {
		this.Page = 1
	}
	if this.PageSize <= 0 {
		this.PageSize = 10
	}

	this.Skip = (this.Page - 1) * this.PageSize
}

func (this *Base) initToken() {
	var token = this.Ctx.URLParam("Token")
	if token == "" {
		token = this.Ctx.GetHeader("Token")
	}
	this.token = token
}

func (this *Base) initParamPost(data interface{}) {
	ct := this.Ctx.GetHeader("Content-Type")
	if strings.Index(ct, "multipart/form-data") == 0 {
		return
	}

	err := this.Ctx.ReadJSON(data)
	if err != nil {
		str := err.Error()
		ThrowBiz("BODY参数和请求头不符合此接口要求，"+StrLR(&str, "."), CodeParamError)
	}
}

func (this *Base) initParamGet(data interface{}) {
	var jsonElemArr []string

	// 遍历结构体data的反射对象，拼接成json字符串
	val := reflect.ValueOf(data)
	if val.Kind() == reflect.Interface || val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	if val.Kind() != reflect.Struct {
		ThrowSys("解析参数过程出错，接口初始化方法传入的必须是结构体的指针")
	}

	for i := 0; i < val.NumField(); i++ {

		// 结构体属性名，大写开头的属性才需要解析
		key := val.Type().Field(i).Name
		if key[0] < 65 || key[0] > 90 {
			continue
		}

		// 属性TAG中读取json标记
		jsonName := val.Type().Field(i).Tag.Get("json")
		if jsonName == "-" {
			continue
		}
		if jsonName != "" {
			key = jsonName
		}

		// 获取参数，不存在的参数则跳过
		if !this.Ctx.URLParamExists(key) {
			continue
		}
		py := this.Ctx.URLParam(key)

		// 拼接单个json元素
		elem := `"` + key + `":`
		if py == "" {
			elem += `""`
		} else {
			ty := val.Type().Field(i).Type.String()
			if ty == "string" || ty == "*string" {
				elem += `"` + py + `"`
			} else {
				elem += py
			}
		}

		// push到json元素数组中
		jsonElemArr = append(jsonElemArr, elem)
	}

	jsonStr := strings.Join(jsonElemArr, ",")
	if jsonStr == "" {
		jsonStr = "{}"
	} else {
		jsonStr = "{" + jsonStr + "}"
	}

	// 解析参数至结构体
	err := json.Unmarshal([]byte(jsonStr), data)
	if err != nil {
		ThrowBiz("URL参数不符合此接口要求："+jsonStr, CodeParamError)
	}
}

func (this *Base) initMyInfo() {
	// 我的ID
	uid := this.Session.Get("UserId")
	if uid != nil {
		this.MyId = uint64(uid.(float64))
	}

	// 我的分组ID
	gid := this.Session.Get("DeptId")
	if gid != nil {
		this.MyDeptId = uint64(gid.(float64))
	}
}
