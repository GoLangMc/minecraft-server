package impl

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	uuid "github.com/satori/go.uuid"

	"minecraft-server/apis"
	"minecraft-server/apis/cmds"
	"minecraft-server/apis/data"
	"minecraft-server/apis/ents"
	"minecraft-server/apis/logs"
	"minecraft-server/apis/task"
	"minecraft-server/apis/util"
	"minecraft-server/impl/conf"
	"minecraft-server/impl/data/plugin"

	"minecraft-server/impl/conn"
	"minecraft-server/impl/cons"
	"minecraft-server/impl/data/system"
	"minecraft-server/impl/data/values"
	"minecraft-server/impl/prot"

	apis_base "minecraft-server/apis/base"
	impl_base "minecraft-server/impl/base"

	apis_event "minecraft-server/apis/game/event"
	impl_event "minecraft-server/impl/game/event"
)

type server struct {
	message chan system.Message

	console *cons.Console

	logging *logs.Logging
	tasking *task.Tasking
	watcher util.Watcher

	command *cmds.CommandManager

	network impl_base.Network
	packets impl_base.Packets

	players *playerAssociation
}

func NewServer(conf conf.ServerConfig) apis.Server {
	message := make(chan system.Message)

	console := cons.NewConsole(message)

	logging := logs.NewLogging("server", logs.EveryLevel...)
	tasking := task.NewTasking(values.MPT)
	watcher := util.NewWatcher()

	join := make(chan impl_base.PlayerAndConnection)
	quit := make(chan impl_base.PlayerAndConnection)

	packets := prot.NewPackets(tasking, join, quit)
	network := conn.NewNetwork(conf.Network.Host, conf.Network.Port, packets, join, quit)

	command := cmds.NewCommandManager()

	return &server{
		message: message,

		console: console,

		logging: logging,
		tasking: tasking,
		watcher: watcher,

		command: command,

		packets: packets,
		network: network,

		players: &playerAssociation{
			uuidToData: make(map[uuid.UUID]*ents.Player),

			connToUUID: make(map[impl_base.Connection]uuid.UUID),
			uuidToConn: make(map[uuid.UUID]impl_base.Connection),
		},
	}
}

