package signals

import (
	"fmt"
	"os"
	"os/signal"
)

var (
	c = make(chan os.Signal, 1)
	shouldGracefulExit = false
)

func ShouldContinue () bool {
	return !shouldGracefulExit;
}

func InstallInterruptSignal() {
	signal.Notify(c, os.Interrupt)

	go func() {
		<- c
		shouldGracefulExit = true
		fmt.Println("Graceful Exit Requested: Ctrl + C")
	}()
}
