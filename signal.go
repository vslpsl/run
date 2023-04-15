package run

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func DefaultSignals() []os.Signal {
	return []os.Signal{os.Interrupt, syscall.SIGTERM}
}

func Signal(signals ...os.Signal) Option {
	ctx, cancel := context.WithCancel(context.Background())
	signalc := make(chan os.Signal, 1)
	signal.Notify(signalc, signals...)
	return func(a *actor) {
		a.execute = func() error {
			select {
			case sig := <-signalc:
				return SignalError{Signal: sig}
			case <-ctx.Done():
				return ctx.Err()
			}
		}
		a.interrupt = func(err error) {
			cancel()
		}
	}
}

type SignalError struct {
	Signal os.Signal
}

func (e SignalError) Error() string {
	return fmt.Sprintf("signal: %s", e.Signal)
}

func IsSignal(err error) bool {
	return errors.Is(err, SignalError{})
}
