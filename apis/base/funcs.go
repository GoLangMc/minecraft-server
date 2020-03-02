package base

import (
	"fmt"
	"strings"
)

func ConvertToString(data ...interface{}) string {
	strs := make([]string, len(data))

	for i, str := range data {
		strs[i] = fmt.Sprintf("%v", str)
	}

	return strings.Join(strs, "")
}

func Attempt(function func()) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("caught: %v", r)
		}
	}()

	function()

	return
}
