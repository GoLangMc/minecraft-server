package tags

type Typ int8

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
