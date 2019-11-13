package util

import (
	"github.com/axgle/mahonia"
	"strings"
)

func StrLR(str *string, subStr string) string {
	pot := strings.Index(*str, subStr)
	if pot == -1 {
		return ""
	}

	subStrPot := pot + len(subStr)
	return (*str)[subStrPot:]
}

func StrLL(str *string, subStr string) string {

	pot := strings.Index(*str, subStr)
	if pot == -1 {
		return ""
	}

	return (*str)[0:pot]
}

func StrRL(str *string, subStr string) string {

	pot := strings.LastIndex(*str, subStr)
	if pot == -1 {
		return ""
	}

	return (*str)[0:pot]
}

func StrRR(str *string, subStr string) string {

	pot := strings.LastIndex(*str, subStr)
	if pot == -1 {
		return ""
	}

	pot += len(subStr)
	return (*str)[pot:]
}

func ConvertToByte(src *string, srcCode string, targetCode string) []byte {
	srcCoder := mahonia.NewDecoder(srcCode)
	srcResult := srcCoder.ConvertString(*src)
	tagCoder := mahonia.NewDecoder(targetCode)
	_, cdata, _ := tagCoder.Translate([]byte(srcResult), true)
	return cdata
}

func Substr(str *string, start int, length ...int) string {
	rs := []rune(*str)
	var strLen = len(rs)

	if start < 0 {
		start = 0
	}

	var end = strLen
	if len(length) != 0 {
		end = start + length[0]
		if end > strLen {
			end = strLen
		}
	}

	return string(rs[start:end])
}

func Replace(str *string, old string, new string) string {
	return strings.Replace(*str, old, new, -1)
}

func Split(str *string, sep string) []string {
	return strings.Split(*str, sep)
}

func SnakeString(s string) string {
	data := make([]byte, 0, len(s)*2)
	num := len(s)
	for i := 0; i < num; i++ {
		d := s[i]
		if i > 0 && d >= 'A' && d <= 'Z' {
			data = append(data, '_')
		}
		data = append(data, d)
	}
	return strings.ToLower(string(data[:]))
}

func StrCut(str string, length int) string {
	rs := []rune(str)
	if len(rs) <= length {
		return str
	}

	return string(rs[0:length]) + "..."
}
