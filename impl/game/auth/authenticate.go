package auth

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/golangmc/minecraft-server/apis/urls"
)

const url = "https://sessionserver.mojang.com/session/minecraft/hasJoined"

type Auth struct {
	UUID string `json:"id"`
	Name string `json:"name"`
	Prop []Prop `json:"properties"`
}

type Prop struct {
	Name string  `json:"name"`
	Data string  `json:"value"`
	Sign *string `json:"signature"`
}

func RunAuthGet(secret []byte, name string, callback func(auth *Auth, err error)) {
	go execute(generateAuthURL(name, generateAuthSHA(secret)), callback)
}

func execute(url string, callback func(auth *Auth, err error)) {
	get, err := urls.GetByte(url)
	if err != nil {
		callback(nil, err)
		return
	}

	var auth Auth

	err = json.Unmarshal(get, &auth)

	if err != nil {
		callback(nil, err)
	} else {
		callback(&auth, nil)
	}
}

func generateAuthURL(name, hash string) string {
	return fmt.Sprintf("%s?username=%s&serverId=%s", url, name, hash)
}

func generateAuthSHA(secret []byte) string {
	sha := sha1.New()

	// update with encoded secret, and encoded public
	_, public := NewCrypt()

	sha.Write(secret)
	sha.Write(public)

	hash := sha.Sum(nil)

	// Check for negative hashes
	negative := (hash[0] & 0x80) == 0x80

	if negative {
		carry := true

		for i := len(hash) - 1; i >= 0; i-- {
			hash[i] = ^hash[i]
			if carry {
				carry = hash[i] == 0xff
				hash[i]++
			}
		}
	}

	// Trim away zeroes
	res := strings.TrimLeft(fmt.Sprintf("%x", hash), "0")
	if negative {
		res = "-" + res
	}

	return res
}
