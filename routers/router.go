package routers

import (
	"github.com/gin-gonic/gin"
	"xh-blog/controller/v1/tag"
	"xh-blog/utils/setting"
)

func InitRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	gin.SetMode(setting.RunMode)

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	apiv1 := r.Group("/api/v1")
	{
		// 获取标签列表
		apiv1.GET("tags", tag.GetTags)
		// 新建标签
		apiv1.POST("tag", tag.AddTag)
		// 修改标签
		apiv1.PUT("tag/:id", tag.EditTag)
		// 删除标签
		apiv1.DELETE("tag/:id", tag.DeleteTag)
	}

	return r
}