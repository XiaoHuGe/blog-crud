package tag

import (
	"fmt"
	"github.com/Unknwon/com"
	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"net/http"
	"xh-blog/models"
	"xh-blog/utils/app"
	"xh-blog/utils/e"
	"xh-blog/utils/setting"
	"xh-blog/utils/util"
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

	code := e.SUCCESS

	data["lists"] = models.GetTags(util.GetPage(ctx), setting.PageSize, maps)
	data["total"] = models.GetTagTotal(maps)
	G.Response(http.StatusOK, code, data)
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
	if models.ExistTagByName(name) {
		code = e.ERROR_EXIST_TAG
	} else {
		code = e.SUCCESS
		models.AddTag(name, state, createdBy)
	}
	G.Response(http.StatusOK, code, nil)
	//ctx.JSON(http.StatusOK, gin.H{
	//	"code": code,
	//	"msg": e.GetMsg(code),
	//	"data": make(map[string]string),
	//})
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
	valid.MaxSize(modifiedBy, 20, "modified_by").Message("名称最多为20")
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

	if models.ExistTagById(id) {
		code = e.SUCCESS
		data := make(map[string]interface{})
		data["modified_by"] = modifiedBy
		if name != "" {
			data["name"] = name
		}
		if state != -1 {
			data["state"] = state
		}
		models.EditTag(id, data)
	} else {
		code = e.ERROR_NOT_EXIST_TAG
	}
	G.Response(http.StatusOK, code, nil)
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

	if models.ExistTagById(id) {
		code = e.SUCCESS
		models.DeleteTag(id)
	}
	G.Response(http.StatusOK, code, nil)
}