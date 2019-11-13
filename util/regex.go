package util

import (
	"regexp"
)

func Match(pattern string, str string) bool {
	ok, err := regexp.Match(pattern, []byte(str))
	if err != nil {
		panic(err)
	}
	return ok
}
