package runner

import "context"

type Stoppable interface {
	Stop()
}

type Resetable interface {
	Reset()
}

type Startable interface {
	Start(parentContext context.Context)
}

type Addable interface {
	Add(run RunFunc)
}

type RunFunc func(ctx context.Context)
