package util

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
)

/**
 * 目录是否存在
 * @param	dirPath		文件路径
 * @return	bool		目录是否存在
 * @return	error		目录存在或不存在时，返回nil；如果目标是非目录或无法读取目标的状态，则返回具体错误信息
 */
func IsExistDir(dirPath string) (bool, error) {
	f, err := os.Stat(dirPath)
	if err == nil {
		if f.IsDir() {
			return true, nil
		}
		return false, errors.New(dirPath + " 不是目录")
	}

	if os.IsNotExist(err) {
		return false, nil
	}

	return false, err
}

/**
 * 创建目录（如果目录存在时，不做任何处理）
 * @param	dirPath		文件路径
 * @return	error		创建文件过程发生的错误
 */
func MkDir(dirPath string) error {
	exist, err := IsExistDir(dirPath)
	if err != nil {
		return err
	}
	if exist {
		return nil
	}

	err = os.MkdirAll(dirPath, os.ModePerm)
	if err != nil {
		return errors.New("目录创建失败：" + err.Error())
	}

	return nil
}

/**
 * 获取当前运行路径
 */
func GetCurrentPath() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		panic(err)
	}
	return strings.Replace(dir, "\\", "/", -1)
}
