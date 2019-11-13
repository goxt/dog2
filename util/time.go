package util

import (
	"time"
)

/**
 * 时间戳转成易读格式：x月x日 时:分
 */
func UnixTimeToMdHi(ut int64) string {
	var t = time.Unix(ut, 0)
	return t.Format("01月02日15:04")
}
