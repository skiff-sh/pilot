package testutil

import (
	"time"

	"github.com/google/go-cmp/cmp"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/testing/protocmp"

	"github.com/stretchr/testify/suite"
)

func ExpectWithin[T any](t *suite.Suite, c chan T, to time.Duration) (out T) {
	timer := time.NewTimer(to)
	select {
	case <-timer.C:
		return
	case hit := <-c:
		return hit
	}
}

// DiffProto returns the diff of two messages and an empty string if there
// is no difference.
func DiffProto(expected, actual proto.Message) string {
	return cmp.Diff(expected, actual, protocmp.Transform())
}
