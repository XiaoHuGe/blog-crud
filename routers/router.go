package routers

import (
	"github.com/gin-gonic/gin"
	"xhblog/controller/v1/article"
	"xhblog/controller/v1/auth"
	"xhblog/controller/v1/tag"
	"xhblog/middleware/xhjwt"
	"xhblog/utils/setting"
)

func InitRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	gin.SetMode(setting.ServerSetting.RunMode)

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.GET("/auth", auth.GetAuth)

	apiv1 := r.Group("/api/v1")
	apiv1.Use(xhjwt.JWT())
	{
		// 获取标签列表
		apiv1.GET("tags", tag.GetTags)
		// 新建标签
		apiv1.POST("tag", tag.AddTag)
		// 修改标签
		apiv1.PUT("tag/:id", tag.EditTag)
		// 删除标签
		apiv1.DELETE("tag/:id", tag.DeleteTag)

		// 获取文章列表
		apiv1.GET("articles", article.GetArticles)
		// 获取指定文章
		apiv1.GET("article/:id", article.GetArticle)
		// 新建标签
		apiv1.POST("article", article.AddArticle)
		// 修改标签
		apiv1.PUT("article/:id", article.EditArticle)
		// 删除标签
		apiv1.DELETE("article/:id", article.DeleteArticle)
	}

	return r
}
