package cmds

import (
	"strings"

	"github.com/golangmc/minecraft-server/apis/ents"
)

type CommandManager struct {
	commands map[string]*Command
}

func NewCommandManager() *CommandManager {
	return &CommandManager{
		commands: make(map[string]*Command),
	}
}

func (c *CommandManager) Load() {

}

func (c *CommandManager) Kill() {
	c.commands = nil
}

func (c *CommandManager) RegisterCommand(command Command) {
	c.commands[command.Name()] = &command

	command.Load()
}

func (c *CommandManager) Register(name string, evaluate func(sender ents.Sender, params []string)) {
	command := simpleCommand{
		name:     name,
		evaluate: evaluate,
	}

	c.RegisterCommand(&command)
}

func (c *CommandManager) Search(named string) *Command {
	for name, command := range c.commands {
		if strings.EqualFold(name, named) {
			return command
		}
	}

	return nil
}

type simpleCommand struct {
	name string

	evaluate func(sender ents.Sender, params []string)
}

func (s *simpleCommand) Name() string {
	return s.name
}

func (s *simpleCommand) Load() {

}

func (s *simpleCommand) Kill() {

}

func (s *simpleCommand) Evaluate(sender ents.Sender, params []string) {
	s.evaluate(sender, params)
}

func (s *simpleCommand) Complete(sender ents.Sender, params []string, output *[]string) {

}
