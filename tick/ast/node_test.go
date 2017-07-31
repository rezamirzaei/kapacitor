package ast

import (
	"fmt"
	"reflect"
	"testing"

	client "github.com/influxdata/kapacitor/client/v1"
	"github.com/stretchr/testify/assert"
)

func TestNumberNode(t *testing.T) {
	assert := assert.New(t)

	type testCase struct {
		Text    string
		Pos     position
		IsInt   bool
		IsFloat bool
		Int64   int64
		Float64 float64
		Err     error
	}

	test := func(tc testCase) {
		n, err := newNumber(tc.Pos, tc.Text, nil)
		if tc.Err != nil {
			assert.Equal(tc.Err, err)
		} else {
			if !assert.NotNil(n) {
				t.FailNow()
			}
			assert.Equal(tc.Pos.pos, n.position.pos)
			assert.Equal(tc.IsInt, n.IsInt)
			assert.Equal(tc.IsFloat, n.IsFloat)
			assert.Equal(tc.Int64, n.Int64)
			assert.Equal(tc.Float64, n.Float64)
		}
	}

	cases := []testCase{
		testCase{
			Text:  "04",
			Pos:   position{pos: 6},
			IsInt: true,
			Int64: 4,
		},
		testCase{
			Text:  "42",
			Pos:   position{pos: 5},
			IsInt: true,
			Int64: 42,
		},
		testCase{
			Text:    "42.21",
			Pos:     position{pos: 4},
			IsFloat: true,
			Float64: 42.21,
		},
		testCase{
			Text:    "42.",
			Pos:     position{pos: 3},
			IsFloat: true,
			Float64: 42.0,
		},
		testCase{
			Text:    "0.42",
			Pos:     position{pos: 2},
			IsFloat: true,
			Float64: 0.42,
		},
		testCase{
			Text: "0.4.2",
			Err:  fmt.Errorf("illegal number syntax: %q", "0.4.2"),
		},
		testCase{
			Text: "0x04",
			Err:  fmt.Errorf("illegal number syntax: %q", "0x04"),
		},
	}

	for _, tc := range cases {
		test(tc)
	}
}

func TestNewBinaryNode(t *testing.T) {
	assert := assert.New(t)

	type testCase struct {
		Left     Node
		Right    Node
		Operator token
	}

	test := func(tc testCase) {
		n := newBinary(position{pos: tc.Operator.pos}, tc.Operator.typ, tc.Left, tc.Right, false, nil)
		if !assert.NotNil(n) {
			t.FailNow()
		}
		assert.Equal(tc.Operator.pos, n.position.pos)
		assert.Equal(tc.Left, n.Right)
		assert.Equal(tc.Right, n.Left)
		assert.Equal(tc.Operator.typ, n.Operator)
	}

	cases := []testCase{
		testCase{
			Left:  nil,
			Right: nil,
			Operator: token{
				pos: 0,
				typ: TokenEqual,
				val: "=",
			},
		},
	}
	for _, tc := range cases {
		test(tc)
	}
}

func TestProgramNodeDBRPs(t *testing.T) {
	type testCase struct {
		name       string
		tickscript string
		dbrps      []client.DBRP
	}

	tt := []testCase{
		{
			name: "one dbrp",
			tickscript: `dbrp "telegraf"."autogen"
			
			stream|from().measurement('m')
			`,
			dbrps: []client.DBRP{
				{
					Database:        "telegraf",
					RetentionPolicy: "auotgen",
				},
			},
		},
		{
			name: "two dbrp",
			tickscript: `dbrp "telegraf"."autogen"
			dbrp "telegraf"."not_autogen"
			
			stream|from().measurement('m')
			`,
			dbrps: []client.DBRP{
				{
					Database:        "telegraf",
					RetentionPolicy: "auotgen",
				},
				{
					Database:        "telegraf",
					RetentionPolicy: "not_autogen",
				},
			},
		},
	}

	for _, tst := range tt {
		t.Run(tst.name, func(t *testing.T) {
			pn, err := NewProgramNodeFromTickscript(tst.tickscript)
			if err != nil {
				t.Fatalf("error parsing tickscript: %v", err)
			}
			if exp, got := tst.dbrps, pn.DBRPs(); reflect.DeepEqual(exp, got) {
				t.Fatalf("DBRPs do not match:\nexp: %v,\ngot %v", exp, got)
			}
		})
	}
}

func TestProgramNodeTaskType(t *testing.T) {
	type testCase struct {
		name     string
		program  string
		taskType client.TaskType
	}
}
