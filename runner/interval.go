package runner

import (
	"context"
	"os"
	"time"

	"github.com/tedsuo/ifrit"
	"github.com/topherbullock/varmint"
)

func NewInterval(
	runFunc RunFunc,
	interval time.Duration,
) Interval {
	return Interval{
		runFunc:  runFunc,
		interval: interval,
	}
}

type RunFunc func(ctx context.Context)

type Interval struct {
	runFunc  RunFunc
	interval time.Duration
	timer    *time.Timer

	statusChan chan varmint.Status
}

func (i *Interval) Stop() {
	if i.timer != nil {
		i.timer.Stop()
	}
	i.status(varmint.Paused)
}

func (i *Interval) Reset() {
	i.Stop()
	i.timer.Reset(i.interval)
	i.status(varmint.Waiting)
}

func (i *Interval) Start(ctx context.Context) {
	timer := time.NewTimer(i.interval)
	i.timer = timer

	go func() {
		defer timer.Stop()
		for {
			select {
			case <-timer.C:
				//TODO: What if I don't want the run func to block subsequent runs
				i.run(ctx)
				timer.Reset(i.interval)
			case <-ctx.Done():
				timer.Stop()
				i.status(varmint.Cancelled)
				return
			}
		}
	}()
	return
}

// Provides an ifrit.Runner
func (i *Interval) Runner(ctx context.Context) ifrit.Runner {
	newCtx, cancel := context.WithCancel(ctx)

	return ifrit.RunFunc(func(signals <-chan os.Signal, ready chan<- struct{}) error {
		close(ready)

		go i.Start(newCtx)

		select {
		case <-signals:
			cancel()
		}

		return nil
	})
}

func (i *Interval) Track() <-chan varmint.Status {
	// don't allocate status channel varmint starts Tracking it
	if i.statusChan == nil {
		i.statusChan = make(chan varmint.Status, 0)
	}
	return i.statusChan
}

func (i *Interval) run(ctx context.Context) {
	i.status(varmint.Running)
	defer i.status(varmint.Waiting)
	i.runFunc(ctx)
}

func (i *Interval) status(update varmint.Status) {
	if i.statusChan != nil {
		i.statusChan <- update
	}
}
