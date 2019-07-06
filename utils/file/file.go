package file

import (
	"io/ioutil"
	"mime/multipart"
	"os"
	"path"
	"xhblog/utils/logging"
)

func GetSize(file multipart.File) (int, error)  {
	content, err := ioutil.ReadAll(file)

	return len(content), err
}

func GetExt(fileName string) string {
	return path.Ext(fileName)
}

func CheckExist(src string) bool {
	_, err := os.Stat(src)
	return os.IsNotExist(err)
}

func CheckPermission(src string) bool {
	_, err := os.Stat(src)
	return os.IsPermission(err)
}

func IsNotExistMkDir(src string) error {
	if exist := CheckExist(src); exist == false {
		if err := MkDir(src); err != nil {
			return err
		}
	}

	return nil
}

func MkDir(src string) error {
	err := os.MkdirAll(src, os.ModePerm)
	if err != nil {
		logging.Error("os.MkdirAll:", err)
		return err
	}
	return nil
}

func Open(name string, flag int, perm os.FileMode) (*os.File, error) {
	f, err := os.OpenFile(name, flag, perm)
	if err != nil {
		return nil, err
	}

	return f, nil
}
