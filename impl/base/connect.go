package base

type Connection interface {
	GetState() PacketState
	SetState(state PacketState)

	Encrypt(data []byte) (output []byte)
	Decrypt(data []byte) (output []byte)

	CertifyName() string

	CertifyData() []byte

	CertifyValues(name string)
	CertifyUpdate(secret []byte)

	Deflate(data []byte) (output []byte)
	Inflate(data []byte) (output []byte)

	Pull(data []byte) (len int, err error)
	Push(data []byte) (len int, err error)

	Stop() (err error)

	SendPacket(packet PacketO)
}
