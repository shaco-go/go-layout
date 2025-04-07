package router

import (
	"github.com/gin-gonic/gin"
)

func RegisterHttpRouter() *gin.Engine {
	r := gin.New()
	return r
}
