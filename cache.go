package ginCore

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"time"
)

type CacheConfig struct {
	Type     string `json:"type"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Name     int    `json:"name"`
	Password string `json:"password"`
}

type RedisService struct {
	Config  *CacheConfig
	Client  *redis.Client
	Context context.Context
}

func NewCacheClient(c *CacheConfig) (rs *RedisService, err error) {
	ctx := context.Background()
	service := RedisService{Config: c, Context: ctx}
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", c.Host, c.Port),
		Password: c.Password,
		DB:       c.Name,
	})
	pong, err := client.Ping(ctx).Result()
	if err != nil {
		fmt.Println(pong, err)
		return nil, err
	}
	fmt.Printf(pong)
	service.Client = client
	return &service, nil
}
func (s *RedisService) Close() {
	err := s.Client.Close()
	if err != nil {
		return
	}
}
func (s *RedisService) Set(key string, value interface{}, t time.Duration) error {
	str, e := json.Marshal(value)
	e = s.Client.Set(s.Context, key, string(str), t).Err()
	if e != nil {
		return e
	}
	return nil
}

func (s *RedisService) Get(key string) (string, error) {
	val, err := s.Client.Get(s.Context, key).Result()
	if err != nil {
		return "", err
	}
	return val, nil
}

func (s *RedisService) ClearAll() {
	s.Client.FlushAll(s.Context)
}

func (s *RedisService) ClearDB() {
	s.Client.FlushDB(s.Context)
}
