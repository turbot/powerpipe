//go:build ignore

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: go run generate_test_mods.go <output_dir> <size>")
		fmt.Println("Sizes: small (10 dashboards), medium (50), large (200), xlarge (500)")
		os.Exit(1)
	}

	outputDir := os.Args[1]
	size := os.Args[2]

	counts := map[string]struct {
		dashboards int
		queries    int
		controls   int
		benchmarks int
	}{
		"small":  {10, 20, 30, 5},
		"medium": {50, 100, 150, 20},
		"large":  {200, 400, 500, 50},
		"xlarge": {500, 1000, 1500, 100},
	}

	c, ok := counts[size]
	if !ok {
		fmt.Printf("Unknown size: %s\n", size)
		os.Exit(1)
	}

	generateMod(outputDir, size, c.dashboards, c.queries, c.controls, c.benchmarks)
}

func generateMod(dir, name string, dashboards, queries, controls, benchmarks int) {
	if err := os.MkdirAll(dir, 0755); err != nil {
		fmt.Printf("Failed to create directory: %v\n", err)
		os.Exit(1)
	}

	// Generate mod.pp
	modContent := fmt.Sprintf(`mod "%s_test" {
  title = "%s Test Mod"
  description = "Generated mod for performance testing with %d dashboards, %d queries, %d controls, %d benchmarks"
}
`, name, strings.ToUpper(name[:1])+name[1:], dashboards, queries, controls, benchmarks)

	if err := os.WriteFile(filepath.Join(dir, "mod.pp"), []byte(modContent), 0644); err != nil {
		fmt.Printf("Failed to write mod.pp: %v\n", err)
		os.Exit(1)
	}

	// Generate queries
	var queryContent strings.Builder
	for i := 0; i < queries; i++ {
		queryContent.WriteString(fmt.Sprintf(`
query "query_%d" {
  title = "Query %d"
  description = "Test query %d for performance benchmarking"
  sql = <<-EOQ
    SELECT
      %d as id,
      'value_%d' as name,
      now() as created_at
    FROM generate_series(1, 100)
  EOQ
}
`, i, i, i, i, i))
	}
	if err := os.WriteFile(filepath.Join(dir, "queries.pp"), []byte(queryContent.String()), 0644); err != nil {
		fmt.Printf("Failed to write queries.pp: %v\n", err)
		os.Exit(1)
	}

	// Generate controls
	var controlContent strings.Builder
	for i := 0; i < controls; i++ {
		controlContent.WriteString(fmt.Sprintf(`
control "control_%d" {
  title = "Control %d"
  description = "Test control %d"
  sql = <<-EOQ
    SELECT
      'resource_%d' as resource,
      'ok' as status,
      'Control %d passed' as reason
  EOQ
  tags = {
    category = "test"
    index    = "%d"
  }
}
`, i, i, i, i, i, i))
	}
	if err := os.WriteFile(filepath.Join(dir, "controls.pp"), []byte(controlContent.String()), 0644); err != nil {
		fmt.Printf("Failed to write controls.pp: %v\n", err)
		os.Exit(1)
	}

	// Generate benchmarks
	var benchmarkContent strings.Builder
	controlsPerBenchmark := controls / benchmarks
	if controlsPerBenchmark == 0 {
		controlsPerBenchmark = 1
	}
	for i := 0; i < benchmarks; i++ {
		children := make([]string, 0, controlsPerBenchmark)
		for j := 0; j < controlsPerBenchmark && (i*controlsPerBenchmark+j) < controls; j++ {
			children = append(children, fmt.Sprintf("control.control_%d", i*controlsPerBenchmark+j))
		}
		benchmarkContent.WriteString(fmt.Sprintf(`
benchmark "benchmark_%d" {
  title = "Benchmark %d"
  description = "Test benchmark %d"
  children = [
    %s
  ]
  tags = {
    category = "test"
  }
}
`, i, i, i, strings.Join(children, ",\n    ")))
	}
	if err := os.WriteFile(filepath.Join(dir, "benchmarks.pp"), []byte(benchmarkContent.String()), 0644); err != nil {
		fmt.Printf("Failed to write benchmarks.pp: %v\n", err)
		os.Exit(1)
	}

	// Generate dashboards
	var dashboardContent strings.Builder
	for i := 0; i < dashboards; i++ {
		dashboardContent.WriteString(fmt.Sprintf(`
dashboard "dashboard_%d" {
  title = "Dashboard %d"
  description = "Test dashboard %d for performance benchmarking"

  tags = {
    category = "test"
    index    = "%d"
  }

  container {
    title = "Overview"

    card {
      width = 2
      sql   = query.query_%d.sql
    }

    card {
      width = 2
      sql   = "SELECT count(*) as total FROM generate_series(1, 100)"
    }

    chart {
      width = 6
      type  = "bar"
      sql   = query.query_%d.sql
    }
  }

  container {
    title = "Details"

    table {
      sql = query.query_%d.sql
    }
  }
}
`, i, i, i, i, i%queries, i%queries, i%queries))
	}
	if err := os.WriteFile(filepath.Join(dir, "dashboards.pp"), []byte(dashboardContent.String()), 0644); err != nil {
		fmt.Printf("Failed to write dashboards.pp: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Generated %s mod in %s:\n", name, dir)
	fmt.Printf("  - %d dashboards\n", dashboards)
	fmt.Printf("  - %d queries\n", queries)
	fmt.Printf("  - %d controls\n", controls)
	fmt.Printf("  - %d benchmarks\n", benchmarks)
}
