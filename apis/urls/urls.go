package urls

import (
	"io/ioutil"
	"net/http"
)

func GetByte(url string) (res []byte, err error) {
	// get
	out, err := http.Get(url)
	if err != nil {
		return
	}

	// read body
	bdy, err := ioutil.ReadAll(out.Body)
	if err != nil {
		return
	}

	// assign response
	res = bdy

	return
}

func GetText(url string) (res string, err error) {
	arr, err := GetByte(url)
	if err != nil {
		return
	}

	res = string(arr)

	return
}
