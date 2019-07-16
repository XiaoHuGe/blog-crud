package article

import (
	"github.com/Unknwon/com"
	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"net/http"
	"xhblog/models"
	"xhblog/service/article_service"
	"xhblog/utils/app"
	"xhblog/utils/e"
	"xhblog/utils/logging"
	"xhblog/utils/setting"
	"xhblog/utils/util"
)

type D map[string]interface{}

// 获取多个文章
func GetArticles(ctx *gin.Context) {
	G := &app.Gin{C: ctx}
	data := make(map[string]interface{})
	maps := make(map[string]interface{})
	valid := validation.Validation{}

	var state int = -1
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

	articleService := article_service.Article{
		TagID:    tagId,
		PageNum:  util.GetPage(ctx),
		PageSize: setting.AppSetting.PageSize,
		State:    state,
	}
	//count, err := models.GetArticleTotal(maps)
	count, err := articleService.GetCount()
	if err != nil {
		logging.Error(err)
		G.Response(http.StatusInternalServerError, e.ERROR_COUNT_ARTICLE_FAIL, nil)
		return
	}

	articleService.Count = count
	article, err := articleService.GetAll()
	//article, err := models.GetArticles(util.GetPage(ctx), setting.AppSetting.PageSize, maps)
	if err != nil {
		logging.Error(err)
		G.Response(http.StatusInternalServerError, e.ERROR_GET_ARTICLES_FAIL, nil)
		return
	}

	code = e.SUCCESS
	data["lists"] = article
	data["total"] = count
	G.Response(http.StatusOK, code, data)
}

// 获取单个文章
func GetArticle(ctx *gin.Context) {
	G := app.Gin{C: ctx}

	id := com.StrTo(ctx.Param("id")).MustInt()
	valid := validation.Validation{}
	valid.Required(id, "id").Message("id不能为空")
	valid.Min(id, 1, "id").Message("id必须大于0")

	code := e.INVALID_PARAMS
	if valid.HasErrors() {
		msg := make([]string, len(valid.Errors))
		for i, err := range valid.Errors {
			msg[i] = err.Error()
		}
		G.Response(http.StatusOK, code, nil)
		return
	}

	//isExist, err := models.ExistArticleByID(id)
	articleService := article_service.Article{ID: id}
	isExist, err := articleService.ExistByID()
	if !isExist {
		G.Response(http.StatusOK, e.ERROR_NOT_EXIST_ARTICLE, nil)
		return
	}
	if err != nil {
		G.Response(http.StatusInternalServerError, e.ERROR_CHECK_EXIST_ARTICLE_FAIL, nil)
		return
	}

	//data, err := models.GetArticle(id)
	data, err := articleService.Get()
	if err != nil {
		G.Response(http.StatusInternalServerError, e.ERROR_GET_ARTICLES_FAIL, nil)
		return
	}
	G.Response(http.StatusOK, e.SUCCESS, data)
}

