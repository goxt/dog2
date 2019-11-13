package util

import (
	"encoding/json"
	"github.com/go-redis/redis"
	"time"
)

var redisClient *redis.Client

func SetDefaultRedisClient(client *redis.Client) {
	redisClient = client
}

/**
 * 设置缓存数据
 * @param	key			缓存key
 * @param	value		数据，将自动转成JSON存储
 * @param	lifetime	有效时间（秒），0表示永久缓存
 * @return				存储缓存过程发生的错误
 */
func SetCache(key string, value interface{}, lifetime time.Duration) error {
	jsonStr, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return redisClient.Set(key, jsonStr, lifetime*time.Second).Err()
}

/**
 * 获取缓存数据
 * @param	key		缓存key
 * @param	data	预期的数据体（结构体、map等类型）
 * @return			取值过程发生的错误，或json字符串解析成预期格式出错
 */
func GetCache(key string, data interface{}) error {

	if key == "" {
		return nil
	}

	rstBytes, err := redisClient.Get(key).Bytes()

	if err == redis.Nil {
		return nil
	}
	if err != nil {
		return err
	}

	err = json.Unmarshal(rstBytes, data)
	if err != nil {
		return err
	}
	return nil
}

/**
 * 设置过期时间
 * @param	key			缓存key
 * @param	lifetime	新的有效时间（秒）
 * @return				设置过程发生的错误
 */
func ExpireCache(key string, lifetime time.Duration) error {
	return redisClient.Expire(key, lifetime*time.Second).Err()
}

/**
 * 删除缓存
 * @param	key		缓存key
 */
func DeleteCache(key string) {
	redisClient.Del(key)
}
