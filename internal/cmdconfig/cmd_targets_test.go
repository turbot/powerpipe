package cmdconfig

import (
	"slices"
	"testing"

	"github.com/turbot/pipe-fittings/v2/parse"
	"github.com/turbot/powerpipe/internal/resources"
)

func TestResourceTargets(t *testing.T) {
	cases := map[string]struct {
		args []string
		want []parse.ResourceTarget
	}{
		"mod qualified":   {[]string{"aws_compliance.benchmark.cis_v300"}, []parse.ResourceTarget{{Mod: "aws_compliance", Name: "benchmark.cis_v300"}}},
		"type and name":   {[]string{"benchmark.cis_v300"}, []parse.ResourceTarget{{Name: "benchmark.cis_v300"}}},
		"name only":       {[]string{"cis_v300"}, []parse.ResourceTarget{{Name: "benchmark.cis_v300"}}},
		"several":         {[]string{"a.benchmark.b1", "b2"}, []parse.ResourceTarget{{Mod: "a", Name: "benchmark.b1"}, {Name: "benchmark.b2"}}},
		"all":             {[]string{"all"}, nil},
		"wrong type":      {[]string{"aws_compliance.control.c1"}, nil},
		"query arguments": {[]string{"benchmark.b1(1)"}, nil},
		"too many parts":  {[]string{"a.b.c.d"}, nil},
		"one bad arg":     {[]string{"benchmark.b1", "dashboard.d"}, nil},
		"no args":         {nil, nil},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			if got := ResourceTargets[*resources.Benchmark](tc.args); !slices.Equal(got, tc.want) {
				t.Fatalf("got %v, want %v", got, tc.want)
			}
		})
	}
}

func TestResourceTargets_Commands(t *testing.T) {
	if got := ResourceTargets[*resources.Control]([]string{"c1"}); !slices.Equal(got, []parse.ResourceTarget{{Name: "control.c1"}}) {
		t.Fatalf("expected a control target, got %v", got)
	}
	// other commands load every resource
	if got := ResourceTargets[*resources.Dashboard]([]string{"dashboard.d"}); got != nil {
		t.Fatalf("expected no targets for a dashboard, got %v", got)
	}
	if got := ResourceTargets[*resources.DetectionBenchmark]([]string{"benchmark.b"}); got != nil {
		t.Fatalf("expected no targets for a detection benchmark, got %v", got)
	}
}
