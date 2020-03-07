package client

import (
	"github.com/golangmc/minecraft-server/apis/buff"
	"github.com/golangmc/minecraft-server/apis/data"
	"github.com/golangmc/minecraft-server/apis/data/msgs"
	"github.com/golangmc/minecraft-server/apis/ents"
	"github.com/golangmc/minecraft-server/apis/game"
	"github.com/golangmc/minecraft-server/apis/game/level"
	"github.com/golangmc/minecraft-server/impl/base"
	"github.com/golangmc/minecraft-server/impl/data/client"
	"github.com/golangmc/minecraft-server/impl/data/plugin"

	apis_conn "github.com/golangmc/minecraft-server/impl/conn"
)

type PacketOChatMessage struct {
	Message         msgs.Message
	MessagePosition msgs.MessagePosition
}

func (p *PacketOChatMessage) UUID() int32 {
	return 0x0F
}

func (p *PacketOChatMessage) Push(writer buff.Buffer, conn base.Connection) {
	message := p.Message

	if p.MessagePosition == msgs.HotBarText {
		message = *msgs.New(message.AsText())
	}

	writer.PushTxt(message.AsJson())
	writer.PushByt(byte(p.MessagePosition))
}

type PacketOJoinGame struct {
	EntityID      int32
	Hardcore      bool
	GameMode      game.GameMode
	Dimension     game.Dimension
	HashedSeed    int64
	MaxPlayers    int
	LevelType     game.LevelType
	ViewDistance  int32
	ReduceDebug   bool
	RespawnScreen bool
}

func (p *PacketOJoinGame) UUID() int32 {
	return 0x26
}

func (p *PacketOJoinGame) Push(writer buff.Buffer, conn base.Connection) {
	writer.PushI32(p.EntityID)
	writer.PushByt(p.GameMode.Encoded(p.Hardcore /* pull this value from somewhere */))
	writer.PushI32(int32(p.Dimension))
	writer.PushI64(p.HashedSeed)
	writer.PushByt(uint8(p.MaxPlayers))
	writer.PushTxt(p.LevelType.String())
	writer.PushVrI(p.ViewDistance)
	writer.PushBit(p.ReduceDebug)
	writer.PushBit(p.RespawnScreen)
}

type PacketOPluginMessage struct {
	Message plugin.Message
}

func (p *PacketOPluginMessage) UUID() int32 {
	return 0x19
}

func (p *PacketOPluginMessage) Push(writer buff.Buffer, conn base.Connection) {
	writer.PushTxt(p.Message.Chan())
	p.Message.Push(writer)
}

type PacketOPlayerLocation struct {
	Location data.Location
	Relative client.Relativity

	ID int32
}

func (p *PacketOPlayerLocation) UUID() int32 {
	return 0x36
}

func (p *PacketOPlayerLocation) Push(writer buff.Buffer, conn base.Connection) {
	writer.PushF64(p.Location.X)
	writer.PushF64(p.Location.Y)
	writer.PushF64(p.Location.Z)

	writer.PushF32(p.Location.AxisX)
	writer.PushF32(p.Location.AxisY)

	p.Relative.Push(writer)

	writer.PushVrI(p.ID)
}

type PacketOKeepAlive struct {
	KeepAliveID int64
}

func (p *PacketOKeepAlive) UUID() int32 {
	return 0x21
}

func (p *PacketOKeepAlive) Push(writer buff.Buffer, conn base.Connection) {
	writer.PushI64(p.KeepAliveID)
}

type PacketOServerDifficulty struct {
	Difficulty game.Difficulty
	Locked     bool // should probably always be true
}

func (p *PacketOServerDifficulty) UUID() int32 {
	return 0x0E
}

func (p *PacketOServerDifficulty) Push(writer buff.Buffer, conn base.Connection) {
	writer.PushByt(byte(p.Difficulty))
	writer.PushBit(p.Locked)
}

type PacketOPlayerAbilities struct {
	Abilities   client.PlayerAbilities
	FlyingSpeed float32
	FieldOfView float32
}

func (p *PacketOPlayerAbilities) UUID() int32 {
	return 0x32
}

func (p *PacketOPlayerAbilities) Push(writer buff.Buffer, conn base.Connection) {
	p.Abilities.Push(writer)

	writer.PushF32(p.FlyingSpeed)
	writer.PushF32(p.FieldOfView)
}

type PacketOHeldItemChange struct {
	Slot client.HotBarSlot
}

func (p *PacketOHeldItemChange) UUID() int32 {
	return 0x40
}

func (p *PacketOHeldItemChange) Push(writer buff.Buffer, conn base.Connection) {
	writer.PushByt(byte(p.Slot))
}

type PacketODeclareRecipes struct {
	// Recipes []*Recipe // this doesn't exist yet ;(
	RecipeCount int32
}

func (p *PacketODeclareRecipes) UUID() int32 {
	return 0x5B
}

func (p *PacketODeclareRecipes) Push(writer buff.Buffer, conn base.Connection) {
	writer.PushVrI(p.RecipeCount)
	// when recipes are implemented, instead of holding a recipe count, simply write the size of the slice, Recipe will implement BufferPush
}

type PacketOChunkData struct {
	Chunk level.Chunk
}

func (p *PacketOChunkData) UUID() int32 {
	return 0x22
}

func (p *PacketOChunkData) Push(writer buff.Buffer, conn base.Connection) {
	writer.PushI32(int32(p.Chunk.ChunkX()))
	writer.PushI32(int32(p.Chunk.ChunkZ()))

	// full chunk (for now >:D)
	writer.PushBit(true)

	chunkData := apis_conn.NewBuffer()
	p.Chunk.Push(chunkData) // write chunk data and primary bit mask

	// pull primary bit mask and push to writer
	writer.PushVrI(chunkData.PullVrI())

	// write height-maps
	writer.PushNbt(p.Chunk.HeightMapNbtCompound())

	biomes := make([]int32, 1024, 1024)
	for i := range biomes {
		biomes[i] = 0 // void biome
	}

	for _, biome := range biomes {
		writer.PushI32(biome)
	}

	// data, prefixed with len
	writer.PushUAS(chunkData.UAS(), true)

	// write block entities
	writer.PushVrI(0)
}

type PacketOPlayerInfo struct {
	Action client.PlayerInfoAction
	Values []client.PlayerInfo
}

func (p *PacketOPlayerInfo) UUID() int32 {
	return 0x34
}

func (p *PacketOPlayerInfo) Push(writer buff.Buffer, conn base.Connection) {
	writer.PushVrI(int32(p.Action))
	writer.PushVrI(int32(len(p.Values)))

	for _, value := range p.Values {
		value.Push(writer)
	}
}

type PacketOEntityMetadata struct {
	Entity ents.Entity
}

func (p *PacketOEntityMetadata) UUID() int32 {
	return 0x44
}

func (p *PacketOEntityMetadata) Push(writer buff.Buffer, conn base.Connection) {
	writer.PushVrI(int32(p.Entity.EntityUUID())) // questionable...

	// only supporting player metadata for now
	_, ok := p.Entity.(ents.Player)
	if ok {

		writer.PushByt(16) // index | displayed skin parts
		writer.PushVrI(0)  // type | byte

		skin := client.SkinParts{
			Cape: true,
			Head: true,
			Body: true,
			ArmL: true,
			ArmR: true,
			LegL: true,
			LegR: true,
		}

		skin.Push(writer)
	}

	writer.PushByt(0xFF)
}
