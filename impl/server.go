package impl

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"minecraft-server/apis/cmds"
	"minecraft-server/apis/data"
	"minecraft-server/apis/ents"
	"minecraft-server/apis/logs"
	"minecraft-server/apis/util"
	"minecraft-server/impl/conn"
	"minecraft-server/impl/cons"

	"minecraft-server/impl/data/server"
	"minecraft-server/impl/prot"
	"minecraft-server/impl/task"

	apis "minecraft-server/apis/base"
	impl "minecraft-server/impl/base"
)

type Server struct {
	message chan server.Message

	Console *cons.Console
	Logging *logs.Logging

	Tasking *task.Tasking

	Command *cmds.CommandManager

	Network impl.Network
	Packets impl.Packets
}

func NewServer(host string, port int) *Server {
	message := make(chan server.Message)

	console := cons.NewConsole(message)
	logging := logs.NewLogging("server", logs.EveryLevel...)

	tasking := task.NewTasking()

	join := make(chan impl.PlayerAndConnection)
	quit := make(chan impl.PlayerAndConnection)

	packets := prot.NewPackets(join, quit)
	network := conn.NewNetwork(host, port, packets, join, quit)

	command := cmds.NewCommandManager()

	return &Server{
		message: message,

		Console: console,

		Logging: logging,
		Tasking: tasking,

		Command: command,

		Packets: packets,
		Network: network,
	}
}

func (s *Server) Load() {

	s.Console.Load()
	s.Command.Load()
	s.Tasking.Load()
	s.Network.Load()

	s.Command.Register("stop", s.stopServerCommand)
	s.Command.Register("time", func(sender ents.Sender, params []string) {
		var seconds int64 = 0

		if len(params) > 0 {
			param, err := strconv.Atoi(params[0])

			if err != nil {
				panic(err)
			}

			if param <= 0 {
				panic(fmt.Errorf("value must be a positive whole number. [1..]"))
			}

			seconds = int64(param)
		}

		sender.SendMessage(util.FormatTime(seconds))
	})

	go func() {
		for {
			// read input from console
			text := strings.Trim(<-s.Console.IChannel, " ")
			if len(text) == 0 {
				continue
			}

			args := strings.Split(text, " ")
			if len(args) == 0 {
				continue
			}

			if command := s.Command.Search(args[0]); command != nil {

				err := apis.Attempt(func() {
					(*command).Evaluate(s.Console, args[1:])
				})

				if err != nil {
					s.Logging.Fail(
						data.Red, "failed to evaluate ",
						data.DarkGray, "`",
						data.White, (*command).Name(),
						data.DarkGray, "`",
						data.Red, ": ", err.Error()[8:])
				}

				continue
			}

			s.Console.SendMessage(text)
		}
	}()

	s.Wait()
}

func (s *Server) Kill() {

	s.Console.Kill()
	s.Command.Kill()
	s.Tasking.Kill()
	s.Network.Kill()

	// push the stop message to the server exit channel
	s.message <- server.Make(server.STOP, "normal stop")
	close(s.message)

	s.Logging.Info(data.DarkRed, "server stopped")
}

func (s *Server) Wait() {
	// select over server commands channel
	select {
	case command := <-s.message:
		switch command.Command {
		// stop selecting when stop is received
		case server.STOP:
			return
		case server.FAIL:
			fmt.Printf("internal server error: %s\n", command.Message)
			return
		}
	}

	s.Wait()
}

func (s *Server) stopServerCommand(sender ents.Sender, params []string) {
	if _, ok := sender.(*cons.Console); !ok {
		s.Logging.FailF("non console sender %s tried to stop the server", sender.Name())
		return
	}

	var after int64 = 0

	if len(params) > 0 {
		param, err := strconv.Atoi(params[0])

		if err != nil {
			panic(err)
		}

		if param <= 0 {
			panic(fmt.Errorf("value must be a positive whole number. [1..]"))
		}

		after = int64(param)
	}

	if after == 0 {

		s.Kill()

	} else {

		// inform future shutdown
		s.Logging.Warn(data.Gold, "stopping server in ", data.Green, util.FormatTime(after))

		// schedule shutdown {after} seconds later
		s.Tasking.AfterTime(after, time.Second, func(task *task.Task) {
			s.Kill()
		})

	}
}
