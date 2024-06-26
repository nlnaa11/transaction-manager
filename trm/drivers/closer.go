package drivers

import "sync"

type Closer interface {
	Closed() <-chan struct{}
	Close()
}

type closer struct {
	once    sync.Once
	closeCh chan struct{}
}

func NewCloser() *closer {
	return &closer{
		closeCh: make(chan struct{}),
		once:    sync.Once{},
	}
}

func (c *closer) Closed() <-chan struct{} {
	return c.closeCh
}

func (c *closer) Close() {
	c.once.Do(func() {
		close(c.closeCh)
	})
}
