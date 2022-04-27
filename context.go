package ginCore

import "github.com/gin-gonic/gin"

type Context interface {
	SetToken(token string)
	GetToken() string
	SetQuery(query map[string]interface{})
	GetQuery() map[string]interface{}
	SetContext(ctx *gin.Context)
}
