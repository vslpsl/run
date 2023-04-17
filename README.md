# run
inspired by github.com/oklog/run

This package provides components for concurrency management.There are three components: Group, Loop and Retry func.

## Group
Provides concurrent execution of goroutines until first error occurred.
For each actor optional failureHandler can be provided.

```go
group := run.Group{}
	{
		server := http.Server{}
		group.Add(
			run.Execute(func() error {
				return server.ListenAndServe()
			}),
			run.Interrupt(func(err error) {
				ctx, cancel := context.WithCancel(context.Background())
				defer cancel()
				_ = server.Shutdown(ctx)
			}),
		)
	}
	{
		group.Add(run.Signal(run.DefaultSignals()...))
	}
	{
		ctx, cancel := context.WithCancel(context.Background())
		group.Add(
			run.Execute(func() error {
				return myCode(ctx, ...)
			}),
			run.FailureHandler(func(err error) {
				handleErr(err)
            }),
			run.Interrupt(func(err error) {
				cancel()
			}),
		)
	}
	
	err := group.Run()
```

## Loop
Provides concurrent repeatable execution of goroutines until loop is stopped.
For each actor optional failureHandler can be provided.

```go
ctx, cancel := context.WithCancel(context.Background())
loop := run.Loop{}
{
    ctx, cancel := context.WithCancel(context.Background())
    loop.Add(
        run.Execute(func() error {
            return repeatableTask(ctx, ...)
        }),
        run.Interrupt(func(err error) {
            cancel()
        }),
    )
}
{
    ctx, cancel := context.WithCancel(context.Background())
    loop.Add(
        run.Execute(func() error {
            return anotherRepeatableTask(ctx, ...)
        }),
        run.Interrupt(func(err error) {
            cancel()
        }),
    )
}
{
    ctx, cancel := context.WithCancel(context.Background())
    loop.Add(
        run.Execute(func() error {
			return yetAnotherRepeatableTask(ctx, ...)
		}),
		run.FailureHandler(func(err error) {
            sentry.CaptureException(err)
        }),
        run.Interrupt(func(err error) {
            cancel()
        }),
    )
}

go func() {
    <-time.After(time.Minute)
    cancel()
}()

loop.Run(ctx.Done())
```

## Retry
Retries provided function with user control via Repeater interface.

```go
ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
defer cancel()
err = run.Retry(func(retrier run.Retrier) error {
    log.Println(retrier.Iteration())
    err := doSomething(ctx, ...)
    if errors.Is(err, ctx.Err()) {
        retrier.Stop()
    } else {
        retrier.SetDelay(time.Second * 3)
    }
    return err
})
```