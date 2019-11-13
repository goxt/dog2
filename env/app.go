package env

import "github.com/kataras/iris"

type DbConfig struct {
	Driver  string
	ConnStr string
}

const (
	// 文件名最大长度
	FileNameMax = 200
)

var (
	Application = iris.New()
)

var (
	FileUploadExt = map[string][]string{
		"package": {
			// "zip", "rar", "7z", "iso", "tar", "gz",
		},
		"word": {
			"doc", "docx", "wps", "wpt", "dot", "rtf", "txt", "dotx", "docm", "dotm",
		},
		"ppt": {
			"ppt", "pptx", "dps", "dpt", "pot", "pps", "pptm", "potx", "potm", "ppsx", "ppsm",
		},
		"excel": {
			"xls", "xlsx", "et", "ett", "xlt", "xlsx", "xlsm", "dbf", "csv", "prn", "dif", "xltx", "xltm",
		},
		"image": {
			"png", "jpg", "jpeg", "bmp", "gif",
		},
		"pdf": {
			"pdf", "tiff", "tif",
		},
		"video": {
			"mp4", "avi", "rm", "wmv", "mov", "asf", "flv", "rmvb", "swf",
		},
		"audio": {
			"mp3", "wav", "wma",
		},
		"html": {
			"xml", "htm", "html", "mht", "mhtml",
		},
	}
)
