package redisx

import (
	"fmt"
	"github.com/go-redis/redis"
)

type Redisx struct {
	Addr     string `json:"addr"`     // format:localhost:6379
	Password string `json:"password"` // 如果没有设置密码则为空
	DB       int    `json:"db"`       // redis连接数据库，默认为0
	Client   *redis.Client
	init     bool
}

func (redisx *Redisx) MarshalDefaults() {
	if redisx.Addr == "" {
		panic("redis Addr cannot be empty")
	}
}

func (redisx *Redisx) Init() {
	if !redisx.init {
		redisx.New()
		redisx.init = true
	}
}

func (redisx *Redisx) New() {
	if redisx.init {
		return
	}

	redisx.MarshalDefaults()

	client := redis.NewClient(&redis.Options{
		Addr:     redisx.Addr,
		Password: redisx.Password,
		DB:       redisx.DB,
	})

	_, err := client.Ping().Result()
	if err != nil {
		panic(fmt.Sprintf("redis Ping err:%v", err))
	}

	redisx.Client = client
}
