package run_test

import (
	"errors"
	"github.com/stretchr/testify/require"
	"github.com/vslpsl/run"
	"testing"
	"time"
)

func TestRetry(t *testing.T) {
	testError := errors.New("test-error")
	count := 0
	err := run.Retry(func(retrier run.Retrier) error {
		count++
		if retrier.Iteration() == 3 {
			retrier.Stop()
			return testError
		}
		retrier.SetDelay(time.Second)
		return testError
	})

	require.ErrorIs(t, err, testError)
	require.Equal(t, count, 3)
}
