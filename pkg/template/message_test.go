package template

import (
	"testing"

	"github.com/skiff-sh/config/ptr"
	test2 "github.com/skiff-sh/pilot/api/go/test"
	"github.com/stretchr/testify/suite"
)

type MessageTestSuite struct {
	suite.Suite
}

func (m *MessageTestSuite) TestApply() {
	type test struct {
		Given            *test2.Primitives
		GivenData        Data
		Expected         *test2.Primitives
		ExpectedErr      string
		ExpectedApplyErr string
	}

	tests := map[string]test{
		"happy path": {
			Given: &test2.Primitives{
				Str: "{{ str }}",
			},
			GivenData: Data{
				"str": "derp",
			},
			Expected: &test2.Primitives{
				Str: "derp",
			},
		},
		"skips primitive values": {
			Given:     &test2.Primitives{StrPtr: ptr.Ptr("hi")},
			GivenData: Data{"hello": "derp"},
			Expected:  &test2.Primitives{StrPtr: ptr.Ptr("hi")},
		},
		"err expression": {
			Given:       &test2.Primitives{Str: "{{ hello }"},
			GivenData:   Data{"hello": "derp"},
			ExpectedErr: `field "str" has invalid expression "{{ hello }": [Error (where: parser) in <string> | Line 1 Col 10 near '}'] '}}' expected`,
		},
	}

	for desc, v := range tests {
		m.Run(desc, func() {
			exprs, err := NewFieldTemplates(v.Given)
			if v.ExpectedErr != "" || !m.NoError(err) {
				m.EqualError(err, v.ExpectedErr)
				return
			}

			err = exprs.Apply(v.Given, v.GivenData)
			if v.ExpectedApplyErr != "" || !m.NoError(err) {
				m.EqualError(err, v.ExpectedApplyErr)
				return
			}

			m.Equal(v.Expected.String(), v.Given.String())
		})
	}
}

func TestMessageTestSuite(t *testing.T) {
	suite.Run(t, new(MessageTestSuite))
}
