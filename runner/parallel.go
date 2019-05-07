package runner

import (
	"context"
	"errors"
	"sync/atomic"
)

type Parallel interface {
	Startable
	Stoppable
}

func NewParallel(maxInFlight int, members ...RunFunc) Parallel {
	return &parallel{
		maxInFlight: maxInFlight,
		members:     members,
	}
}

type parallel struct {
	maxInFlight int
	rootContext context.Context
	rootCancel  context.CancelFunc

	members  []RunFunc
	children int64
}

func (p *parallel) Start(parentContext context.Context) {
	p.rootContext, p.rootCancel = context.WithCancel(parentContext)
}

func (p *parallel) Add(run RunFunc) error {
	if int(p.children) < p.maxInFlight {
		go p.spawnChild(run)
		return nil
	}
	return errors.New("max in flight reached")
}

func (p *parallel) spawnChild(run RunFunc) {
	atomic.AddInt64(&p.children, 1)
	run(p.rootContext)
	atomic.AddInt64(&p.children, -1)
}

func (p *parallel) Stop() {
	p.rootCancel()
}
