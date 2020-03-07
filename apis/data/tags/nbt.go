package tags

type Typ byte

const (
	TAG_End Typ = iota
	TAG_Byte
	TAG_Short
	TAG_Int
	TAG_Long
	TAG_Float
	TAG_Double
	TAG_Byte_Array
	TAG_String
	TAG_List
	TAG_Compound
	TAG_Int_Array
	TAG_Long_Array
)

type Nbt interface {
	Type() Typ

	Name() string
}

// end
type NbtEnd struct{}

func (n *NbtEnd) Type() Typ {
	return TAG_End
}

func (n *NbtEnd) Name() string {
	return "TAG_End"
}

// byte
type NbtByt struct {
	Value int8
}

func (n *NbtByt) Type() Typ {
	return TAG_Byte
}

func (n *NbtByt) Name() string {
	return "TAG_Byte"
}

// short
type NbtI16 struct {
	Value int16
}

func (n *NbtI16) Type() Typ {
	return TAG_Short
}

func (n *NbtI16) Name() string {
	return "TAG_Short"
}

/*func (n *NbtI16) Push(writer buff.Buffer) {
	writer.PushI16(n.Value)
}

func (n *NbtI16) Pull(reader buff.Buffer) {
	n.Value = reader.PullI16()
}*/

// int
type NbtI32 struct {
	Value int32
}

func (n *NbtI32) Type() Typ {
	return TAG_Int
}

func (n *NbtI32) Name() string {
	return "TAG_Int"
}

/*func (n *NbtI32) Push(writer buff.Buffer) {
	writer.PushI32(n.Value)
}

func (n *NbtI32) Pull(reader buff.Buffer) {
	n.Value = reader.PullI32()
}*/

// long
type NbtI64 struct {
	Value int64
}

func (n *NbtI64) Type() Typ {
	return TAG_Long
}

func (n *NbtI64) Name() string {
	return "TAG_Long"
}

/*func (n *NbtI64) Push(writer buff.Buffer) {
	writer.PushI64(n.Value)
}

func (n *NbtI64) Pull(reader buff.Buffer) {
	n.Value = reader.PullI64()
}*/

// float
type NbtF32 struct {
	Value float32
}

func (n *NbtF32) Type() Typ {
	return TAG_Float
}

func (n *NbtF32) Name() string {
	return "TAG_Float"
}

/*func (n *NbtF32) Push(writer buff.Buffer) {
	writer.PushF32(n.Value)
}

func (n *NbtF32) Pull(reader buff.Buffer) {
	n.Value = reader.PullF32()
}*/

// double
type NbtF64 struct {
	Value float64
}

func (n *NbtF64) Type() Typ {
	return TAG_Double
}

func (n *NbtF64) Name() string {
	return "TAG_Double"
}

/*func (n *NbtF64) Push(writer buff.Buffer) {
	writer.PushF64(n.Value)
}

func (n *NbtF64) Pull(reader buff.Buffer) {
	n.Value = reader.PullF64()
}*/

// byte array
type NbtArrByt struct {
	Value []int8
}

func (n *NbtArrByt) Type() Typ {
	return TAG_Byte_Array
}

func (n *NbtArrByt) Name() string {
	return "TAG_Byte_Array"
}

/*func (n *NbtArrByt) Push(writer buff.Buffer) {
	writer.PushSAS(n.Value, true)
}

func (n *NbtArrByt) Pull(reader buff.Buffer) {
	n.Value = reader.PullSAS()
}*/

// string
type NbtTxt struct {
	Value string
}

func (n *NbtTxt) Type() Typ {
	return TAG_String
}

func (n *NbtTxt) Name() string {
	return "TAG_String"
}

/*func (n *NbtTxt) Push(writer buff.Buffer) {
	writer.PushTxt(n.Value)
}

func (n *NbtTxt) Pull(reader buff.Buffer) {
	n.Value = reader.PullTxt()
}*/

// typed list
type NbtArrAny struct {
	NType Typ
	Value []Nbt
}

func (n *NbtArrAny) Type() Typ {
	return TAG_List
}

func (n *NbtArrAny) Name() string {
	return "TAG_List"
}

/*func (n *NbtArrAny) Push(writer buff.Buffer) {
	if len(n.Value) == 0 {
		writer.PushByt(0)
	} else {
		writer.PushByt(byte(n.NType))
	}

	writer.PushI32(int32(len(n.Value)))

	for _, nbt := range n.Value {
		nbt.Push(writer)
	}
}

func (n *NbtArrAny) Pull(reader buff.Buffer) {
	nType := Typ(reader.PullByt()) // this can probably fail...

	size := reader.PullI32()
	value := make([]Nbt, size, size)

	for i := int32(0); i < size; i++ {
		inst := typeToInst[nType]()
		inst.Pull(reader)

		value[i] = inst
	}
}*/

// compound (map)
type NbtCompound struct {
	Named string
	Value map[string]Nbt
}

func (n *NbtCompound) Type() Typ {
	return TAG_Compound
}

func (n *NbtCompound) Name() string {
	return "TAG_Compound"
}

func (n *NbtCompound) Set(name string, data Nbt) {
	n.Value[name] = data
}

func (n *NbtCompound) Get(name string) (nbt Nbt, con bool) {
	nbt, con = n.Value[name]
	return
}

// int list
type NbtArrI32 struct {
	Value []int32
}

func (n *NbtArrI32) Type() Typ {
	return TAG_Int_Array
}

func (n *NbtArrI32) Name() string {
	return "TAG_Int_Array"
}

/*func (n *NbtArrI32) Push(writer buff.Buffer) {
	writer.PushI32(int32(len(n.Value)))

	for _, value := range n.Value {
		writer.PushI32(value)
	}
}

func (n *NbtArrI32) Pull(reader buff.Buffer) {
	value := make([]int32, reader.PullI32())

	for i := 0; i < len(value); i++ {
		value[i] = reader.PullI32()
	}
}*/

// long list
type NbtArrI64 struct {
	Value []int64
}

func (n *NbtArrI64) Type() Typ {
	return TAG_Long_Array
}

func (n *NbtArrI64) Name() string {
	return "TAG_Long_Array"
}
