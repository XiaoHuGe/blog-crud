package article

import (
	"github.com/Unknwon/com"
	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"net/http"
	"xhblog/models"
	"xhblog/utils/app"
	"xhblog/utils/e"
	"xhblog/utils/setting"
	"xhblog/utils/util"
)

type D map[string]interface{}
// 获取多个文章
func GetArticles(ctx *gin.Context) {
	G := &app.Gin{C:ctx}
	data := make(map[string]interface{})
	maps := make(map[string]interface{})
	valid := validation.Validation{}

	var state int
	if arg := ctx.Query("state"); arg != "" {
		state = com.StrTo(arg).MustInt()
		maps["state"] = state
		valid.Range(state, 0, 1, "state").Message("状态只允许0或1")
	}

	var tagId int = -1
	if arg := ctx.Query("tag_id"); arg != "" {
		tagId = com.StrTo(arg).MustInt()
		maps["tag_id"] = tagId
		valid.Min(tagId, 1, "tag_id").Message("标签必须大于0")
	}

	code := e.INVALID_PARAMS
	if valid.HasErrors() {
		msg := make([]string, len(valid.Errors))
		for i, err := range valid.Errors {
			msg[i] = err.Message
		}
		G.Response(http.StatusBadRequest, code, nil)
		return
	}

	code = e.SUCCESS
	data["lists"] = models.GetArticles(util.GetPage(ctx), setting.PageSize, maps)
	data["total"] = models.GetArticleTotal(maps)
	G.Response(http.StatusOK, code, data)
}

// 获取单个文章
func GetArticle(ctx *gin.Context) {
	G := app.Gin{C:ctx}

	id := com.StrTo(ctx.Param("id")).MustInt()

	valid := validation.Validation{}
	valid.Required(id, "id").Message("id不能为空")
	valid.Min(id,1, "id").Message("id必须大于0")

	code := e.INVALID_PARAMS
	if valid.HasErrors() {
		msg := make([]string, len(valid.Errors))
		for i, err := range valid.Errors {
			msg[i] = err.Error()
		}
		G.Response(http.StatusOK, code, nil)
		return
	}

	if !models.ExistArticleByID(id) {
		code = e.ERROR_NOT_EXIST_ARTICLE
		G.Response(http.StatusOK, code, nil)
		return
	}

	code = e.SUCCESS
	data := models.GetArticle(id)

	G.Response(http.StatusOK, code, data)
}

// 新增文章
func AddArticle(ctx *gin.Context) {
	G := &app.Gin{C:ctx}

	tagId := com.StrTo(ctx.Query("tag_id")).MustInt()
	title := ctx.Query("title")
	desc := ctx.Query("desc")
	content := ctx.Query("content")
	createdBy := ctx.Query("created_by")
	state := com.StrTo(ctx.DefaultQuery("state", "0")).MustInt()

	valid := validation.Validation{}
	valid.Min(tagId,1, "tag_id").Message("标签id必须大于0")
	valid.Required(title, "title").Message("标题不能为空")
	valid.Required(desc, "desc").Message("描述不能为空")
	valid.Required(content, "content").Message("内容不能为空")
	valid.Required(createdBy, "created_by").Message("创建人不能为空")
	valid.Range(state,0, 1, "state").Message("状态只能为0或1")

	code := e.INVALID_PARAMS
	if valid.HasErrors() {
		msg := make([]string, len(valid.Errors))
		for i, err := range valid.Errors {
			msg[i] = err.Message
		}
		G.Response(http.StatusBadRequest, code, msg)
		return
	}

	if !models.ExistTagById(tagId) {
		code = e.ERROR_NOT_EXIST_TAG
		G.Response(http.StatusBadRequest, code, nil)
		return
	}
	code = e.SUCCESS
	data := make(D)
	data["tag_id"] = tagId
	data["title"] = title
	data["desc"] = desc
	data["content"] = content
	data["created_by"] = createdBy
	data["state"] = state
	models.AddArticle(data)
	G.Response(http.StatusOK, code, nil)
}

// 修改文章
func EditArticle(ctx *gin.Context) {
	G := app.Gin{C:ctx}

	id := com.StrTo(ctx.Param("id")).MustInt()
	TagId := com.StrTo(ctx.Query("tag_id")).MustInt()
	title := ctx.Query("title")
	desc := ctx.Query("desc")
	content := ctx.Query("content")
	modifiedBy := ctx.Query("modified_by")

	valid := validation.Validation{}
	var state = -1
	if arg := ctx.Query("state"); arg != "" {
		state = com.StrTo(arg).MustInt()
		valid.Range(state, 0, 1, "state").Message("状态只允许0或1")
	}
	valid.Required(id, "id").Message("id不能为空")
	valid.Required(modifiedBy, "modified_by").Message("修改人不能为空")
	valid.MaxSize(modifiedBy, 20, "modified_by").Message("修改人长度最多为20")

	code := e.INVALID_PARAMS
	if valid.HasErrors() {
		msg := make([]string, len(valid.Errors))
		for i, err := range valid.Errors{
			msg[i] = err.Error()
		}
		G.Response(http.StatusBadRequest, code, msg)
		return
	}

	if !models.ExistArticleByID(id) {
		code = e.ERROR_NOT_EXIST_ARTICLE
		G.Response(http.StatusBadRequest, code, nil)
		return
	}

	code = e.SUCCESS
	maps := make(map[string]interface{})
	if TagId > 0 { maps["tag_id"] = id }
	if title != "" { maps["title"] = title }
	if desc != "" { maps["desc"] = desc }
	if content != "" { maps["content"] = content }
	if modifiedBy != "" { maps["modified_by"] = modifiedBy }

	models.EditArticle(id, maps)

	G.Response(http.StatusOK, code, nil)
}

// 删除文章
func DeleteArticle(ctx *gin.Context) {
	G := &app.Gin{C:ctx}

	id := com.StrTo(ctx.Param("id")).MustInt()

	valid := validation.Validation{}
	valid.Required(id, "id").Message("id不能为空")
	valid.Min(id,1, "id").Message("id必须大于0")

	code := e.SUCCESS
	if valid.HasErrors() {
		code = e.INVALID_PARAMS
		msg := make([]string, len(valid.Errors))
		for i, err := range valid.Errors{
			msg[i] = err.Message
		}
		G.Response(http.StatusBadRequest, code, msg)
		return
	}

	if !models.ExistArticleByID(id) {
		code = e.ERROR_NOT_EXIST_ARTICLE
		G.Response(http.StatusBadRequest, code, nil)
		return
	}
	models.DeleteArticle(id)
	G.Response(http.StatusOK, code, nil)
}