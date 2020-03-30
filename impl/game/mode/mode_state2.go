package mode

import (
	"bytes"
	"fmt"

	"github.com/golangmc/minecraft-server/apis/data/chat"
	"github.com/golangmc/minecraft-server/apis/data/msgs"
	"github.com/golangmc/minecraft-server/apis/game"
	"github.com/golangmc/minecraft-server/apis/util"
	"github.com/golangmc/minecraft-server/apis/uuid"
	"github.com/golangmc/minecraft-server/impl/base"
	"github.com/golangmc/minecraft-server/impl/game/auth"
	"github.com/golangmc/minecraft-server/impl/game/ents"
	"github.com/golangmc/minecraft-server/impl/prot/client"
	"github.com/golangmc/minecraft-server/impl/prot/server"
)

/**
 * login
 */

func HandleState2(watcher util.Watcher, join chan base.PlayerAndConnection, offline bool) {

	watcher.SubAs(func(packet *server.PacketILoginStart, conn base.Connection) {
		conn.CertifyValues(packet.PlayerName)

		_, public := auth.NewCrypt()

		if offline {
			addNewPlayer(conn, join, &auth.Auth{
				Name: conn.CertifyName(),
			})
			return
		}

		conn.SendPacket(&client.PacketOEncryptionRequest{
			Server: "",
			Public: public,
			Verify: conn.CertifyData(),
		})
	})

	watcher.SubAs(func(packet *server.PacketIEncryptionResponse, conn base.Connection) {
		defer func() {
			if err := recover(); err != nil {
				conn.SendPacket(&client.PacketODisconnect{
					Reason: *msgs.New(fmt.Sprintf("Authentication failed: %v", err)).SetColor(chat.Red),
				})
			}
		}()

		ver, err := auth.Decrypt(packet.Verify)
		if err != nil {
			panic(fmt.Errorf("failed to decrypt token: %s\n%v\n", conn.CertifyName(), err))
		}

		if !bytes.Equal(ver, conn.CertifyData()) {
			panic(fmt.Errorf("encryption failed, tokens are different: %s\n%v | %v", conn.CertifyName(), ver, conn.CertifyData()))
		}

		sec, err := auth.Decrypt(packet.Secret)
		if err != nil {
			panic(fmt.Errorf("failed to decrypt secret: %s\n%v\n", conn.CertifyName(), err))
		}

		conn.CertifyUpdate(sec) // enable encryption on the connection

		auth.RunAuthGet(sec, conn.CertifyName(), func(auth *auth.Auth, err error) {
			defer func() {
				if err := recover(); err != nil {
					conn.SendPacket(&client.PacketODisconnect{
						Reason: *msgs.New(fmt.Sprintf("Authentication failed: %v", err)).SetColor(chat.Red),
					})
				}
			}()

			if err != nil {
				panic(fmt.Errorf("failed to authenticate: %s\n%v\n", conn.CertifyName(), err))
			}

			addNewPlayer(conn, join, auth)
		})

	})

}

// addNewPlayer
// 1. reads an authentication response
// 2. sends a Login Success packet to the client
// 3. adds the resultant player to our game
func addNewPlayer(conn base.Connection, join chan base.PlayerAndConnection, auth *auth.Auth) {
	prof := game.Profile{
		UUID: uuid.TextToUUID(auth.UUID),
		Name: auth.Name,
	}

	for _, prop := range auth.Prop {
		prof.Properties = append(prof.Properties, &game.ProfileProperty{
			Name:      prop.Name,
			Value:     prop.Data,
			Signature: prop.Sign,
		})
	}

	player := ents.NewPlayer(&prof, conn)

	conn.SendPacket(&client.PacketOLoginSuccess{
		PlayerName: player.Name(),
		PlayerUUID: player.UUID().String(),
	})

	conn.SetState(base.PLAY)

	join <- base.PlayerAndConnection{
		Player:     player,
		Connection: conn,
	}
}
