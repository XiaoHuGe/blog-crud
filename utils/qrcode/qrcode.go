package qrcode

import (
	"github.com/EDDYCJY/go-gin-example/pkg/file"
	"github.com/boombuler/barcode"
	"image/jpeg"
	"xhblog/utils/util"
	"github.com/boombuler/barcode/qr"
	"xhblog/utils/setting"
)

type QrCode struct {
	URL    string
	Width  int
	Height int
	Ext    string
	Level  qr.ErrorCorrectionLevel
	Mode   qr.Encoding
}

const EXT_JPG  = ".jpg"

func NewQrCode(url string, width, height int,
	level qr.ErrorCorrectionLevel, mode qr.Encoding) *QrCode {
	return &QrCode{
		URL:url,
		Width:width,
		Height:height,
		Level:level,
		Mode:mode,
		Ext:EXT_JPG,
	}
}

func GetQrCodePath() string {
	return setting.AppSetting.QrCodeSavePath
}

func GetQrCodeFullUrl(name string) string {
	return setting.AppSetting.PrefixUrl + "/" + GetQrCodePath() + name
}

func GetQrCodeFileName(value string) string {
	return util.EncodeMd5(value)
}

func (this *QrCode)GetQrCodeExt() string {
	return this.Ext
}

func (this *QrCode)EnCode(path string) (string, string, error) {
	name := GetQrCodeFileName(this.URL) + this.GetQrCodeExt()
	src := path + name
	if file.CheckNotExist(src) {
		code , err := qr.Encode(this.URL, this.Level, this.Mode)
		if err != nil {
			return "", "", err
		}
		code , err = barcode.Scale(code, this.Width, this.Height)
		if err != nil {
			return "", "", err
		}
		f, err := file.MustOpen(name, path)
		if err != nil {
			return "", "", err
		}
		defer f.Close()

		err = jpeg.Encode(f, code, nil)
		if err != nil {
			return "", "", err
		}
	}
	return  name, path, nil
}