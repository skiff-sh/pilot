package testutil

import (
	"testing"
	"time"

	"github.com/skiff-sh/config/ptr"
	"github.com/skiff-sh/pilot/api/go/pilot"
	"github.com/stretchr/testify/suite"
	"google.golang.org/protobuf/proto"
)

type UtilTestSuite struct {
	suite.Suite
}

func (u *UtilTestSuite) TestDiffProto() {
	type test struct {
		Expected proto.Message
		Actual   proto.Message
		Same     bool
	}

	tests := map[string]test{
		"matches": {
			Expected: &pilot.ProvokeBehavior_Request{Name: "name"},
			Actual:   &pilot.ProvokeBehavior_Request{Name: "name"},
			Same:     true,
		},
		"no match": {
			Expected: &pilot.ProvokeBehavior_Request{Name: "name"},
			Actual:   &pilot.ProvokeBehavior_Request{Name: "name1"},
		},
	}

	for desc, v := range tests {
		u.Run(desc, func() {
			u.Equal(v.Same, len(DiffProto(v.Expected, v.Actual)) == 0)
		})
	}
}

func (u *UtilTestSuite) TestExpectWithin() {
	type test struct {
		TO       time.Duration
		Given    *int
		Expected int
	}

	tests := map[string]test{
		"value": {
			TO:       time.Second,
			Given:    ptr.Ptr(1),
			Expected: 1,
		},
		"timeout": {
			TO: time.Millisecond,
		},
	}

	for desc, v := range tests {
		u.Run(desc, func() {
			c := make(chan int, 1)

			if v.Given != nil {
				c <- *v.Given
			}

			u.Equal(v.Expected, ExpectWithin(&u.Suite, c, v.TO))
		})
	}
}

func TestUtilTestSuite(t *testing.T) {
	suite.Run(t, new(UtilTestSuite))
}
