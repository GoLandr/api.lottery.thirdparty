package global

import (
	"mycommon/logs"
	"time"

	"api.lottery.thirdparty/config"
	"gopkg.in/redis.v3"
)

const max_redis = 10

var redisChan chan *redis.Client

func init() {
	redisChan = make(chan *redis.Client, max_redis)
	for i := 0; i < max_redis; i++ {
		redisChan <- nil
	}
}

// 获取redis
func getRedis() *redis.Client {
	var client *redis.Client

	select {
	case <-time.After(time.Second * 3):
	case client = <-redisChan:
	}

	// 检测client有效性
	if nil != client {
		_, err := client.Ping().Result()
		if nil != err {
			client.Close()
			client = nil
		}
	}

	// 是否需要新建连接
	if nil != client {
		return client
	}
	client = redis.NewClient(config.GetRedisOption())
	return client
}

func relaseRedis(client *redis.Client) {
	select {
	case <-time.After(time.Second * 3):
		client.Close()
		logs.Debug("delete redis client")
	case redisChan <- client:
	}
}

func RedisGet(key string) (string, error) {
	client := getRedis()
	defer relaseRedis(client)

	val, err := client.Get(key).Result()
	if nil != err {
		//		return "", err
		return "", nil
	}

	return val, nil
}

func RedisDel(key string) error {
	client := getRedis()
	defer relaseRedis(client)

	xx := client.Del(key)

	return xx.Err()
}

// redis查询
func RedisLike(key string) ([]string, error) {
	client := getRedis()
	defer relaseRedis(client)

	val, err := client.Keys(key).Result()
	if nil != err {
		var nullResult = []string{""}
		return nullResult, err
	}

	return val, nil
}

func RedisSet(key string, val string) error {
	client := getRedis()
	defer relaseRedis(client)

	_, err := client.Set(key, val, 0).Result()

	return err
}

func RedisTimeSet(key string, val string, expirations time.Duration) error {
	client := getRedis()
	defer relaseRedis(client)

	_, err := client.Set(key, val, expirations).Result()

	return err
}

func RedisHashGet(key string, hash string) (string, error) {
	client := getRedis()
	defer relaseRedis(client)

	val, err := client.HGet(key, hash).Result()
	if nil != err {
		return "", err
	}

	return val, nil
}

// 获取时间
func redisTime() (int64, error) {
	return time.Now().Unix(), nil
}
