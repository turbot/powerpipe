//go:build ignore

// generate_lazy_test_mods.go generates test mods specifically designed for
// lazy loading testing. It creates mods with various reference patterns,
// hierarchies, and edge cases.
//
// Usage: go run generate_lazy_test_mods.go <output_dir> <preset>
// Presets: small, medium, large, stress

package main

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
)

type Config struct {
	Dashboards       int
	Queries          int
	Controls         int
	Benchmarks       int
	MaxBenchDepth    int
	MaxContainerNest int
	CrossRefPercent  int // Percentage of resources that reference others
}

var presets = map[string]Config{
	"small": {
		Dashboards:       10,
		Queries:          20,
		Controls:         30,
		Benchmarks:       5,
		MaxBenchDepth:    3,
		MaxContainerNest: 2,
		CrossRefPercent:  30,
	},
	"medium": {
		Dashboards:       50,
		Queries:          100,
		Controls:         150,
		Benchmarks:       25,
		MaxBenchDepth:    5,
		MaxContainerNest: 3,
		CrossRefPercent:  40,
	},
	"large": {
		Dashboards:       200,
		Queries:          400,
		Controls:         600,
		Benchmarks:       75,
		MaxBenchDepth:    7,
		MaxContainerNest: 4,
		CrossRefPercent:  50,
	},
	"stress": {
		Dashboards:       500,
		Queries:          1000,
		Controls:         2000,
		Benchmarks:       200,
		MaxBenchDepth:    10,
		MaxContainerNest: 5,
		CrossRefPercent:  60,
	},
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: go run generate_lazy_test_mods.go <output_dir> <preset>")
		fmt.Println("Presets: small, medium, large, stress")
		fmt.Println()
		fmt.Println("This generator creates test mods specifically designed for lazy loading testing,")
		fmt.Println("including varied reference patterns, deep hierarchies, and cross-references.")
		os.Exit(1)
	}

	outputDir := os.Args[1]
	preset := os.Args[2]

	config, ok := presets[preset]
	if !ok {
		fmt.Printf("Unknown preset: %s\n", preset)
		fmt.Println("Available presets: small, medium, large, stress")
		os.Exit(1)
	}

	// Use deterministic seed for reproducible output
	rand.Seed(42)

	generateLazyMod(outputDir, preset, config)
}

func generateLazyMod(dir, name string, c Config) {
	if err := os.MkdirAll(dir, 0755); err != nil {
		fmt.Printf("Failed to create directory: %v\n", err)
		os.Exit(1)
	}

	// Generate mod.pp with locals
	modContent := fmt.Sprintf(`mod "lazy_%s" {
  title       = "Lazy Loading %s Test"
  description = "Generated mod for lazy loading testing: %d dashboards, %d queries, %d controls, %d benchmarks"
}

# Common tags for all resources
locals {
  common_tags = {
    generator = "lazy_test"
    preset    = "%s"
    test      = "true"
  }

  severity_levels = ["low", "medium", "high", "critical"]
}

# Variables for parameterized resources
variable "test_region" {
  type        = string
  default     = "us-east-1"
  description = "Test region for lazy loading tests"
}

variable "test_threshold" {
  type    = number
  default = 100
}
`, name, strings.Title(name), c.Dashboards, c.Queries, c.Controls, c.Benchmarks, name)

	writeFile(filepath.Join(dir, "mod.pp"), modContent)

	// Generate queries with various patterns
	generateQueries(dir, c)

	// Generate controls with cross-references
	generateControls(dir, c)

	// Generate benchmarks with hierarchies
	generateBenchmarks(dir, c)

	// Generate dashboards with nested containers
	generateDashboards(dir, c)

	fmt.Printf("Generated lazy loading test mod '%s' in %s:\n", name, dir)
	fmt.Printf("  - %d dashboards (up to %d container nesting)\n", c.Dashboards, c.MaxContainerNest)
	fmt.Printf("  - %d queries (%d%% with cross-refs)\n", c.Queries, c.CrossRefPercent)
	fmt.Printf("  - %d controls (%d%% using query refs)\n", c.Controls, c.CrossRefPercent)
	fmt.Printf("  - %d benchmarks (up to %d depth)\n", c.Benchmarks, c.MaxBenchDepth)
}

