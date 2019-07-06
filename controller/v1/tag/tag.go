package tag

import (
	"fmt"
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

//获取多个文章标签
func GetTags(ctx *gin.Context)  {
	G := &app.Gin{C:ctx}
	// 标签名称
	name := ctx.Query("name")

	maps := make(map[string]interface{})
	data := make(map[string]interface{})

	if name != "" {
		maps["name"] = name
	}

	var state int = -1
	if arg := ctx.Query("state"); arg != "" {
		state = com.StrTo(arg).MustInt()
		maps["state"] = state
	}

	tags, err := models.GetTags(util.GetPage(ctx), setting.AppSetting.PageSize, maps)
	if err != nil {
		G.Response(http.StatusInternalServerError, e.ERROR_GET_TAGS_FAIL, data)
		return
	}

	count, err := models.GetTagTotal(maps)
	if err != nil {
		G.Response(http.StatusInternalServerError, e.ERROR_COUNT_TAG_FAIL, data)
		return
	}

	data["lists"] = tags
	data["total"] = count
	G.Response(http.StatusOK, e.SUCCESS, data)
}

// 新增文章标签
func AddTag(ctx *gin.Context)  {
	G := &app.Gin{C:ctx}
	// 标签名称
	name := ctx.Query("name")
	state := com.StrTo(ctx.DefaultQuery("state", "0")).MustInt()
	createdBy := ctx.Query("created_by")

	valid := validation.Validation{}
	valid.Required(name, "name").Message("名称不能为空")
	valid.MaxSize(name, 100, "name").Message("名称最长为100字符")
	valid.Required(createdBy, "created_by").Message("创建人不能为空")
	valid.MaxSize(createdBy, 20,"created_by").Message("创建人最长为20字符")
	valid.Range(state, 0, 1, "state").Message("状态值允许0或1")

	code := e.INVALID_PARAMS
	if valid.HasErrors() {
		msg := make([]string, len(valid.Errors))
		for i, err := range valid.Errors {
			msg[i] = err.Message
		}
		fmt.Println("err : ", valid.Errors)
		G.Response(http.StatusBadRequest, code, nil)
		return
	}

	// 判断标签是否存在
	isExist, err := models.ExistTagByName(name)
	if err != nil {
		G.Response(http.StatusInternalServerError, e.ERROR_EXIST_TAG_FAIL, nil)
		return
	}
	if isExist {
		G.Response(http.StatusInternalServerError, e.ERROR_EXIST_TAG, nil)
		return
	}

	err = models.AddTag(name, state, createdBy)
	if err != nil {
		G.Response(http.StatusInternalServerError, e.ERROR_ADD_TAG_FAIL, nil)
		return
	}
	G.Response(http.StatusOK, e.SUCCESS, nil)
}

// 修改文章标签
func EditTag(ctx *gin.Context)  {
	G := &app.Gin{C:ctx}
	valid := validation.Validation{}

	id := com.StrTo(ctx.Param("id")).MustInt()
	name := ctx.Query("name")
	modifiedBy := ctx.Query("modified_by")

	var state = -1
	if arg := ctx.Query("state"); arg != "" {
		state = com.StrTo(arg).MustInt()
		valid.Range(state, 0, 1, "state").Message("状态只允许0或1")
	}
	valid.Required(id, "id").Message("id不能为空")
	valid.Required(modifiedBy, "modified_by").Message("修改人不能为空")
	valid.MaxSize(modifiedBy, 20, "modified_by").Message("修改人长度最多为20")
	valid.MaxSize(name, 20, "name").Message("名称最多为20")

	code := e.INVALID_PARAMS
	if valid.HasErrors() {
		msg := make([]string, len(valid.Errors))
		for i, err := range valid.Errors {
			msg[i] = err.Message
		}
		G.Response(http.StatusBadRequest, code, nil)
		return
	}

	isExist, err := models.ExistTagById(id)
	if err != nil {
		G.Response(http.StatusInternalServerError, e.ERROR_EXIST_TAG_FAIL, nil)
		return
	}
	if !isExist {
		G.Response(http.StatusOK, e.ERROR_NOT_EXIST_TAG, nil)
		return
	}

	data := make(map[string]interface{})
	data["modified_by"] = modifiedBy
	if name != "" {
		data["name"] = name
	}
	if state != -1 {
		data["state"] = state
	}

	err = models.EditTag(id, data)
	if err != nil {
		G.Response(http.StatusInternalServerError, e.ERROR_EDIT_TAG_FAIL, nil)
		return
	}
	G.Response(http.StatusOK, e.SUCCESS, nil)
}

// 删除文章标签
func DeleteTag(ctx *gin.Context)  {
	G := &app.Gin{C:ctx}
	valid := validation.Validation{}

	id := com.StrTo(ctx.Param("id")).MustInt()
	valid.Required(id, "id").Message("id不能为空")
	valid.Min(id, 1, "id").Message("id必须大于0")

	code := e.INVALID_PARAMS
	if valid.HasErrors() {
		msg := make([]string, len(valid.Errors))
		for i, err := range valid.Errors {
			msg[i] = err.Message
		}
		G.Response(http.StatusBadRequest, code, nil)
		return
	}

	isExist, err := models.ExistTagById(id)
	if err != nil {
		G.Response(http.StatusInternalServerError, e.ERROR_EXIST_TAG_FAIL, nil)
		return
	}
	if !isExist {
		G.Response(http.StatusOK, e.ERROR_NOT_EXIST_TAG, nil)
		return
	}

	err = models.DeleteTag(id)
	if err != nil {
		G.Response(http.StatusInternalServerError, e.ERROR_DELETE_TAG_FAIL, nil)
		return
	}
	G.Response(http.StatusOK, e.SUCCESS, nil)
}