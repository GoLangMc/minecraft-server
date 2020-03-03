package apis

import (
	"minecraft-server/apis/base"
	"minecraft-server/apis/cmds"
	"minecraft-server/apis/logs"
	"minecraft-server/apis/task"
)

type Server interface {
	base.State

	Logging() *logs.Logging

	Command() *cmds.CommandManager

	Tasking() *task.Tasking
}
