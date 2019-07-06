package file

import (
	"log"
	"mime/multipart"
	"path"
	"strings"
	"xhblog/utils/logging"
	"xhblog/utils/setting"
	"xhblog/utils/util"
)

func GetImageName(name string) string {
	ext := path.Ext(name)
	fileName := strings.TrimSuffix(name, ext)
	fileName = util.EncodeMd5(fileName)

	return fileName + ext
}

func CheckImageExt(fileName string) bool {

	ext := GetExt(fileName)
	for _, allowExt := range setting.AppSetting.ImageAllowExts {
		if strings.ToUpper(allowExt) == strings.ToUpper(ext) {
			return true
		}
	}

	return false
}

func CheckImageSize(f multipart.File) bool {
	size, err := GetSize(f)
	if err != nil {
		log.Println(err)
		logging.Warn(err)
		return false
	}

	return size <= setting.AppSetting.ImageMaxSize
}
