package client

type StatusAction int

const (
	Respawn StatusAction = iota
	Request
)
