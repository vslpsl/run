package run

type Option func(a *actor)

func Execute(fn func() error) Option {
	return func(a *actor) {
		a.execute = fn
	}
}

func FailureHandler(handler func(err error)) Option {
	return func(a *actor) {
		a.failureHandler = handler
	}
}

func Interrupt(fn func(err error)) Option {
	return func(a *actor) {
		a.interrupt = fn
	}
}

func NoOpInterrupt() Option {
	return func(a *actor) {
		a.interrupt = func(err error) {}
	}
}
