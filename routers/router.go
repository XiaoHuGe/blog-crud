package routers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"xhblog/controller/v1/article"
	"xhblog/controller/v1/auth"
	"xhblog/controller/v1/poster"
	"xhblog/controller/v1/tag"
	"xhblog/controller/v1/upload"
	"xhblog/middleware/xhjwt"
	"xhblog/utils/file"
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
	r.StaticFS("/resource/upload/images", http.Dir(file.GetImagePath()))
	r.StaticFS("/resource/export", http.Dir(file.GetExclePath()))
	r.GET("/auth", auth.GetAuth)

	apiv1 := r.Group("/api/v1")
	apiv1.Use(xhjwt.JWT())
	{
		// 获取标签列表
		apiv1.GET("/tags", tag.GetTags)
		// 新建标签
		apiv1.POST("/tag", tag.AddTag)
		// 修改标签
		apiv1.PUT("/tag/:id", tag.EditTag)
		// 删除标签
		apiv1.DELETE("/tag/:id", tag.DeleteTag)

		// 获取文章列表
		apiv1.GET("/articles", article.GetArticles)
		// 获取指定文章
		apiv1.GET("/article/:id", article.GetArticle)
		// 新建标签
		apiv1.POST("/article", article.AddArticle)
		// 修改标签
		apiv1.PUT("/article/:id", article.EditArticle)
		// 删除标签
		apiv1.DELETE("/article/:id", article.DeleteArticle)

		//上传图片
		apiv1.POST("/upload", upload.UploadImage)

		apiv1.POST("tags/export", tag.ExportTag)
		apiv1.POST("tags/import", tag.ImportTag)

		apiv1.POST("poster/generate", poster.GenerateArticlePoster)
	}

	return r
}
