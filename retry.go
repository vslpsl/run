package run

import "time"

type (
	RetryFunc     func() error
	TestErrorFunc func(iter Iterator, err error) (retry bool)
)

type Iterator interface {
	Iteration() int
	SetDelay(duration time.Duration)
}

type iterator struct {
	iteration int
	delay     *time.Duration
}

func (i *iterator) Iteration() int {
	return i.iteration
}

func (i *iterator) SetDelay(duration time.Duration) {
	i.delay = &duration
}

func Retry(fn RetryFunc, testErrorFn TestErrorFunc) error {
	i := &iterator{}

	for {
		i.iteration++
		i.delay = nil

		err := fn()
		if err == nil {
			return nil
		}

		retry := testErrorFn(i, err)
		if retry == false {
			return err
		}

		if i.delay != nil {
			<-time.After(*i.delay)
		}
	}
}
