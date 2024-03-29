package tag

import (
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/Unknwon/com"
	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"github.com/tealeg/xlsx"
	"net/http"
	"strconv"
	"time"
	"xhblog/models"
	"xhblog/service/tag_service"
	"xhblog/utils/app"
	"xhblog/utils/e"
	"xhblog/utils/file"
	"xhblog/utils/logging"
	"xhblog/utils/setting"
	"xhblog/utils/util"
)

//获取多个文章标签
func GetTags(ctx *gin.Context) {
	G := &app.Gin{C: ctx}
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
func AddTag(ctx *gin.Context) {
	G := &app.Gin{C: ctx}

	tagService := tag_service.AddTagServer{}
	ctx.ShouldBind(&tagService)

	valid := validation.Validation{}
	valid.Required(tagService.Name, "name").Message("名称不能为空")
	valid.MaxSize(tagService.Name, 100, "name").Message("名称最长为100字符")
	valid.Required(tagService.CreatedBy, "created_by").Message("创建人不能为空")
	valid.MaxSize(tagService.CreatedBy, 20, "created_by").Message("创建人最长为20字符")
	valid.Range(tagService.State, 0, 1, "state").Message("状态值允许0或1")

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
	isExist, err := models.ExistTagByName(tagService.Name)
	if err != nil {
		G.Response(http.StatusInternalServerError, e.ERROR_EXIST_TAG_FAIL, nil)
		return
	}
	if isExist {
		G.Response(http.StatusInternalServerError, e.ERROR_EXIST_TAG, nil)
		return
	}

	err = models.AddTag(tagService.Name, tagService.State, tagService.CreatedBy)
	if err != nil {
		G.Response(http.StatusInternalServerError, e.ERROR_ADD_TAG_FAIL, nil)
		return
	}
	G.Response(http.StatusOK, e.SUCCESS, nil)
}

// 修改文章标签
func EditTag(ctx *gin.Context) {
	G := &app.Gin{C: ctx}
	valid := validation.Validation{}

	id := com.StrTo(ctx.Param("id")).MustInt()
	editService := tag_service.EditTagServer{}
	ctx.ShouldBind(&editService)

	var state = -1
	//if arg := ctx.Query("state"); arg != "" {
	if arg := com.StrTo(editService.State); arg != "" {
		state = com.StrTo(editService.State).MustInt()
		valid.Range(state, 0, 1, "state").Message("状态只允许0或1")
	}

	//}
	valid.Required(id, "id").Message("id不能为空")
	valid.Required(editService.ModifiedBy, "modified_by").Message("修改人不能为空")
	valid.MaxSize(editService.ModifiedBy, 20, "modified_by").Message("修改人长度最多为20")
	valid.MaxSize(editService.Name, 20, "name").Message("名称最多为20")

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
	data["modified_by"] = editService.ModifiedBy
	if editService.Name != "" {
		data["name"] = editService.Name
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
func DeleteTag(ctx *gin.Context) {
	G := &app.Gin{C: ctx}
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

// 导出标签
func ExportTag(ctx *gin.Context) {
	G := app.Gin{C: ctx}

	tags, err := models.GetTags(util.GetPage(ctx), setting.AppSetting.PageSize, make(map[string]interface{}))

	xlsxFile := xlsx.NewFile()
	sheet, err := xlsxFile.AddSheet("标签")
	if err != nil {
		logging.Error(err)
		return
	}

	titles := []string{"ID", "名称", "创建人", "创建时间", "修改人", "修改时间"}
	row := sheet.AddRow()

	var cell *xlsx.Cell
	for _, title := range titles {
		cell = row.AddCell()
		cell.Value = title
	}

	for _, tag := range tags {
		valus := []string{
			strconv.Itoa(tag.ID),
			tag.Name,
			tag.CreatedBy,
			strconv.Itoa(tag.CreatedOn),
			tag.ModifiedBy,
			strconv.Itoa(tag.ModifiedOn),
		}

		row = sheet.AddRow()
		for _, v := range valus {
			cell = row.AddCell()
			cell.Value = v
		}
	}

	time := strconv.Itoa(int(time.Now().Unix()))
	fileName := "tag-" + time + ".xlsx"
	path := file.GetExclePath() + "/" + fileName
	if err := xlsxFile.Save(path); err != nil {
		G.Response(http.StatusOK, e.ERROR, nil)
	}

	G.Response(http.StatusOK, e.SUCCESS, map[string]string{
		"export_url": file.GetExcleUrl(fileName),
	})
}

func ImportTag(ctx *gin.Context) {
	G := app.Gin{C: ctx}
	file, _, err := ctx.Request.FormFile("file")
	if err != nil {
		logging.Error(err)
		G.Response(http.StatusOK, e.INVALID_PARAMS, map[string]string{
			"state": "failed",
		})
		return
	}
	xlsx, err := excelize.OpenReader(file)
	if err != nil {
		logging.Error(err)
		G.Response(http.StatusOK, e.ERROR, map[string]string{
			"state": "failed",
		})
		return
	}
	rows, err := xlsx.GetRows("标签")
	for irow, row := range rows {
		if irow > 0 {
			var data []string
			for _, cell := range row {
				data = append(data, cell)
			}
			models.AddTag(data[1], 1, data[2])
		}
	}
	G.Response(http.StatusOK, e.SUCCESS, map[string]string{
		"state": "success",
	})
}
