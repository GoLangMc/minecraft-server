package base

import (
	"crypto/sha256"
	"encoding/binary"
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

func JavaStringHashCode(value string) int32 {
	var h int32

	if len(value) > 0 {
		for _, r := range value {
			h = 31*h + r
		}
	}

	return h
}

func JavaSHA256HashLong(value int64) []byte {
	bytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(bytes, uint64(value))

	hash := sha256.Sum256(bytes)

	return hash[:]
}
