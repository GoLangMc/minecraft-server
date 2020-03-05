package mode

import (
	"bytes"
	"fmt"

	"minecraft-server/apis/data"
	"minecraft-server/apis/util"
	"minecraft-server/apis/uuid"
	"minecraft-server/impl/base"
	"minecraft-server/impl/game/auth"
	"minecraft-server/impl/game/ents"
	"minecraft-server/impl/prot/states"
)

/**
 * login
 */

func HandleState2(watcher util.Watcher, join chan base.PlayerAndConnection) {

	watcher.SubAs(func(packet *states.PacketILoginStart, conn base.Connection) {
		conn.CertifyValues(packet.PlayerName)

		_, public := auth.NewCrypt()

		response := states.PacketOEncryptionRequest{
			Server: "",
			Public: public,
			Verify: conn.CertifyData(),
		}

		conn.SendPacket(&response)
	})

	watcher.SubAs(func(packet *states.PacketIEncryptionResponse, conn base.Connection) {
		defer func() {
			if err := recover(); err != nil {
				conn.SendPacket(&states.PacketODisconnect{
					Reason: data.NewMessage(fmt.Sprintf("Authentication failed: %v", err)).SetColor(data.Red),
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
					conn.SendPacket(&states.PacketODisconnect{
						Reason: data.NewMessage(fmt.Sprintf("Authentication failed: %v", err)).SetColor(data.Red),
					})
				}
			}()

			if err != nil {
				panic(fmt.Errorf("failed to authenticate: %s\n%v\n", conn.CertifyName(), err))
			}

			uuid, err := uuid.TextToUUID(auth.UUID)
			if err != nil {
				panic(fmt.Errorf("failed to decode uuid for %s: %s\n%v\n", conn.CertifyName(), auth.UUID, err))
			}

			player := ents.NewPlayer()

			// store the player object somewhere

			player.SetUUID(uuid)
			player.SetName(auth.Name)

			conn.SendPacket(&states.PacketOLoginSuccess{
				PlayerName: player.Name(),
				PlayerUUID: player.UUID().String(),
			})

			conn.SetState(base.PLAY)

			join <- base.PlayerAndConnection{
				Player:     &player,
				Connection: conn,
			}
		})

	})

}
