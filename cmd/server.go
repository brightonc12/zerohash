package main

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
	vwapHdl "zerohash/internal/handlers"
)

func main() {
	ctx, cancelCtx := context.WithCancel(context.Background())
	wg := &sync.WaitGroup{}

	vwapHdl.RunVWapAgainstTrade(ctx, wg)

	termChan := make(chan os.Signal)

	signal.Notify(termChan, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<- termChan
	cancelCtx()
	wg.Wait()
}