// 新增文章
func AddArticle(ctx *gin.Context) {
	G := &app.Gin{C: ctx}

	//tagId := com.StrTo(ctx.Query("tag_id")).MustInt()
	//title := ctx.Query("title")
	//desc := ctx.Query("desc")
	//content := ctx.Query("content")
	//createdBy := ctx.Query("created_by")
	//state := com.StrTo(ctx.DefaultQuery("state", "0")).MustInt()

	addArticleService := article_service.AddArticleService{}
	ctx.ShouldBind(&addArticleService)

	valid := validation.Validation{}
	valid.Min(addArticleService.TagID, 1, "tag_id").Message("标签id必须大于0")
	valid.Required(addArticleService.Title, "title").Message("标题不能为空")
	valid.Required(addArticleService.Desc, "desc").Message("描述不能为空")
	valid.Required(addArticleService.Content, "content").Message("内容不能为空")
	valid.Required(addArticleService.CreatedBy, "created_by").Message("创建人不能为空")
	valid.Range(addArticleService.State, 0, 1, "state").Message("状态只能为0或1")

	code := e.INVALID_PARAMS
	if valid.HasErrors() {
		msg := make([]string, len(valid.Errors))
		for i, err := range valid.Errors {
			msg[i] = err.Message
		}
		G.Response(http.StatusBadRequest, code, msg)
		return
	}

	isExist, err := models.ExistTagById(addArticleService.TagID)
	if err != nil {
		G.Response(http.StatusInternalServerError, e.ERROR_EXIST_TAG_FAIL, nil)
		return
	}
	if !isExist {
		G.Response(http.StatusBadRequest, e.ERROR_NOT_EXIST_TAG, nil)
		return
	}

	articleService := article_service.Article{
		TagID:     addArticleService.TagID,
		Title:     addArticleService.Title,
		Desc:      addArticleService.Desc,
		Content:   addArticleService.Content,
		CreatedBy: addArticleService.CreatedBy,
		State:     addArticleService.State,
	}
	err = articleService.Add()
	//data := make(D)
	//data["tag_id"] = tagId
	//data["title"] = title
	//data["desc"] = desc
	//data["content"] = content
	//data["created_by"] = createdBy
	//data["state"] = state
	//err = models.AddArticle(data)
	if err != nil {
		G.Response(http.StatusOK, e.ERROR_ADD_ARTICLE_FAIL, nil)
	}
	G.Response(http.StatusOK, e.SUCCESS, nil)
}

// 修改文章
func EditArticle(ctx *gin.Context) {
	G := app.Gin{C: ctx}

	id := com.StrTo(ctx.Param("id")).MustInt()
	tagId := com.StrTo(ctx.Query("tag_id")).MustInt()
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
		for i, err := range valid.Errors {
			msg[i] = err.Error()
		}
		G.Response(http.StatusBadRequest, code, msg)
		return
	}

	isExist, err := models.ExistArticleByID(id)
	if !isExist {
		G.Response(http.StatusOK, e.ERROR_NOT_EXIST_ARTICLE, nil)
		return
	}
	if err != nil {
		G.Response(http.StatusInternalServerError, e.ERROR_CHECK_EXIST_ARTICLE_FAIL, nil)
		return
	}

	articleService := article_service.Article{
		ID:         id,
		TagID:      tagId,
		Title:      title,
		Desc:       desc,
		Content:    content,
		ModifiedBy: modifiedBy,
		State:      state,
	}
	err = articleService.Edit()

	if err != nil {
		G.Response(http.StatusOK, e.ERROR_EDIT_ARTICLE_FAIL, nil)
		return
	}
	G.Response(http.StatusOK, e.SUCCESS, nil)
}

// 删除文章
func DeleteArticle(ctx *gin.Context) {
	G := &app.Gin{C: ctx}

	id := com.StrTo(ctx.Param("id")).MustInt()

	valid := validation.Validation{}
	valid.Required(id, "id").Message("id不能为空")
	valid.Min(id, 1, "id").Message("id必须大于0")

	code := e.SUCCESS
	if valid.HasErrors() {
		code = e.INVALID_PARAMS
		msg := make([]string, len(valid.Errors))
		for i, err := range valid.Errors {
			msg[i] = err.Message
		}
		G.Response(http.StatusBadRequest, code, msg)
		return
	}
	articleService := article_service.Article{ID: id}
	isExist, err := articleService.ExistByID()
	//isExist, err := models.ExistArticleByID(id)
	if !isExist {
		G.Response(http.StatusOK, e.ERROR_NOT_EXIST_ARTICLE, nil)
		return
	}
	if err != nil {
		G.Response(http.StatusInternalServerError, e.ERROR_CHECK_EXIST_ARTICLE_FAIL, nil)
		return
	}

	//err = models.DeleteArticle(id)
	err = articleService.Delete()
	if err != nil {
		G.Response(http.StatusOK, e.ERROR_DELETE_ARTICLE_FAIL, nil)
	}
	G.Response(http.StatusOK, code, nil)
}
