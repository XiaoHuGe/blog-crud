package auth

import (
	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"xhblog/models"
	"xhblog/utils/app"
	"xhblog/utils/e"
	"xhblog/utils/jwt"
)

type auth struct {
	Username string `valid:"Required; MaxSize(20)"`
	Password string `valid:"Required; MaxSize(20)"`
}

func GetAuth(ctx *gin.Context) {
	G := app.Gin{C: ctx}

	username := ctx.Query("username")
	password := ctx.Query("password")

	valid := validation.Validation{}
	a := auth{username, password}
	ok, err := valid.Valid(&a)

	data := make(map[string]interface{})
	code := e.INVALID_PARAMS
	if ok {
		isExist, err := models.CheckAuth(username, password)
		if err != nil {
			log.Println("models.CheckAuth err:", err)
		}
		//isExist := models.CheckAuth(username, password)
		if isExist {
			token, err := jwt.GenerateToken(username, password)
			if err != nil {
				code = e.ERROR_AUTH_CHECK_TOKEN_FAIL
			} else {
				data["token"] = token
				code = e.SUCCESS
			}
		} else {
			code = e.ERROR_AUTH_TOKEN
		}
	}
	if err != nil {
		errs := make([]string, len(valid.Errors))
		for i, err := range valid.Errors {
			errs[i] = err.Message
		}
		G.Response(http.StatusBadRequest, code, errs)
		return
	}
	G.Response(http.StatusOK, code, data)
}
