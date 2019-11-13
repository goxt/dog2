package util

import (
	"bufio"
	"errors"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

/**
 * 文件是否存在
 * @param	filePath	文件路径
 * @return	bool		文件是否存在
 * @return	error		文件存在或不存在时，返回nil；如果目标是目录或无法读取目标的状态，则返回具体错误信息
 */
func IsExistFile(filePath string) (bool, error) {
	f, err := os.Stat(filePath)
	if err == nil {
		if f.IsDir() {
			return false, errors.New(filePath + " 不是文件")
		}
		return true, nil
	}

	if os.IsNotExist(err) {
		return false, nil
	}

	return false, err
}

/**
 * 创建空文件（并关闭文件，如果文件存在时，不做任何处理）
 * @param	filePath	文件路径
 * @return	error		创建文件过程发生的错误
 */
func MkFile(filePath string) error {
	exist, err := IsExistFile(filePath)
	if err != nil {
		return err
	}
	if exist {
		return nil
	}

	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return errors.New("文件创建失败：" + err.Error())
	}
	_ = file.Close()

	return nil
}

/**
 * 读取文件的内容（一次性读取完）
 */
func ReadFile(path string) ([]byte, bool) {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		LogError("文件读取失败", err)
		return nil, false
	}
	return bytes, true
}

/**
 * 读取文件大小
 * @param	filePath	文件路径
 * @return	int64		文件大小
 * @return	bool		是否读取成功
 */
func GetFileSize(filePath string) (int64, bool) {
	f, err := os.Stat(filePath)
	if err == nil {
		if f.IsDir() {
			return 0, false
		}
		return f.Size(), true
	}
	return 0, false
}

/**
 * 打开文件，没有则创建
 * @param	filePath	文件路径（相对或绝对）
 */
func OpenFile(filePath string) *os.File {
	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		if file != nil {
			_ = file.Close()
		}
		panic(errors.New("文件打开失败：" + err.Error()))
	}
	return file
}

/**
 * 逐行读取文件内容
 * @param	fileName	文件路径
 * @param	handler		回调方法，传入一行文件内容
 */
func ReadLine(fileName string, handler func(string)) {
	f, err := os.Open(fileName)
	defer func() {
		if f != nil {
			_ = f.Close()
		}
	}()
	if err != nil {
		panic(err)
	}
	buf := bufio.NewReader(f)
	for {
		line, e := buf.ReadString('\n')
		if e != nil {
			if e == io.EOF {
				break
			}
			panic(e)
		}
		line = strings.TrimSpace(line)
		handler(line)
	}
}

/**
 * 处理文件名中特殊的字符，将特殊字符转换成下划线
 * @param	fileName	旧文件名
 * @return				转换后的文件名（只是字符串处理，并非正真修改文件名）
 */
func DealFileName(fileName string) string {
	var s = fileName
	s = strings.ReplaceAll(s, `:`, `_`)
	s = strings.ReplaceAll(s, `*`, `_`)
	s = strings.ReplaceAll(s, `?`, `_`)
	s = strings.ReplaceAll(s, `<`, `_`)
	s = strings.ReplaceAll(s, `>`, `_`)
	s = strings.ReplaceAll(s, `|`, `_`)
	s = strings.ReplaceAll(s, `"`, `_`)
	s = strings.ReplaceAll(s, ` `, `_`)
	return s
}
