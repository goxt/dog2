package util

import (
	"github.com/goxt/dog2/env"
	"github.com/goxt/dog2/util/crypt"
	"net/http"
)

func CheckPwdStr(input string, store string) bool {
	var salt = Substr(&store, 3, 6)

	rst, err := crypt.Sha512WithSalt([]byte(input), []byte(salt))
	if err != nil {
		panic(err.Error())
	}

	return rst == store
}

func EncryptPwdStr(input string) string {
	var salt = RandomStr(6)
	rst, err := crypt.Sha512WithSalt([]byte(input), []byte(salt))
	if err != nil {
		panic(err.Error())
	}
	return rst
}

func GetUidByTokenInSession(res http.ResponseWriter, req *http.Request) uint64 {
	token := req.URL.Query().Get("Token")
	if token == "" {
		token = req.Header.Get("Token")
	}
	if token == "" {
		http.Error(res, "客户端没有传入token，无法分析用户信息", 404)
		return 0
	}

	ud := SessionData{}
	err := GetCache(env.Config.App.AppName+":session_"+token, &ud)
	if err != nil {
		LogException("从Redis中读取会话数据失败：" + token)
		http.Error(res, MsgSysError, 500)
		return 0
	}
	if !ud.Flag {
		http.Error(res, "请先登录", 404)
		return 0
	}
	if ud.User == nil {
		LogException("用户的会话数据异常，不存在User属性：" + token)
		http.Error(res, MsgSysError, 500)
		return 0
	}

	var idFloat = ud.User["UserId"].(float64)
	if idFloat == 0 {
		LogException("用户的会话数据异常，没有UserId字段：" + token)
		http.Error(res, MsgSysError, 500)
		return 0
	}

	return uint64(idFloat)
}
