package base

type Compacter struct {
	Values []int64

	bpb int
	max int
}

func NewCompacter(bits, size int) *Compacter {
	sliceSize := iDontKnowWhatThisDoes(size*bits, 64) / 64

	return &Compacter{
		bpb:    bits,
		max:    (1 << bits) - 1,
		Values: make([]int64, sliceSize, sliceSize),
	}
}

func (c *Compacter) Set(index int, value int) int {
	bIndex := index * c.bpb

	sIndex := bIndex >> 0x06
	eIndex := (((index + 1) * c.bpb) - 1) >> 0x06

	uIndex := bIndex ^ (sIndex << 0x06)

	previousValue := int64(uint64(c.Values[sIndex])>>uIndex) & int64(c.max)

	c.Values[sIndex] = c.Values[sIndex]&int64(^(c.max<<uIndex)) | int64((value&c.max)<<uIndex)

	if sIndex != eIndex {
		zIndex := 64 - uIndex
		pIndex := c.bpb - 1

		previousValue |= (c.Values[eIndex] << zIndex) & int64(c.max)

		c.Values[eIndex] = int64(((uint64(c.Values[eIndex]) >> pIndex) << pIndex) | uint64((value&c.max)>>zIndex))
	}

	return int(previousValue)
}

func (c *Compacter) Get(index int) int {
	bIndex := index * c.bpb

	sIndex := bIndex >> 0x06
	eIndex := (((index + 1) * c.bpb) - 1) >> 0x06

	uIndex := bIndex ^ (sIndex << 0x06)

	if sIndex == eIndex {
		return int((uint64(c.Values[sIndex]) >> uIndex) & uint64(c.max))
	}

	zIndex := 64 - uIndex

	return int(uint64(c.Values[sIndex]>>uIndex) | uint64(c.Values[eIndex]<<zIndex)&uint64(c.max))
}

func iDontKnowWhatThisDoes(var0, var1 int) int {
	if var1 == 0 {
		return 0
	} else if var0 == 0 {
		return var1
	} else {
		if var0 < 0 {
			var1 *= -1
		}

		var2 := var0 % var1

		if var2 == 0 {
			return var0
		} else {
			return var0 + var1 - var2
		}
	}
}
