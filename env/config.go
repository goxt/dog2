package env

import (
	"github.com/go-redis/redis"
	"time"
)

var Config = struct {
	App struct {
		// 项目名
		AppName string `json:"-"`

		// 项目名(中文名)
		AppCnName string `json:"name"`

		// 运行端口
		Port int `json:"port"`

		// webSocket端口
		WsPort int `json:"wsPort"`

		// 版本号
		Version string `json:"version"`

		// 是否启动session
		OpenSession bool `json:"-"`

		// 会话超时时间 (单位:秒)
		SessionLifeTime time.Duration `json:"-"`

		// 上传文件大小限制 (单位:M)
		FileUploadMax int `json:"-"`
	}

	Tls struct {
		// crt文件
		Crt string `json:"crt"`

		// key文件
		Key string `json:"key"`
	}

	Database struct {
		// 默认数据库地址
		Host string `json:"host"`

		// 默认数据库库名
		Name string `json:"name"`

		// 默认数据库账号
		Account string `json:"account"`

		// 默认数据库密码
		Password string `json:"password"`
	}

	Redis struct {
		// redis地址
		Host string `json:"host"`

		// redis密码
		Password string `json:"password"`

		// redis库号
		Db int `json:"db"`
	}

	File struct {
		// 文件存储驱动
		Type string `json:"type"`

		// 文件服务器地址
		Host string `json:"host"`

		// 文件服务器地址（下载）
		Down string `json:"down"`

		// 文件服务器密码
		Password string `json:"password"`

		// 扩展数据
		Ext string `json:"ext"`
	}

	Preview struct {

		// 文件预览服务器接口地址
		Host string `json:"host"`

		// 文件预览服务器密码
		Password string `json:"password"`

		// 扩展数据
		Ext string `json:"ext"`
	}

	Log struct {
		// 日志驱动
		Type string `json:"type"`

		// 异常日志是否推送给运维人员
		ErrorEmailOn bool `json:"errorEmailOn"`

		// 异常日志推送的运维人员邮箱地址
		ErrorMailAddr string `json:"errorEmailAddr"`
	}

	// 消息路由器地址
	MsgRouter string `json:"msgRouter"`

	// 配置文件是否已初始化的标记值
	ConfigInitFlag bool `json:"-"`

	// 数据库连接（多个）
	DbConfigs map[string]*DbConfig `json:"-"`

	// Redis连接（多个）
	RedisConfigs map[string]*redis.Client `json:"-"`
}{}
