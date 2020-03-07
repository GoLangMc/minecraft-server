package uuid

import "github.com/satori/go.uuid"

type UUID = uuid.UUID

func NewUUID() UUID {
	gen, err := uuid.NewV4()
	if err != nil {
		panic(err)
	}

	return gen
}

func TextToUUID(text string) (data UUID, err error) {
	return uuid.FromString(text)
}

func UUIDToText(uuid UUID) (text string, err error) {
	data, err := uuid.MarshalText()

	if err == nil {
		text = string(data)
	}

	return
}
