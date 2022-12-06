package fileutil

import "os"

func Mkdir(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		// 先创建文件夹
		err = os.Mkdir(path, 0777)
		if err != nil {
			return err
		}
		// 再修改权限
		err = os.Chmod(path, 0777)
		if err != nil {
			return err
		}
	} else {
		return err
	}
	return nil
}

func RemoveAll(paths ...string) error {
	var err error
	for _, path := range paths {
		if e := os.RemoveAll(path); e != nil {
			err = e
		}
	}
	return err
}

func RemoveFile(path string) error {
	return os.Remove(path)
}

func Exist(filepath string) bool {
	_, err := os.Stat(filepath)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}
