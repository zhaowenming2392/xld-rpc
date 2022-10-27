package redis

import (
	"errors"
	"strconv"
	"time"

	"github.com/gomodule/redigo/redis"
)

//Conn 接口有一个通用的方法来执行 Redis 命令：
//Do(commandName string, args ...interface{}) (reply interface{}, err error)

//TODO 可能需要一个助手结构，可以新建以及提供默认函数

//NewRedis New Redis
func NewRedis(host string, port int, username, password string) (redis.Conn, error) {
	var err error
	var conn redis.Conn

	//c, err := redis.Dial("tcp", ":6379")
	if username == "" {
		conn, err = redis.Dial("tcp", host+":"+strconv.Itoa(port))
	} else {
		conn, err = redis.Dial("tcp", host+":"+strconv.Itoa(port), redis.DialUsername(username), redis.DialPassword(password))
	}

	if err != nil {
		// handle error
		return nil, err
	}

	return conn, nil
}

//配置
type RedisConfig struct {
	host     string
	port     int
	username string
	password string
}

//默认配置
var defaultRedisConfig = &RedisConfig{
	host: "127.0.0.1",
	port: 6379,
}

func SetFastRedisConfig(host string, port int, username, password string) {
	defaultRedisConfig.host = host
	defaultRedisConfig.port = port
	defaultRedisConfig.username = username
	defaultRedisConfig.password = password
}

//连接池
var redisPool *redis.Pool

func newPool(config *RedisConfig) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		// Dial or DialContext must be set. When both are set, DialContext takes precedence over Dial.
		Dial: func() (redis.Conn, error) {
			return NewRedis(config.host, config.port, config.username, config.password)
		},
	}
}

var defaultConn *redis.Conn

//快速执行redis语句，采用默认的配置
func FastRedisDo(commandName string, args ...interface{}) (reply interface{}, err error) {
	if defaultRedisConfig == nil {
		return nil, errors.New("请先配置默认的redis配置")
	}
	//连接池初始话
	if redisPool == nil {
		redisPool = newPool(defaultRedisConfig)
	}

	if defaultConn == nil {
		//获取一个连接
		*defaultConn = redisPool.Get()
	}

	defer (*defaultConn).Close()

	return (*defaultConn).Do(commandName, args...)
}
