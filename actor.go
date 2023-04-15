package run

type actor struct {
	execute        func() error
	failureHandler func(err error)
	interrupt      func(err error)
}

func execute(a actor) error {
	err := a.execute()
	if err != nil && a.failureHandler != nil {
		a.failureHandler(err)
	}
	return err
}

type execResult struct {
	actor actor
	err   error
}
