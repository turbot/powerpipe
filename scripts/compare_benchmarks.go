//go:build ignore

package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type BenchmarkResult struct {
	Name        string  `json:"name"`
	Iterations  int     `json:"iterations"`
	NsPerOp     float64 `json:"ns_per_op"`
	MsPerOp     float64 `json:"ms_per_op"`
	BytesPerOp  int64   `json:"bytes_per_op,omitempty"`
	AllocsPerOp int64   `json:"allocs_per_op,omitempty"`
}

type BenchmarkReport struct {
	Timestamp string            `json:"timestamp"`
	Results   []BenchmarkResult `json:"results"`
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: go run compare_benchmarks.go <before.json> <after.json>")
		fmt.Println("")
		fmt.Println("Compares two benchmark result files and shows performance changes.")
		os.Exit(1)
	}

	before := loadResults(os.Args[1])
	after := loadResults(os.Args[2])

	fmt.Println("=== Benchmark Comparison ===")
	fmt.Printf("Before: %s\n", os.Args[1])
	fmt.Printf("After:  %s\n", os.Args[2])
	fmt.Println("")

	fmt.Printf("%-50s %12s %12s %10s\n", "Benchmark", "Before", "After", "Change")
	fmt.Println(strings.Repeat("-", 90))

	totalBefore := 0.0
	totalAfter := 0.0
	count := 0

	for name, b := range before {
		if a, ok := after[name]; ok {
			change := (a.NsPerOp - b.NsPerOp) / b.NsPerOp * 100
			changeStr := fmt.Sprintf("%+.1f%%", change)

			// Color coding hint
			indicator := ""
			if change < -5 {
				indicator = " ✓" // improvement
			} else if change > 5 {
				indicator = " ✗" // regression
			}

			fmt.Printf("%-50s %10.2fms %10.2fms %10s%s\n",
				name,
				b.NsPerOp/1e6,
				a.NsPerOp/1e6,
				changeStr,
				indicator)

			totalBefore += b.NsPerOp
			totalAfter += a.NsPerOp
			count++
		}
	}

	if count > 0 {
		fmt.Println(strings.Repeat("-", 90))
		overallChange := (totalAfter - totalBefore) / totalBefore * 100
		fmt.Printf("%-50s %10.2fms %10.2fms %+9.1f%%\n",
			"TOTAL",
			totalBefore/1e6,
			totalAfter/1e6,
			overallChange)
	}

	// Show memory comparison
	fmt.Println("")
	fmt.Println("=== Memory Comparison ===")
	fmt.Printf("%-50s %12s %12s %10s\n", "Benchmark", "Before", "After", "Change")
	fmt.Println(strings.Repeat("-", 90))

	for name, b := range before {
		if a, ok := after[name]; ok {
			if b.BytesPerOp > 0 && a.BytesPerOp > 0 {
				change := float64(a.BytesPerOp-b.BytesPerOp) / float64(b.BytesPerOp) * 100
				fmt.Printf("%-50s %10s %10s %+9.1f%%\n",
					name,
					formatBytes(b.BytesPerOp),
					formatBytes(a.BytesPerOp),
					change)
			}
		}
	}

	// Report any benchmarks that only exist in one file
	fmt.Println("")

	newBenchmarks := []string{}
	for name := range after {
		if _, ok := before[name]; !ok {
			newBenchmarks = append(newBenchmarks, name)
		}
	}
	if len(newBenchmarks) > 0 {
		fmt.Println("New benchmarks (not in before):")
		for _, name := range newBenchmarks {
			fmt.Printf("  + %s\n", name)
		}
	}

	removedBenchmarks := []string{}
	for name := range before {
		if _, ok := after[name]; !ok {
			removedBenchmarks = append(removedBenchmarks, name)
		}
	}
	if len(removedBenchmarks) > 0 {
		fmt.Println("Removed benchmarks (not in after):")
		for _, name := range removedBenchmarks {
			fmt.Printf("  - %s\n", name)
		}
	}
}

func loadResults(path string) map[string]BenchmarkResult {
	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading %s: %v\n", path, err)
		os.Exit(1)
	}

	var report BenchmarkReport
	if err := json.Unmarshal(data, &report); err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing %s: %v\n", path, err)
		os.Exit(1)
	}

	results := make(map[string]BenchmarkResult)
	for _, r := range report.Results {
		results[r.Name] = r
	}
	return results
}

func formatBytes(b int64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(b)/float64(div), "KMGTPE"[exp])
}
