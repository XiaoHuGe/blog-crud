package file

import "xhblog/utils/setting"

func GetExclePath() string {
	return setting.AppSetting.ExportSavePath
}

func GetExcleUrl(name string) string {
	return setting.AppSetting.PrefixUrl + "/" + GetExclePath() + name
}
