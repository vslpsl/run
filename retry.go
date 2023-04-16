package run

import "time"

type Retrier interface {
	Iteration() int
	SetDelay(duration time.Duration)
	Stop()
}

type retrier struct {
	iteration int
	delay     *time.Duration
	stopped   bool
}

// Iteration returns number of passed iterations.
func (r *retrier) Iteration() int {
	return r.iteration
}

// SetDelay sets delay for next iteration.
func (r *retrier) SetDelay(duration time.Duration) {
	r.delay = &duration
}

// Stop preventing further executions.
func (r *retrier) Stop() {
	r.stopped = true
}

func Retry(fn func(retrier Retrier) error) error {
	r := &retrier{}

	for {
		r.iteration++
		r.delay = nil

		err := fn(r)
		if err == nil {
			return nil
		}

		if r.stopped {
			return err
		}

		if r.delay != nil {
			<-time.After(*r.delay)
		}
	}
}
