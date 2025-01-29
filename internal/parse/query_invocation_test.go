package parse

import (
	"fmt"
	"testing"

	"github.com/turbot/pipe-fittings/v2/utils"
	"github.com/turbot/powerpipe/internal/resources"
)

// NOTE: all query arg values must be JSON representations
type parseQueryInvocationTest struct {
	input    string
	expected parseQueryInvocationResult
}

type parseQueryInvocationResult struct {
	queryName string
	args      *resources.QueryArgs
}

var emptyArgs = resources.NewQueryArgs()
var testCasesParseQueryInvocation = map[string]parseQueryInvocationTest{
	"no brackets": {
		input:    `query.q1`,
		expected: parseQueryInvocationResult{"query.q1", emptyArgs},
	},
	"no params": {
		input:    `query.q1()`,
		expected: parseQueryInvocationResult{"query.q1", emptyArgs},
	},
	"invalid params 1": {
		input: `query.q1(foo)`,
		expected: parseQueryInvocationResult{
			queryName: `query.q1`,
			args:      &resources.QueryArgs{},
		},
	},
	"invalid params 4": {
		input: `query.q1("foo",  "bar"])`,
		expected: parseQueryInvocationResult{
			queryName: `query.q1`,

			args: &resources.QueryArgs{},
		},
	},

	"single positional param": {
		input: `query.q1("foo")`,
		expected: parseQueryInvocationResult{
			queryName: `query.q1`,
			args:      &resources.QueryArgs{ArgList: []*string{utils.ToStringPointer("foo")}},
		},
	},
	"single positional param extra spaces": {
		input: `query.q1("foo"   )   `,
		expected: parseQueryInvocationResult{
			queryName: `query.q1`,
			args:      &resources.QueryArgs{ArgList: []*string{utils.ToStringPointer("foo")}},
		},
	},
	"multiple positional params": {
		input: `query.q1("foo", "bar", "foo-bar")`,
		expected: parseQueryInvocationResult{
			queryName: `query.q1`,
			args:      &resources.QueryArgs{ArgList: []*string{utils.ToStringPointer("foo"), utils.ToStringPointer("bar"), utils.ToStringPointer("foo-bar")}},
		},
	},
	"multiple positional params extra spaces": {
		input: `query.q1("foo",   "bar",    "foo-bar"   )`,
		expected: parseQueryInvocationResult{
			queryName: `query.q1`,
			args:      &resources.QueryArgs{ArgList: []*string{utils.ToStringPointer("foo"), utils.ToStringPointer("bar"), utils.ToStringPointer("foo-bar")}},
		},
	},
	"single named param": {
		input: `query.q1(p1 => "foo")`,
		expected: parseQueryInvocationResult{
			queryName: `query.q1`,
			args:      &resources.QueryArgs{ArgMap: map[string]string{"p1": "foo"}},
		},
	},
	"single named param extra spaces": {
		input: `query.q1(  p1  =>  "foo"  ) `,
		expected: parseQueryInvocationResult{
			queryName: `query.q1`,
			args:      &resources.QueryArgs{ArgMap: map[string]string{"p1": "foo"}},
		},
	},
	"multiple named params": {
		input: `query.q1(p1 => "foo", p2 => "bar")`,
		expected: parseQueryInvocationResult{
			queryName: `query.q1`,
			args:      &resources.QueryArgs{ArgMap: map[string]string{"p1": "foo", "p2": "bar"}},
		},
	},
	"multiple named params extra spaces": {
		input: ` query.q1 ( p1 => "foo" ,  p2  => "bar"     ) `,
		expected: parseQueryInvocationResult{
			queryName: `query.q1`,
			args:      &resources.QueryArgs{ArgMap: map[string]string{"p1": "foo", "p2": "bar"}},
		},
	},
	"named param with dot in value": {
		input: `query.q1(p1 => "foo.bar")`,
		expected: parseQueryInvocationResult{
			queryName: `query.q1`,
			args:      &resources.QueryArgs{ArgMap: map[string]string{"p1": "foo.bar"}},
		},
	},
}

func TestParseQueryInvocation(t *testing.T) {
	for name, test := range testCasesParseQueryInvocation {
		queryName, args, _ := ParseQueryInvocation(test.input)
		if args == nil {
			args = emptyArgs
		}
		if queryName != test.expected.queryName || !test.expected.args.Equals(args) {
			//nolint:forbidigo // acceptable
			fmt.Printf("")
			t.Errorf("Test: '%s'' FAILED : expected:\nquery: %s params: %s\n\ngot:\nquery: %s params: %s",
				name,
				test.expected.queryName,
				test.expected.args,
				queryName, args)
		}
	}
}
