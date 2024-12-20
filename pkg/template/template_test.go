package template

import (
	"testing"

	test2 "github.com/skiff-sh/pilot/api/go/test"
	"github.com/stretchr/testify/suite"
	"google.golang.org/protobuf/types/known/structpb"
)

type TemplateTestSuite struct {
	suite.Suite
}

func (t *TemplateTestSuite) TestSet() {
	type test struct {
		Data     Data
		Key      string
		Val      any
		Expected Data
	}

	tests := map[string]test{
		"no dot": {
			Data:     Data{},
			Key:      "derp",
			Val:      1,
			Expected: Data{"derp": 1},
		},
		"overwrite": {
			Data:     Data{"derp": 1},
			Key:      "derp",
			Val:      2,
			Expected: Data{"derp": 2},
		},
		"multiple dots": {
			Data:     Data{},
			Key:      "derp.flerp.blerp",
			Val:      1,
			Expected: Data{"derp": Data{"flerp": Data{"blerp": 1}}},
		},
	}

	for desc, v := range tests {
		t.Run(desc, func() {
			v.Data.Set(v.Key, v.Val)
			t.Equal(v.Expected, v.Data)
		})
	}
}

func (t *TemplateTestSuite) TestUnmarshal() {
	type test struct {
		Given       any
		Expected    Data
		ExpectedErr string
	}

	tests := map[string]test{
		"protobuf": {
			Given:    &test2.Primitives{Str: "derp"},
			Expected: Data{"str": "derp"},
		},
		"string map": {
			Given:    map[string]string{"derp": "flerp"},
			Expected: Data{"derp": "flerp"},
		},
		"any map": {
			Given:    map[string]any{"derp": 1},
			Expected: Data{"derp": 1},
		},
		"data": {
			Given:    Data{"derp": 2},
			Expected: Data{"derp": 2},
		},
	}

	for desc, v := range tests {
		t.Run(desc, func() {
			actual, err := Unmarshal(v.Given)
			if v.ExpectedErr != "" || !t.NoError(err) {
				t.EqualError(err, v.ExpectedErr)
				return
			}

			t.Equal(v.Expected, actual)
		})
	}
}

func (t *TemplateTestSuite) TestIsTruthy() {
	type test struct {
		Given    any
		Expected bool
	}

	tests := map[string]test{
		"str": {
			Given:    "asd",
			Expected: true,
		},
		"str false": {Given: "falsE "},
		"str true": {
			Given:    " TrUe",
			Expected: true,
		},
		"str one": {
			Given:    " 1",
			Expected: true,
		},
		"str zero": {Given: "0"},
		"str nil":  {Given: "nil "},
		"str null": {Given: "nulL "},
		"bool": {
			Given:    true,
			Expected: true,
		},
		"bool false": {Given: false},
		"nil val":    {},
		"zero":       {Given: 0},
	}

	for desc, v := range tests {
		t.Run(desc, func() {
			t.Equal(v.Expected, IsTruthy(v.Given))
		})
	}
}

func (t *TemplateTestSuite) TestToProto() {
	type test struct {
		Given    Data
		Expected *structpb.Struct
	}

	tests := map[string]test{
		"happy": {
			Given: Data{"derp": 1, "hi": Data{"qwe": true}},
			Expected: &structpb.Struct{
				Fields: map[string]*structpb.Value{
					"derp": structpb.NewNumberValue(1),
					"hi": structpb.NewStructValue(&structpb.Struct{Fields: map[string]*structpb.Value{
						"qwe": structpb.NewBoolValue(true),
					}}),
				},
			},
		},
	}

	for desc, v := range tests {
		t.Run(desc, func() {
			t.Equal(v.Expected.String(), v.Given.ToProto().String())
		})
	}
}

func (t *TemplateTestSuite) TestGet() {
	type test struct {
		Given    Data
		Key      string
		Expected any
	}

	tests := map[string]test{
		"exists": {
			Key:      "derp.hi",
			Given:    Data{"derp": Data{"hi": "there"}},
			Expected: "there",
		},
		"empty key": {
			Given: Data{"hi": "there"},
		},
		"key not found": {
			Key:   "derp",
			Given: Data{"derp1": 1},
		},
		"double dot": {
			Given: Data{"derp": 1},
			Key:   "..",
		},
	}

	for desc, v := range tests {
		t.Run(desc, func() {
			actual, _ := v.Given.Get(v.Key)
			t.Equal(v.Expected, actual)
		})
	}
}

func TestTemplateTestSuite(t *testing.T) {
	suite.Run(t, new(TemplateTestSuite))
}
