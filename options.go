package run

type Option func(a *actor)

func Execute(fn func() error) Option {
	return func(a *actor) {
		a.execute = fn
	}
}

func ExecuteWithFailureHandler(executeFn func() error, failureHandler func(err error)) Option {
	return func(a *actor) {
		a.execute = executeFn
		a.failureHandler = failureHandler
	}
}

func Interrupt(fn func(err error)) Option {
	return func(a *actor) {
		a.interrupt = fn
	}
}