func generateQueries(dir string, c Config) {
	var sb strings.Builder
	sb.WriteString("# Generated queries for lazy loading testing\n\n")

	for i := 0; i < c.Queries; i++ {
		// Vary the query patterns
		switch i % 4 {
		case 0:
			// Simple query
			sb.WriteString(fmt.Sprintf(`query "query_%d" {
  title       = "Query %d"
  description = "Simple query for lazy loading test"
  sql         = "SELECT %d as id, 'query_%d' as name"
  tags        = local.common_tags
}

`, i, i, i, i))
		case 1:
			// Query with heredoc SQL
			sb.WriteString(fmt.Sprintf(`query "query_%d" {
  title       = "Query %d with Heredoc"
  description = "Query using heredoc SQL syntax"
  sql = <<-EOQ
    SELECT
      %d as id,
      'query_%d' as name,
      'heredoc' as style,
      now() as generated_at
    FROM generate_series(1, 10)
  EOQ
  tags = local.common_tags
}

`, i, i, i, i))
		case 2:
			// Parameterized query
			sb.WriteString(fmt.Sprintf(`query "query_%d" {
  title       = "Parameterized Query %d"
  description = "Query with parameters for lazy loading test"
  sql         = "SELECT * FROM data WHERE region = $1 AND count > $2"

  param "region" {
    description = "Region filter"
    default     = var.test_region
  }

  param "threshold" {
    description = "Threshold value"
    default     = var.test_threshold
  }
}

`, i, i))
		case 3:
			// Control-ready query
			sb.WriteString(fmt.Sprintf(`query "query_%d" {
  title       = "Control Query %d"
  description = "Query formatted for control usage"
  sql         = "SELECT 'pass' as status, 'resource_%d' as resource, 'Query %d check passed' as reason"
}

`, i, i, i, i))
		}
	}

	writeFile(filepath.Join(dir, "queries.pp"), sb.String())
}

func generateControls(dir string, c Config) {
	var sb strings.Builder
	sb.WriteString("# Generated controls for lazy loading testing\n\n")

	severities := []string{"low", "medium", "high", "critical"}

	for i := 0; i < c.Controls; i++ {
		severity := severities[i%len(severities)]

		// Decide if this control references a query
		usesQueryRef := rand.Intn(100) < c.CrossRefPercent

		if usesQueryRef && c.Queries > 0 {
			// Control referencing a query
			queryIdx := (i * 3) % c.Queries // Deterministic but varied
			sb.WriteString(fmt.Sprintf(`control "control_%d" {
  title       = "Control %d (Query Ref)"
  description = "Control referencing query.query_%d"
  query       = query.query_%d
  severity    = "%s"
  tags        = local.common_tags
}

`, i, i, queryIdx, queryIdx, severity))
		} else {
			// Control with inline SQL
			sb.WriteString(fmt.Sprintf(`control "control_%d" {
  title       = "Control %d (Inline)"
  description = "Control with inline SQL"
  sql         = "SELECT 'pass' as status, 'resource_%d' as resource, 'Control %d passed' as reason"
  severity    = "%s"
  tags        = local.common_tags
}

`, i, i, i, i, severity))
		}
	}

	writeFile(filepath.Join(dir, "controls.pp"), sb.String())
}

