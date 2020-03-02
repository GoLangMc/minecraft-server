package rand

import "crypto/rand"

func RandomByteArray(len int) []byte {
	array := make([]byte, 4)
	_, _ = rand.Read(array)

	return array
}