func (s *server) Load() {
	apis.SetMinecraftServer(s)

	s.console.Load()
	s.command.Load()
	s.tasking.Load()
	s.network.Load()

	s.command.Register("stop", s.stopServerCommand)
	s.command.Register("time", func(sender ents.Sender, params []string) {
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
	s.command.Register("send", func(sender ents.Sender, params []string) {
		message := strings.Join(params, " ")

		for _, player := range s.Players() {
			player.SendMessage(message)
		}
	})

	s.watcher.SubAs(func(event apis_event.PlayerJoinEvent) {
		s.logging.InfoF("player %s logged in with uuid:%v", event.Player.Name(), event.Player.UUID())
	})
	s.watcher.SubAs(func(event apis_event.PlayerQuitEvent) {
		s.logging.InfoF("%s disconnected!", event.Player.Name())
	})

	s.watcher.SubAs(func(event impl_event.PlayerConnJoinEvent) {
		s.players.addData(event.Conn)

		s.watcher.PubAs(apis_event.PlayerJoinEvent{Player: event.Conn.Player})
	})
	s.watcher.SubAs(func(event impl_event.PlayerConnQuitEvent) {
		player := s.players.playerByConn(event.Conn.Connection)

		if player != nil {
			s.watcher.PubAs(apis_event.PlayerQuitEvent{Player: *player})
		}

		s.players.delData(event.Conn)
	})

	s.watcher.SubAs(func(event impl_event.PlayerPluginMessagePullEvent) {
		s.logging.DataF("received message on channel '%s' from player %s:%s", event.Channel, event.Conn.Name(), event.Conn.UUID())

		switch event.Channel {
		case plugin.CHANNEL_BRAND:
			s.logging.DataF("their client's brand is '%s'", event.Message.(*plugin.Brand).Name)
		}
	})

	go func() {
		for {
			// read input from console
			text := strings.Trim(<-s.console.IChannel, " ")
			if len(text) == 0 {
				continue
			}

			args := strings.Split(text, " ")
			if len(args) == 0 {
				continue
			}

			if command := s.command.Search(args[0]); command != nil {

				err := apis_base.Attempt(func() {
					(*command).Evaluate(s.console, args[1:])
				})

				if err != nil {
					s.logging.Fail(
						data.Red, "failed to evaluate ",
						data.DarkGray, "`",
						data.White, (*command).Name(),
						data.DarkGray, "`",
						data.Red, ": ", err.Error()[8:])
				}

				continue
			}

			s.console.SendMessage(text)
		}
	}()

	s.Wait()
}

func (s *server) Kill() {

	s.console.Kill()
	s.command.Kill()
	s.tasking.Kill()
	s.network.Kill()

	// push the stop message to the server exit channel
	s.message <- system.Make(system.STOP, "normal stop")
	close(s.message)

	s.logging.Info(data.DarkRed, "server stopped")
}

func (s *server) Wait() {
	// select over server commands channel
	select {
	case command := <-s.message:
		switch command.Command {
		// stop selecting when stop is received
		case system.STOP:
			return
		case system.FAIL:
			fmt.Printf("internal server error: %s\n", command.Message)
			return
		}
	}

	s.Wait()
}

func (s *server) Logging() *logs.Logging {
	return s.logging
}

func (s *server) Command() *cmds.CommandManager {
	return s.command
}

func (s *server) Tasking() *task.Tasking {
	return s.tasking
}

func (s *server) Watcher() util.Watcher {
	return s.watcher
}

func (s *server) Players() []ents.Player {
	players := make([]ents.Player, 0)

	for _, player := range s.players.uuidToData {
		players = append(players, *player)
	}

	return players
}

func (s *server) ConnByUUID(uuid uuid.UUID) impl_base.Connection {
	return s.players.uuidToConn[uuid]
}

func (s *server) PlayerByUUID(uuid uuid.UUID) ents.Player {
	return *s.players.uuidToData[uuid]
}

func (s *server) PlayerByConn(conn impl_base.Connection) ents.Player {
	uuid, con := s.players.connToUUID[conn]
	if !con {
		return nil
	}

	return s.PlayerByUUID(uuid)
}

func (s *server) ServerVersion() string {
	return "0.0.1-SNAPSHOT"
}

func (s *server) stopServerCommand(sender ents.Sender, params []string) {
	if _, ok := sender.(*cons.Console); !ok {
		s.logging.FailF("non console sender %s tried to stop the server", sender.Name())
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
		s.logging.Warn(data.Gold, "stopping server in ", data.Green, util.FormatTime(after))

		// schedule shutdown {after} seconds later
		s.tasking.AfterTime(after, time.Second, func(task *task.Task) {
			s.Kill()
		})

	}
}

type playerAssociation struct {
	uuidToData map[uuid.UUID]*ents.Player

	connToUUID map[impl_base.Connection]uuid.UUID
	uuidToConn map[uuid.UUID]impl_base.Connection
}

func (p *playerAssociation) addData(data impl_base.PlayerAndConnection) {
	p.uuidToData[data.Player.UUID()] = &data.Player

	p.connToUUID[data.Connection] = data.Player.UUID()
	p.uuidToConn[data.Player.UUID()] = data.Connection
}

func (p *playerAssociation) delData(data impl_base.PlayerAndConnection) {
	player := p.playerByConn(data.Connection)

	uuid := p.connToUUID[data.Connection]

	delete(p.connToUUID, data.Connection)
	delete(p.uuidToConn, uuid)

	if player != nil {
		delete(p.uuidToData, (*player).UUID())
	}
}

func (p *playerAssociation) playerByUUID(uuid uuid.UUID) *ents.Player {
	return p.uuidToData[uuid]
}

func (p *playerAssociation) playerByConn(conn impl_base.Connection) *ents.Player {
	uuid, con := p.connToUUID[conn]

	if !con {
		return nil
	}

	data, con := p.uuidToData[uuid]

	if !con {
		return nil
	}

	return data
}
