package game

type LevelType int

const (
	DEFAULT LevelType = iota
	FLAT
	LARGEBIOMES
	AMPLIFIED
	CUSTOMIZED
	BUFFET
	DEFAULT11
)

var typeToName = map[LevelType]string{
	DEFAULT:     "default",
	FLAT:        "flat",
	LARGEBIOMES: "largeBiomes",
	AMPLIFIED:   "amplified",
	CUSTOMIZED:  "customized",
	BUFFET:      "buffet",
	DEFAULT11:   "default_1_1",
}

func (l LevelType) String() string {
	return typeToName[l]
}
