package run

type Group struct {
	actors []actor
}

func (g *Group) Add(options ...Option) {
	a := actor{}
	for _, option := range options {
		option(&a)
	}

	if a.interrupt == nil {
		panic("interrupt must be declared")
	}

	g.actors = append(g.actors, a)
}

func (g *Group) Run() error {
	if len(g.actors) == 0 {
		return nil
	}

	var err error
	execResultc := make(chan execResult, len(g.actors))
	for _, a := range g.actors {
		go func(a actor) {
			execResultc <- execResult{
				actor: a,
				err:   execute(a),
			}
		}(a)
	}

	result := <-execResultc
	err = result.err
	for _, a := range g.actors {
		a.interrupt(result.err)
	}

	for i := 0; i < len(g.actors)-1; i++ {
		<-execResultc
	}

	return err
}
