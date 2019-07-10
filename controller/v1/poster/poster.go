package poster

import (
	"github.com/boombuler/barcode/qr"
	"github.com/gin-gonic/gin"
	"net/http"
	"xhblog/utils/app"
	"xhblog/utils/e"
	"xhblog/utils/qrcode"
)

const QRCODE_URL  = "https://github.com/XiaoHuGe?tab=repositories"
func GenerateArticlePoster(ctx *gin.Context) {
	G := app.Gin{C:ctx}
	qrc := qrcode.NewQrCode(QRCODE_URL, 300, 300, qr.M, qr.Auto)
	path := qrcode.GetQrCodePath()
	_, _, err := qrc.EnCode(path)
	if err != nil {
		G.Response(http.StatusOK, e.ERROR, nil)
		return
	}

	G.Response(http.StatusOK, e.SUCCESS, nil)
}
