package plugin

import (
	"minecraft-server/apis/data"
	"minecraft-server/impl/base"
)

type Message interface {
	Chan() string

	base.BufferPush
	base.BufferPull
}

var registry = createMessageRegistry()

type MessageRegistry struct {
	channels map[string]func() Message
}

func createMessageRegistry() MessageRegistry {
	registry := MessageRegistry{make(map[string]func() Message)}

	registry.channels["minecraft:brand"] = func() Message {
		return &Brand{}
	}

	registry.channels["minecraft:debug/paths"] = func() Message {
		return &DebugPaths{}
	}

	registry.channels["minecraft:debug/neighbors_update"] = func() Message {
		return &DebugNeighbors{}
	}

	return registry
}

func GetMessageForChannel(channel string) Message {
	creator := registry.channels[channel]
	if creator == nil {
		return nil
	}

	return creator()
}

// look, they're like cute little packets :D

type Brand struct {
	Name string
}

func (b *Brand) Chan() string {
	return "minecraft:brand"
}

func (b *Brand) Push(writer base.Buffer) {
	writer.PushTxt(b.Name)
}

func (b *Brand) Pull(reader base.Buffer) {
	b.Name = reader.PullTxt()
}

type DebugPaths struct { // unused? honestly why did I do this
	UnknownValue1 int32
	UnknownValue2 float32
	PathEntity    PathEntity
}

type PathEntity struct {
	Index   int
	Target  PathPoint
	PSetLen int
	PSet    []PathPoint
	OSetLen int
	OSet    []PathPoint
	CSetLen int
	CSet    []PathPoint
}

func (p *PathEntity) Push(writer base.Buffer) {
	writer.PushI32(int32(p.Index))

	p.Target.Push(writer)

	writer.PushI32(int32(p.PSetLen))
	for _, point := range p.PSet {
		point.Push(writer)
	}

	writer.PushI32(int32(p.OSetLen))
	for _, point := range p.OSet {
		point.Push(writer)
	}

	writer.PushI32(int32(p.CSetLen))
	for _, point := range p.CSet {
		point.Push(writer)
	}
}

func (p *PathEntity) Pull(reader base.Buffer) {
	p.Index = int(reader.PullI32())

	target := PathPoint{}
	target.Pull(reader)

	p.Target = target

	p.PSet = make([]PathPoint, 0)
	p.PSetLen = int(reader.PullI32())

	for i := 0; i < p.PSetLen; i++ {
		point := PathPoint{}
		point.Pull(reader)

		p.PSet = append(p.PSet, point)
	}

	p.OSet = make([]PathPoint, 0)
	p.OSetLen = int(reader.PullI32())

	for i := 0; i < p.OSetLen; i++ {
		point := PathPoint{}
		point.Pull(reader)

		p.OSet = append(p.OSet, point)
	}

	p.CSet = make([]PathPoint, 0)
	p.CSetLen = int(reader.PullI32())

	for i := 0; i < p.CSetLen; i++ {
		point := PathPoint{}
		point.Pull(reader)

		p.CSet = append(p.CSet, point)
	}
}

type PathPoint struct {
	X int32
	Y int32
	Z int32

	DistanceOrigin float32
	Cost           float32
	CostMalus      float32
	Visited        bool
	NodeType       NodeType
	DistanceTarget float32
}

func (p *PathPoint) Push(writer base.Buffer) {
	writer.PushI32(p.X)
	writer.PushI32(p.Y)
	writer.PushI32(p.Z)
	writer.PushF32(p.DistanceOrigin)
	writer.PushF32(p.Cost)
	writer.PushF32(p.CostMalus)
	writer.PushBit(p.Visited)
	writer.PushI32(int32(p.NodeType))
	writer.PushF32(p.DistanceTarget)
}

func (p *PathPoint) Pull(reader base.Buffer) {
	p.X = reader.PullI32()
	p.Y = reader.PullI32()
	p.Z = reader.PullI32()
	p.DistanceOrigin = reader.PullF32()
	p.Cost = reader.PullF32()
	p.CostMalus = reader.PullF32()
	p.Visited = reader.PullBit()
	p.NodeType = NodeType(reader.PullI32())
	p.DistanceTarget = reader.PullF32()
}

type NodeType int

const (
	BLOCKED NodeType = iota
	OPEN
	WALKABLE
	TRAPDOOR
	FENCE
	LAVA
	WATER
	RAIL
	DANGER_FIRE
	DAMAGE_FIRE
	DANGER_CACTUS
	DAMAGE_CACTUS
	DANGER_OTHER
	DAMAGE_OTHER
	DOOR_OPEN
	DOOR_WOOD_CLOSED
	DOOR_IRON_CLOSED
)

func (d *DebugPaths) Chan() string {
	return "minecraft:debug/paths"
}

func (d *DebugPaths) Push(writer base.Buffer) {
	writer.PushI32(d.UnknownValue1)
	writer.PushF32(d.UnknownValue2)
	d.PathEntity.Push(writer)
}

func (d *DebugPaths) Pull(reader base.Buffer) {
	d.UnknownValue1 = reader.PullI32()
	d.UnknownValue2 = reader.PullF32()

	entity := PathEntity{}
	entity.Pull(reader)

	d.PathEntity = entity
}

type DebugNeighbors struct {
	Time     int64
	Location data.PositionI
}

func (d *DebugNeighbors) Chan() string {
	return "minecraft:debug/neighbors_update"
}

func (d *DebugNeighbors) Push(writer base.Buffer) {
	writer.PushVrL(d.Time)
	writer.PushPos(d.Location)
}

func (d *DebugNeighbors) Pull(reader base.Buffer) {
	d.Time = reader.PullVrL()
	d.Location = reader.PullPos()
}
