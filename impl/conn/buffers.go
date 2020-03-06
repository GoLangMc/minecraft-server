package conn

import (
	"encoding/binary"
	"fmt"
	"math"
	"strconv"

	"minecraft-server/apis/data"
	"minecraft-server/apis/data/tags"
	"minecraft-server/impl/base"
)

/*
 Language used:
	- Len = Length
	- Arr = Array

	- Bit = Boolean
	- Byt = Byte
	- Int = int
	- VrI = VarInt
	- Srt = Short
	- Txt = String
*/

type buffer struct {
	iIndex int32
	oIndex int32

	bArray []byte
}

func (b *buffer) String() string {
	return fmt.Sprintf("Buffer[%d](i: %d, o: %d)%v", b.Len(), b.iIndex, b.oIndex, b.bArray)
}

// new
func NewBuffer() base.Buffer {
	return NewBufferWith(make([]byte, 0))
}

func NewBufferWith(bArray []byte) base.Buffer {
	return &buffer{bArray: bArray}
}

// server_data
func (b *buffer) Len() int32 {
	return int32(len(b.bArray))
}

func (b *buffer) SAS() []int8 {
	return asSArray(b.bArray)
}

func (b *buffer) UAS() []byte {
	return b.bArray
}

func (b *buffer) InI() int32 {
	return b.iIndex
}

func (b *buffer) InO() int32 {
	return b.oIndex
}

func (b *buffer) SkpAll() {
	b.SkpLen(b.Len() - 1)
}

func (b *buffer) SkpLen(delta int32) {
	b.iIndex += delta
}

// pull
func (b *buffer) PullBit() bool {
	return b.pullNext() != 0
}

func (b *buffer) PullByt() byte {
	return b.pullNext()
}

func (b *buffer) PullI16() int16 {
	return int16(binary.BigEndian.Uint16(b.pullSize(4)))
}

func (b *buffer) PullU16() uint16 {
	return uint16(b.pullNext())<<8 | uint16(b.pullNext())
}

func (b *buffer) PullI32() int32 {
	return int32(binary.BigEndian.Uint32(b.pullSize(4)))
}

func (b *buffer) PullI64() int64 {
	return int64(b.PullU64())
}

func (b *buffer) PullU64() uint64 {
	return binary.BigEndian.Uint64(b.pullSize(8))
}

func (b *buffer) PullF32() float32 {
	return math.Float32frombits(binary.BigEndian.Uint32(b.pullSize(4)))
}

func (b *buffer) PullF64() float64 {
	return math.Float64frombits(binary.BigEndian.Uint64(b.pullSize(8)))
}

func (b *buffer) PullVrI() int32 {
	return int32(b.pullVariable(5))
}

func (b *buffer) PullVrL() int64 {
	return b.pullVariable(10)
}

func (b *buffer) PullTxt() string {
	return string(b.PullUAS())
}

func (b *buffer) PullUAS() []byte {
	sze := b.PullVrI()
	arr := b.bArray[b.iIndex : b.iIndex+sze]

	b.iIndex += sze

	return arr
}

func (b *buffer) PullSAS() []int8 {
	return asSArray(b.PullUAS())
}

func (b *buffer) PullPos() data.PositionI {
	val := b.PullU64()

	x := int64(val) >> 38
	y := int64(val) & 0xFFF
	z := int64(val) << 26 >> 38

	return data.PositionI{
		X: x,
		Y: y,
		Z: z,
	}
}

func (b *buffer) PullNbt() tags.Nbt {
	panic("implement me")
}

// push
func (b *buffer) PushBit(data bool) {
	if data {
		b.pushNext(byte(0x01))
	} else {
		b.pushNext(byte(0x00))
	}
}

func (b *buffer) PushByt(data byte) {
	b.pushNext(data)
}

func (b *buffer) PushI16(data int16) {
	b.pushNext(
		byte(data)>>8,
		byte(data))
}

func (b *buffer) PushI32(data int32) {
	b.pushNext(
		byte(data>>24),
		byte(data>>16),
		byte(data>>8),
		byte(data))
}

func (b *buffer) PushI64(data int64) {
	b.pushNext(
		byte(data>>56),
		byte(data>>48),
		byte(data>>40),
		byte(data>>32),
		byte(data>>24),
		byte(data>>16),
		byte(data>>8),
		byte(data))
}

func (b *buffer) PushF32(data float32) {
	b.PushI32(int32(math.Float32bits(data)))
}

func (b *buffer) PushF64(data float64) {
	b.PushI64(int64(math.Float64bits(data)))
}

func (b *buffer) PushVrI(data int32) {
	for {
		temp := data & 0x7F
		data >>= 7

		if data != 0 {
			temp |= 0x80
		}

		b.pushNext(byte(temp))

		if data == 0 {
			break
		}
	}
}

func (b *buffer) PushVrL(data int64) {
	for {
		temp := data & 0x7F
		data >>= 7

		if data != 0 {
			temp |= 0x80
		}

		b.pushNext(byte(temp))

		if data == 0 {
			break
		}
	}
}

func (b *buffer) PushTxt(data string) {
	b.PushUAS([]byte(data), true)
}

func (b *buffer) PushUAS(data []byte, prefixWithLen bool) {
	if prefixWithLen {
		b.PushVrI(int32(len(data)))
	}

	b.pushNext(data...)
}

func (b *buffer) PushSAS(data []int8, prefixWithLen bool) {
	b.PushUAS(asUArray(data), prefixWithLen)
}

func (b *buffer) PushPos(data data.PositionI) {
	b.PushI64(((data.X & 0x3FFFFFF) << 38) | ((data.Z & 0x3FFFFFF) << 12) | (data.Y & 0xFFF))
}

func (b *buffer) PushNbt(data tags.Nbt) {
	panic("implement me")
}

// internal
func (b *buffer) pullNext() byte {

	if b.iIndex >= b.Len() {
		return 0
		// panic("reached end of buffer")
	}

	next := b.bArray[b.iIndex]
	b.iIndex++

	if b.oIndex > 0 {
		b.oIndex--
	}

	return next
}

func (b *buffer) pullSize(next int) []byte {
	bytes := make([]byte, next)

	for i := 0; i < next; i++ {
		bytes[i] = b.pullNext()
	}

	return bytes
}

func (b *buffer) pushNext(bArray ...byte) {
	b.oIndex += int32(len(bArray))
	b.bArray = append(b.bArray, bArray...)
}

func (b *buffer) pullVariable(max int) int64 {
	var num int
	var res int64

	for {
		tmp := int64(b.pullNext())
		res |= (tmp & 0x7F) << uint(num*7)

		if num++; num > max {
			panic("VarInt > " + strconv.Itoa(max))
		}

		if tmp&0x80 != 0x80 {
			break
		}
	}

	return res
}

func asSArray(bytes []byte) []int8 {
	array := make([]int8, 0)

	for _, b := range bytes {
		array = append(array, int8(b))
	}

	return array
}

func asUArray(bytes []int8) []byte {
	array := make([]byte, 0)

	for _, b := range bytes {
		array = append(array, byte(b))
	}

	return array
}
