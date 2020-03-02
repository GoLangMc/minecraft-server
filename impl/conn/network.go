package conn

import (
	"fmt"
	"net"
	"reflect"
	"strconv"

	"minecraft-server/apis/logs"
	"minecraft-server/impl/base"
	"minecraft-server/impl/data/server"
)

type network struct {
	host string
	port int

	logger  *logs.Logging
	packets base.Packets

	join chan base.PlayerAndConnection
	quit chan base.PlayerAndConnection

	report chan server.Message
}

func NewNetwork(host string, port int, packet base.Packets, join chan base.PlayerAndConnection, quit chan base.PlayerAndConnection) base.Network {
	return &network{
		host: host,
		port: port,

		join: join,
		quit: quit,

		logger:  logs.NewLogging("network", logs.EveryLevel...),
		packets: packet,
	}
}

func (n *network) Load() {
	if err := n.startListening(); err != nil {
		n.report <- server.Make(server.FAIL, err)
		return
	}
}

func (n *network) Kill() {

}

func (n *network) startListening() error {
	ser, err := net.ResolveTCPAddr("tcp", n.host+":"+strconv.Itoa(n.port))
	if err != nil {
		return fmt.Errorf("resolution failed [%v]", err)
	}

	tcp, err := net.ListenTCP("tcp", ser)
	if err != nil {
		return fmt.Errorf("binding failed [%v]", err)
	}

	n.logger.InfoF("listening on %s:%d", n.host, n.port)

	go func() {
		for {
			con, err := tcp.AcceptTCP()

			if err != nil {
				n.report <- server.Make(server.FAIL, err)
				break
			}

			_ = con.SetNoDelay(true)
			_ = con.SetKeepAlive(true)

			go handleConnect(n, NewConnection(con))
		}
	}()

	return nil
}

func handleConnect(network *network, conn base.Connection) {
	// network.logger.Data("New Connection from &6%v", conn.TCP.RemoteAddr())

	var inf []byte

	for {
		inf = make([]byte, 1024)
		sze, err := conn.Pull(inf)

		if err != nil && err.Error() == "EOF" {
			network.quit <- base.PlayerAndConnection{
				Player:     nil,
				Connection: conn,
			}

			break
		}

		if err != nil || sze == 0 {
			_ = conn.Stop()
			break
		}

		buf := NewBufferWith(conn.Decrypt(inf[:sze]))

		// decompression
		// decryption

		if buf.Arr()[0] == 0xFE { // LEGACY PING
			continue
		}

		packetLen := buf.PullVrI()

		bufI := NewBufferWith(buf.Arr()[buf.InI() : buf.InI()+packetLen])
		bufO := NewBuffer()

		handleReceive(network, conn, bufI, bufO)

		if bufO.Len() > 1 {
			temp := NewBuffer()
			temp.PushVrI(bufO.Len())

			comp := NewBuffer()
			comp.PushArr(conn.Deflate(bufO.Arr()), false)

			temp.PushArr(comp.Arr(), false)

			_, err := conn.Push(conn.Encrypt(temp.Arr()))

			if err != nil {
				network.logger.Fail("Failed to push client bound packet: %v", err)
			}
		}
	}
}

func handleReceive(network *network, conn base.Connection, bufI base.Buffer, bufO base.Buffer) {
	uuid := bufI.PullVrI()

	network.logger.InfoF("UUID: %d", uuid)

	packetI := network.packets.GetPacketI(uuid, conn.GetState())
	if packetI == nil {
		network.logger.WarnF("unable to decode %v packet with uuid: %d", conn.GetState(), uuid)
		return
	}

	network.logger.DataF("GET packet: %d | %v | %v", packetI.UUID(), reflect.TypeOf(packetI), conn.GetState())

	// populate incoming packet
	packetI.Pull(bufI, conn)

	network.packets.PubAs(packetI)
	network.packets.PubAs(packetI, conn)
}
