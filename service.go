package ginCore

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"log"
)

type BaseService struct {
	Context    *gin.Context
	Containter *Container
}

type ServiceInterface interface {
}

func New(ctx *gin.Context, ct *Container) *BaseService {
	return &BaseService{
		Context:    ctx,
		Containter: ct,
	}
}

func (s *BaseService) GetServiceContext() *gin.Context {
	return s.Context
}

func (s *BaseService) SetServiceContext(ctx *gin.Context) *BaseService {
	s.Context = ctx
	return s
}

func (s *BaseService) GetServiceConfig() *Container {
	return s.Containter
}

func (s *BaseService) SetServiceConfig(containter *Container) *BaseService {
	s.Containter = containter
	return s
}

func (s *BaseService) Cache(key string, handle func() (interface{}, error)) (interface{}, error) {
	var result interface{}
	if s.Containter.Config.GetCacheConfig().Type != "" {
		temp, e := s.Containter.Cache.Get(key)
		if e == nil {
			err := json.Unmarshal([]byte(temp), &result)
			if err == nil {
				return result, nil
			} else {
				return temp, nil
			}
		}
	}
	object, err := handle()
	if err != nil {
		log.Panicln(err.Error())
	}
	if s.Containter.Cache.Config.Type != "" {
		err := s.Containter.Cache.Set(key, StrVal(object), 0)
		if err != nil {
			return "", err
		}
	}
	err = json.Unmarshal([]byte(StrVal(object)), &result)
	if err == nil {
		return result, nil
	} else {
		return object, nil
	}
}
