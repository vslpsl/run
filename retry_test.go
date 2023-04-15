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
	err := run.Retry(
		func() error {
			count++
			return testError
		},
		func(iter run.Iterator, err error) (retry bool) {
			if iter.Iteration() == 3 {
				return false
			}
			iter.SetDelay(time.Second * 3)
			return true
		})

	require.ErrorIs(t, err, testError)
	require.Equal(t, count, 3)
}
