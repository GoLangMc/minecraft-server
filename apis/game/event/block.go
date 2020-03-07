package event

import (
	"github.com/golangmc/minecraft-server/apis/game/level"
)

type BlockEvent struct {
	level.Block
}

type BlockBreakEvent struct {
	BlockEvent
	PlayerEvent
	Cancellable
}
