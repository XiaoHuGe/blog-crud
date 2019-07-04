package jwt

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
	"xhblog/utils/app"
	"xhblog/utils/e"
	"xhblog/utils/jwt"
)

func Jwt() gin.HandlerFunc {
	return func(context *gin.Context) {
		G := app.Gin{C: context}

		var code int
		//var data interface{}

		code = e.SUCCESS
		token := context.Query("token")
		if token == "" {
			code = e.INVALID_PARAMS
		} else {
			claims, err := jwt.ParseToken(token)
			if err != nil {
				code = e.ERROR_AUTH_CHECK_TOKEN_FAIL
			} else if claims.ExpiresAt < time.Now().Unix() {
				code = e.ERROR_AUTH_CHECK_TOKEN_TIMEOUT
			}
		}

		if code != e.SUCCESS {
			G.Response(http.StatusUnauthorized, code, nil)
			context.Abort()
			return
		}

		context.Next()
	}
}
