package runner_test

import (
	"context"
	//"fmt"
	"time"

	"github.com/topherbullock/varmint"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/topherbullock/varmint/runner"
)

var _ = Describe("Interval", func() {

	var (
		intervalRunner runner.Interval
		runFunc        runner.RunFunc
		interval       = 1 * time.Second
		messages       []int64
		ctx            context.Context
		cancel         context.CancelFunc
	)

	BeforeEach(func() {
		ctx, cancel = context.WithCancel(context.Background())

		runFunc = func(ctx context.Context) {
			now := time.Now()
			messages = append(messages, now.UnixNano())
		}

		intervalRunner = runner.NewInterval(runFunc, interval)
	})

	AfterEach(func() {
		cancel()
	})

	Describe("Start", func() {
		JustBeforeEach(func() {
			intervalRunner.Start(ctx)
		})

		It("runs on the desired interval", func() {
			time.AfterFunc(3050*time.Millisecond, cancel)
			<-ctx.Done()
			Expect(messages).To(HaveLen(3))
		})
	})

	Describe("Track", func() {
		var (
			statusChan <-chan varmint.Status
		)

		JustBeforeEach(func() {
			statusChan = intervalRunner.Track()
			intervalRunner.Start(ctx)
		})

		It("Emits statuses when running until cancelled", func() {
			time.AfterFunc(2010*time.Millisecond, cancel)
			event := <-statusChan
			Expect(event.String()).To(Equal("Running"))
			event = <-statusChan
			Expect(event.String()).To(Equal("Waiting"))
			event = <-statusChan
			Expect(event.String()).To(Equal("Running"))
			event = <-statusChan
			Expect(event.String()).To(Equal("Waiting"))
			<-ctx.Done()
			event = <-statusChan
			Expect(event.String()).To(Equal("Cancelled"))
		})

		It("Emits status when stopped", func() {
			time.AfterFunc(2010*time.Millisecond, cancel)
			time.AfterFunc(1010*time.Millisecond, intervalRunner.Stop)
			event := <-statusChan
			Expect(event.String()).To(Equal("Running"))

			event = <-statusChan
			Expect(event.String()).To(Equal("Waiting"))

			event = <-statusChan
			Expect(event.String()).To(Equal("Paused"))

			<-ctx.Done()
			event = <-statusChan
			Expect(event.String()).To(Equal("Cancelled"))
		})
	})
})
