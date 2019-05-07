package runner_test

import (
	"context"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/topherbullock/varmint/runner"
)

var _ = Describe("Parallel", func() {

	var (
		parallelRunner runner.Parallel
		runFunc        runner.RunFunc
		ctx            context.Context
		cancel         context.CancelFunc
		maxInFlight    = 3
		messages       chan string
	)

	Context("With an initial set of members", func() {
		BeforeEach(func() {
			messages = make(chan string, 0)
			ctx, cancel = context.WithCancel(context.Background())
			runFunc = func(ctx context.Context) {
				messages <- "foo"
			}
			parallelRunner = runner.NewParallel(maxInFlight, runFunc)
		})

		Describe("Start", func() {
			BeforeEach(func() {
				parallelRunner.Start(ctx)
			})

			It("runs each of its members", func() {
				Eventually(messages).Should(Receive("foo"))
				Expect(true).To(BeFalse())
			})
		})
		AfterEach(func() {
			cancel()
		})
	})

	// BeforeEach(func() {
	// 	ctx, cancel = context.WithCancel(context.Background())
	// 	runFunc = func(ctx context.Context) {}
	// 	parallelRunner = runner.NewParallel(maxInFlight)
	// })
	//
	// AfterEach(func() {
	// 	cancel()
	// })
})
