package upload

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"xhblog/utils/app"
	"xhblog/utils/e"
	"xhblog/utils/file"
	"xhblog/utils/logging"
	"xhblog/utils/setting"
)

func UploadImage(ctx *gin.Context) {
	G := app.Gin{C: ctx}
	code := e.SUCCESS
	data := make(map[string]interface{})
	f, image, err := ctx.Request.FormFile("image")
	if err != nil {
		logging.Error(err)
		code = e.ERROR
		G.Response(http.StatusOK, code, data)
	}
	fileName := image.Filename
	if !file.CheckImageExt(fileName) {
		logging.Error("file.CheckImageExt err")
		code = e.ERROR_UPLOAD_CHECK_IMAGE_FORMAT
		G.Response(http.StatusOK, code, data)
		return
	}
	if !file.CheckImageSize(f) {
		logging.Error("CheckImageSize err")
		code = e.ERROR_UPLOAD_CHECK_IMAGE_FORMAT
		G.Response(http.StatusOK, code, data)
		return
	}
	fileName = file.GetImageName(image.Filename)
	if err := ctx.SaveUploadedFile(image, setting.AppSetting.ImageSavaPath + fileName); err != nil {
		ctx.String(http.StatusBadRequest, fmt.Sprintf("upload file err: %string", err.Error()))
		return
	}
	data["image_url"] = setting.AppSetting.ImagePrefixUrl + "/" + setting.AppSetting.ImageSavaPath + fileName
	data["state"] = "图片上传成功！"
	G.Response(http.StatusOK, code, data)
}
