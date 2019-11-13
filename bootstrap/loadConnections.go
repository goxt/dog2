package bootstrap

import (
	"config"
	"github.com/go-redis/redis"
	"github.com/goxt/dog2/env"
	"github.com/goxt/dog2/util"
	"strconv"
)

func loadConnections() {
	loadDbConnections()
	loadRedisConnections()
}

/**
 * 定义DB连接
 */
func loadDbConnections() {

	db := env.Config.Database

	env.Config.DbConfigs = map[string]*env.DbConfig{}
	env.Config.DbConfigs["default"] = &env.DbConfig{
		Driver: "mysql",
		ConnStr: db.Account + ":" + db.Password + "@(" + db.Host + ")/" + db.Name +
			"?charset=utf8mb4&collation=utf8mb4_unicode_ci",
	}

	for name, v := range config.DbConnections {
		var conf = &env.DbConfig{
			Driver: v["driver"],
		}
		if v["driver"] != "mysql" {
			panic("暂不支持DB驱动：" + v["driver"])
		}
		conf.ConnStr = v["username"] + ":" + v["password"] + "@(" + v["host"] + ")/" + v["database"] +
			"?charset=" + v["charset"] +
			"&collation=" + v["collation"]
		env.Config.DbConfigs[name] = conf
	}
}

/**
 * 定义Redis连接
 */
func loadRedisConnections() {

	env.Config.RedisConfigs = map[string]*redis.Client{}

	db := env.Config.Redis

	defRedis := redis.NewClient(&redis.Options{
		Addr:     db.Host,
		Password: db.Password,
		DB:       db.Db,
	})
	env.Config.RedisConfigs["default"] = defRedis
	util.SetDefaultRedisClient(defRedis)

	for name, v := range config.RedisConnections {
		DB, err := strconv.Atoi(v["DB"])
		if err != nil {
			panic("配置项错误，config/database.go -> RedisConnections -> DB : " + v["DB"])
		}
		var client = redis.NewClient(&redis.Options{
			Addr:     v["host"],
			Password: v["password"],
			DB:       DB,
		})

		env.Config.RedisConfigs[name] = client
	}

}
