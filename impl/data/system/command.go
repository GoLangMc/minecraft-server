package system

type Command int

type Message struct {
	Command
	Message interface{}
}

const (
	// stops the server entirely
	STOP Command = iota
	FAIL
)

func Make(command Command, message interface{}) Message {
	return struct {
		Command
		Message interface{}
	}{Command: command, Message: message}
}