func generateBenchmarks(dir string, c Config) {
	var sb strings.Builder
	sb.WriteString("# Generated benchmarks for lazy loading testing\n\n")

	controlsPerBenchmark := c.Controls / c.Benchmarks
	if controlsPerBenchmark < 1 {
		controlsPerBenchmark = 1
	}

	// Create flat benchmarks first
	flatBenchmarks := c.Benchmarks / 2
	if flatBenchmarks < 1 {
		flatBenchmarks = 1
	}

	for i := 0; i < flatBenchmarks; i++ {
		children := make([]string, 0, controlsPerBenchmark)
		for j := 0; j < controlsPerBenchmark && (i*controlsPerBenchmark+j) < c.Controls; j++ {
			children = append(children, fmt.Sprintf("control.control_%d", i*controlsPerBenchmark+j))
		}
		sb.WriteString(fmt.Sprintf(`benchmark "benchmark_%d" {
  title       = "Flat Benchmark %d"
  description = "Flat benchmark with %d controls"
  children = [
    %s
  ]
  tags = local.common_tags
}

`, i, i, len(children), strings.Join(children, ",\n    ")))
	}

	// Create nested benchmarks
	// Build a hierarchy: root -> level_1 -> level_2 -> ... -> controls
	for depth := 1; depth <= c.MaxBenchDepth; depth++ {
		benchIdx := flatBenchmarks + depth - 1
		if benchIdx >= c.Benchmarks {
			break
		}

		if depth == 1 {
			// Root of hierarchy
			sb.WriteString(fmt.Sprintf(`benchmark "nested_root" {
  title       = "Nested Hierarchy Root"
  description = "Root of %d-level deep hierarchy"
  children = [
    benchmark.nested_level_1
  ]
  tags = local.common_tags
}

`, c.MaxBenchDepth))
			// Also generate nested_level_1 if there are more levels
			if c.MaxBenchDepth > 1 {
				sb.WriteString(`benchmark "nested_level_1" {
  title       = "Nested Level 1"
  description = "Intermediate level in hierarchy"
  children = [
    benchmark.nested_level_2
  ]
}

`)
			}
		} else if depth < c.MaxBenchDepth {
			// Intermediate level
			sb.WriteString(fmt.Sprintf(`benchmark "nested_level_%d" {
  title       = "Nested Level %d"
  description = "Intermediate level in hierarchy"
  children = [
    benchmark.nested_level_%d
  ]
}

`, depth, depth, depth+1))
		} else {
			// Leaf level with controls
			leafControls := make([]string, 0, 3)
			for j := 0; j < 3 && (c.Controls-3+j) < c.Controls; j++ {
				leafControls = append(leafControls, fmt.Sprintf("control.control_%d", c.Controls-3+j))
			}
			sb.WriteString(fmt.Sprintf(`benchmark "nested_level_%d" {
  title       = "Nested Level %d (Leaf)"
  description = "Deepest level with controls"
  children = [
    %s
  ]
}

`, depth, depth, strings.Join(leafControls, ",\n    ")))
		}
	}

	// Create a wide benchmark
	wideChildren := make([]string, 0, 20)
	for i := 0; i < 20 && i < flatBenchmarks; i++ {
		wideChildren = append(wideChildren, fmt.Sprintf("benchmark.benchmark_%d", i))
	}
	if len(wideChildren) > 0 {
		sb.WriteString(fmt.Sprintf(`benchmark "wide_root" {
  title       = "Wide Benchmark Root"
  description = "Wide benchmark with many child benchmarks"
  children = [
    %s
  ]
  tags = local.common_tags
}

`, strings.Join(wideChildren, ",\n    ")))
	}

	writeFile(filepath.Join(dir, "benchmarks.pp"), sb.String())
}

func generateDashboards(dir string, c Config) {
	var sb strings.Builder
	sb.WriteString("# Generated dashboards for lazy loading testing\n\n")

	for i := 0; i < c.Dashboards; i++ {
		nestingDepth := (i % c.MaxContainerNest) + 1
		queryIdx := i % c.Queries

		sb.WriteString(fmt.Sprintf(`dashboard "dashboard_%d" {
  title       = "Dashboard %d"
  description = "Dashboard with %d levels of container nesting"
  tags        = local.common_tags

`, i, i, nestingDepth))

		// Generate nested containers
		generateNestedContainer(&sb, nestingDepth, queryIdx, c.Queries, 1)

		sb.WriteString("}\n\n")
	}

	writeFile(filepath.Join(dir, "dashboards.pp"), sb.String())
}

func generateNestedContainer(sb *strings.Builder, maxDepth, queryIdx, totalQueries, currentDepth int) {
	indent := strings.Repeat("  ", currentDepth)

	sb.WriteString(fmt.Sprintf("%scontainer {\n", indent))
	sb.WriteString(fmt.Sprintf("%s  title = \"Level %d Container\"\n\n", indent, currentDepth))

	// Add a card at each level
	sb.WriteString(fmt.Sprintf("%s  card {\n", indent))
	sb.WriteString(fmt.Sprintf("%s    title = \"Level %d Card\"\n", indent, currentDepth))
	sb.WriteString(fmt.Sprintf("%s    width = 4\n", indent))
	sb.WriteString(fmt.Sprintf("%s    sql   = query.query_%d.sql\n", indent, (queryIdx+currentDepth)%totalQueries))
	sb.WriteString(fmt.Sprintf("%s  }\n\n", indent))

	if currentDepth < maxDepth {
		generateNestedContainer(sb, maxDepth, queryIdx, totalQueries, currentDepth+1)
	} else {
		// Add a table at the deepest level
		sb.WriteString(fmt.Sprintf("%s  table {\n", indent))
		sb.WriteString(fmt.Sprintf("%s    title = \"Data Table\"\n", indent))
		sb.WriteString(fmt.Sprintf("%s    sql   = query.query_%d.sql\n", indent, queryIdx))
		sb.WriteString(fmt.Sprintf("%s  }\n", indent))
	}

	sb.WriteString(fmt.Sprintf("%s}\n", indent))
}

func writeFile(path, content string) {
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		fmt.Printf("Failed to write %s: %v\n", path, err)
		os.Exit(1)
	}
}
