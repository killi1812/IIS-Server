package scheduler

import (
	"context"
	"iis_server/httpserver"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"go.uber.org/zap"
)

var signalNotificationCh = make(chan os.Signal, 1)

func Start() {
	// relay selected signals to channel
	// - os.Interrupt, ctrl-c
	// - syscall.SIGTERM, program termination
	signal.Notify(signalNotificationCh, os.Interrupt, syscall.SIGTERM)

	// create scheduler
	schedulerWg := sync.WaitGroup{}
	schedulerCtx := context.Background()
	schedulerCtx, schedulerCancel := context.WithCancel(schedulerCtx)
	zap.S().Debugf("Created scheduler context")

	schedulerWg.Add(1)
	go CheckInterrupt(schedulerCtx, &schedulerWg, schedulerCancel)
	zap.S().Debugf("Started CheckInterrupt")

	schedulerWg.Add(1)
	go httpserver.Start(schedulerCtx, &schedulerWg, schedulerCancel)
	zap.S().Debugf("Started HTTP server")

	schedulerWg.Wait()

	zap.S().Debugf("Terminated program")
}

func CheckInterrupt(ctx context.Context, wg *sync.WaitGroup, schedulerCancel context.CancelFunc) {
	defer wg.Done()

	for {
		select {

		case <-ctx.Done():
			zap.S().Debugf("Terminated CheckInterrupt")
			return

		case sig := <-signalNotificationCh:
			zap.S().Debugf("Received signal on notification channel, signal = %v", sig)
			schedulerCancel()
		}
	}
}
