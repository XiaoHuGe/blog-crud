package app

import (
	"github.com/gin-gonic/gin"
	"sketch/utils/e"
)

type Gin struct {
	C *gin.Context
}

func (this *Gin)Response(httpCode, errCode int, data interface{})  {
	this.C.JSON(httpCode, gin.H{
		"code": httpCode,
		"msg": e.GetMsg(errCode),
		"data": data,
	})
}
