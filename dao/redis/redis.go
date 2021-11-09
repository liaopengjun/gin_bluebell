package redis

import (
	"fmt"
	"gin_bluebell/settings"
	"github.com/go-redis/redis"
)
var rdb *redis.Client
func Init(cfg *settings.RedisConfig)(err error){
	//连接参数配置
	rdb = redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d",cfg.Host, cfg.Port),
		Password:     cfg.Password, // no password set
		DB:           cfg.DB,       // use default DB
		PoolSize:     cfg.PoolSize,
	})
	_,err = rdb.Ping().Result()
	return
}
//释放资源
func Close()  {
	_ = rdb.Close()
}
