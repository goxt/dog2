package util

import (
	"fmt"
	"github.com/goxt/dog2/env"
	"strconv"
)

// ====================================== 会话结构定义 ======================================

/**
 * 存储在redis中的某个用户的会话数据结构
 */
type SessionData struct {
	// 是否已经登录
	Flag bool

	// 用户信息
	User map[string]interface{}

	// 拥有的权限key的集合
	Auth []string

	// 所属角色id的集合
	Roles []string
}

/**
 * 某个请求生命周期内，业务逻辑的会话数据结构
 */
type Session struct {
	// 本次请求是否需要登录
	needLogin bool

	// 本次请求是否改变了会话数据，关系到请求结束前更新缓存数据库中的[会话内容和过期时间]or[仅过期时间]
	isChanged bool

	// 用户会否已经登录
	isLogged bool

	// 用户的token
	token string

	// 缓存数据库-key
	key string

	// 会话内容
	data SessionData
}

/**
 * 接收到请求后调此方法，为接口生成本次请求的会话数据
 */
func NewSession(needLogin bool, token string) *Session {

	if !env.Config.App.OpenSession && needLogin {
		ThrowBiz("后端服务没有启动会话功能，但是此接口要求登录，请联系技术人员修改配置！")
	}

	var s = &Session{
		needLogin: needLogin,
		token:     token,
	}
	if s.token == "" {
		s.token = Uuid()
	}
	s.key = env.Config.App.AppName + ":session_" + s.token
	if !needLogin {
		return s
	}

	// 从缓存数据库中获取用户的会话数据
	s.initSessionFromCache()

	// 判断是否已经登录
	if s.data.User == nil {
		s.isLogged = false
	} else {
		s.isLogged = s.data.Flag
	}

	return s
}

func (s *Session) initSessionFromCache() {
	if err := GetCache(s.key, &s.data); err != nil {
		ThrowSys("从Redis中读取会话失败：" + err.Error())
	}
}

// ====================================== 获取会话数据 ======================================

func (s *Session) IsNeedLogin() bool {
	return s.needLogin
}

func (s *Session) IsChanged() bool {
	return s.isChanged
}

func (s *Session) IsLogged() bool {
	return s.isLogged
}

func (s *Session) GetToken() string {
	return s.token
}

func (s *Session) GetKey() string {
	return s.key
}

func (s *Session) GetUserData() map[string]interface{} {
	return s.data.User
}

func (s *Session) GetUserAuth() []string {
	return s.data.Auth
}

func (s *Session) GetUserRoles() []string {
	return s.data.Roles
}

func (s *Session) All() SessionData {
	return s.data
}

/**
 * 获取用户信息中指定key的值
 * @param	key		用户信息的某个字段
 * @return			用户信息的某个字段的值
 */
func (s *Session) Get(key string) interface{} {
	return s.data.User[key]
}

// ====================================== 操作会话数据 ======================================

/**
 * 设置或新增用户信息的某个字段和值
 * @param	key		字段名
 * @param	value	字段值
 */
func (s *Session) Set(key string, value interface{}) {
	s.isChanged = true
	s.data.User[key] = value
}

func (s *Session) ClearData() {
	s.data = SessionData{}
	s.isChanged = true
}

func (s *Session) SignOut() {
	s.isLogged = false
	s.ClearData()
}

/**
 * 注册会话
 * @param	uIdField	用户ID字段名，一般都是"UserId"
 * @param	user		用户信息map结构
 * @param	auth		权限key的集合
 * @param	roles		角色Id的集合
 */
func (s *Session) SignIn(uIdField string, user map[string]interface{}, auth, roles []string) {

	// 获取用户ID
	uId1 := user[uIdField]
	if uId1 == nil {
		ThrowSys("注册会话时，传入的用户ID字段不在userMap中")
	}
	uId := strconv.FormatUint(uId1.(uint64), 10)
	if uId == "" {
		ThrowSys("注册会话时，传入的用户ID值在userMap中为空值")
	}

	// 用户会话
	s.data.Flag = true
	s.data.Roles = roles
	s.data.Auth = auth
	s.data.User = user

	// 变更状态
	s.isLogged = true
	s.isChanged = true

	// 用户关系与登录维护设备关系表
	s.AddDeviceList(uId)
}

func (s *Session) AddDeviceList(uId string) {
	var deviceKey = env.Config.App.AppName + ":userDevice_" + uId
	var currDeviceList []string

	err := GetCache(deviceKey, &currDeviceList)
	if err != nil {
		ThrowSys("添加设备：从缓存中读取用户当前登录的设备列表失败", err)
	}

	var isNew = true
	for _, v := range currDeviceList {
		if v == s.key {
			isNew = false
			break
		}
	}

	if isNew {
		currDeviceList = append(currDeviceList, s.key)
		if err := SetCache(deviceKey, currDeviceList, 0); err != nil {
			ThrowSys("添加设备：将用户当前登录的设备列表失败写入缓存时失败", err)
		}
	}
}

func (s *Session) RemoveDeviceList(uId string, key string) {
	var deviceKey = env.Config.App.AppName + ":userDevice_" + uId
	var currDeviceList []string
	var newDeviceList []string

	err := GetCache(deviceKey, &currDeviceList)
	if err != nil {
		ThrowSys("移除设备：从缓存中读取用户当前登录的设备列表失败", err)
	}

	var isReSave = false
	for _, v := range currDeviceList {
		if v == key {
			fmt.Println(key)
			isReSave = true
			continue
		}
		newDeviceList = append(newDeviceList, v)
	}

	if isReSave {
		if err := SetCache(deviceKey, newDeviceList, 0); err != nil {
			ThrowSys("移除设备：将用户当前登录的设备列表失败写入缓存时失败", err)
		}
	}
}

func (s *Session) UpdateSession() {

	if !env.Config.App.OpenSession || s == nil {
		return
	}

	if !s.isChanged && !s.needLogin {
		return
	}

	// 会话数据没变，更新失效时间即可
	lifeTime := env.Config.App.SessionLifeTime
	if !s.isChanged {
		if err := ExpireCache(s.key, lifeTime); err != nil {
			ThrowSys("更新会话过期时间失败："+err.Error(), s.data)
		}
		return
	}

	// 会话数据发生变更，需要更新会话数据及过期时间
	if err := SetCache(s.key, s.data, lifeTime); err != nil {
		ThrowSys("更新会话数据失败："+err.Error(), s.data)
	}
	return
}
