package filesystem

import (
	"os"
	"path/filepath"
)

// IsFile 判断是否是一个文件
func IsFile(fname string) bool {
	_, err := os.Stat(fname)
	if err == nil {
		return true
	}
	return false
}

// CreateFile 创建一个空文件,并且自动创建目录
func CreateFile(fullname string) (bool, error) {
	fpath, fname := filepath.Split(fullname)
	if fpath != "" {
		if err := os.MkdirAll(fpath, os.ModePerm); err != nil {
			return false, err
		}
	}
	f, err := os.Create(fname)
	defer f.Close()
	if err != nil {
		return false, err
	}
	return true, nil
}

// OpenFile 打开指定的文件
func OpenFile(name string, flag int, perm os.FileMode) (*os.File, error) {
	if !IsFile(name) {
		_, err := CreateFile(name)
		if err != nil {
			return nil, err
		}
	}
	return os.OpenFile(name, flag, perm)
}
