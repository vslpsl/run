package run

type Loop struct {
	actors []actor
}

func (l *Loop) Add(options ...Option) {
	a := actor{}
	for _, option := range options {
		option(&a)
	}

	if a.interrupt == nil {
		panic("interrupt must be implemented")
	}

	l.actors = append(l.actors, a)
}

func (l *Loop) Run(stopc <-chan struct{}) {
	if len(l.actors) == 0 {
		<-stopc
		return
	}

	execResultc := make(chan execResult)
	for _, a := range l.actors {
		go func(a actor) {
			execResultc <- execResult{
				actor: a,
				err:   execute(a),
			}
		}(a)
	}

	count := len(l.actors)
	stopped := false

	for {
		if count == 0 {
			return
		}

		select {
		case result := <-execResultc:
			if stopped {
				count--
				break
			}
			go func(a actor) {
				execResultc <- execResult{
					actor: a,
					err:   execute(a),
				}
			}(result.actor)
		case <-stopc:
			stopc = nil
			stopped = true
			go func(actors []actor) {
				for _, a := range actors {
					a.interrupt(nil)
				}
			}(l.actors)
		}
	}
}
