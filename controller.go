package ginCore

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
)

type BaseController struct {
	Containter *Container
}

func (c *BaseController) SetConfig(ct *Container) *BaseController {
	c.Containter = ct
	return c
}

func (c *BaseController) Validate(v interface{}) (bool, error) {
	err := validator.New().Struct(v)
	if err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			fmt.Println(err)
			return false, err
		}
		for _, err := range err.(validator.ValidationErrors) {
			fmt.Println(err.Namespace())
			fmt.Println(err.Field())
			fmt.Println(err.StructNamespace())
			fmt.Println(err.StructField())
			fmt.Println(err.Tag())
			fmt.Println(err.ActualTag())
			fmt.Println(err.Kind())
			fmt.Println(err.Type())
			fmt.Println(err.Value())
			fmt.Println(err.Param())
			fmt.Println()
		}
		return false, err
	}
	return true, nil
}

func (c *BaseController) ResponseSuccess(g *gin.Context) {
	g.JSON(200, gin.H{
		"code": 200,
		"data": "",
		"msg":  "操作成功",
	})
}

func (c *BaseController) ResponseData(g *gin.Context, data interface{}) {
	g.JSON(200, gin.H{
		"code": 2000,
		"data": data,
		"msg":  "操作成功",
	})
}

func (c *BaseController) ResponsePageData(g *gin.Context, data interface{}, total int, last int) {
	g.JSON(200, gin.H{
		"code":  2000,
		"data":  data,
		"total": total,
		"last":  last,
		"msg":   "操作成功",
	})
}

func (c *BaseController) ResponseFailureForParameter(g *gin.Context, err interface{}) {
	g.JSON(403, gin.H{
		"code": 4003,
		"data": "",
		"msg":  err,
	})
}

func (c *BaseController) ResponseFailureForFuncErr(g *gin.Context, err interface{}) {
	g.JSON(500, gin.H{
		"code": 5000,
		"data": "",
		"msg":  err,
	})
}

func (c *BaseController) ResponseFailure(g *gin.Context, httpCode, code int, err error) {
	g.JSON(httpCode, gin.H{
		"code": code,
		"data": "",
		"msg":  err.Error(),
	})
}
