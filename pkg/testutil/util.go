package testutil

import (
	"time"

	"github.com/stretchr/testify/suite"
)

func ExpectWithin[T any](t *suite.Suite, c chan T, to time.Duration) (out T) {
	timer := time.NewTimer(to)
	select {
	case <-timer.C:
		t.Fail("took too long to receive request")
		return
	case hit := <-c:
		return hit
	}
}
