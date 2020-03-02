package base

type Loads interface {
	Load()
}

type Kills interface {
	Kill()
}

type State interface {
	Loads
	Kills
}
