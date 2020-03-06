package event

type Cancellable struct {
	cancelled bool
}

func (c *Cancellable) GetCancelled() bool {
	return c.cancelled
}

func (c *Cancellable) SetCancelled(cancelled bool) {
	c.cancelled = cancelled
}
