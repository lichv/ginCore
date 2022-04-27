package ginCore

import (
	"gorm.io/gorm"
	"log"
)

type Container struct {
	Config Configure
	DB     *gorm.DB
	Cache  *RedisService
}

func InitContainer(config Configure) (*Container, error) {
	dbConfig := config.GetDBConfig()
	cacheConfig := config.GetCacheConfig()
	db, e := GetDBByConfig(dbConfig)
	if e != nil {
		log.Panicln("数据库初始化连接失败")
		return nil, e
	}
	var client *RedisService
	if cacheConfig.Type != "" {
		client, e = NewCacheClient(cacheConfig)
		if e != nil {
			log.Panicln("redis 缓存初始化连接失败")
			return nil, e
		}
	} else {
		client = nil
	}
	return &Container{Config: config, DB: db, Cache: client}, nil
}
