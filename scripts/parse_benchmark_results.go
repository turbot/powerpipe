//go:build ignore

package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strconv"
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
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "Usage: go run parse_benchmark_results.go <results_file>")
		os.Exit(1)
	}

	file, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	// Regex to parse benchmark output lines
	// Example: BenchmarkLoadWorkspace_Small-10   50   23456789 ns/op   12345678 B/op   123456 allocs/op
	benchRegex := regexp.MustCompile(`^(Benchmark\S+)-(\d+)\s+(\d+)\s+([\d.]+)\s+ns/op(?:\s+(\d+)\s+B/op)?(?:\s+(\d+)\s+allocs/op)?`)

	report := BenchmarkReport{
		Timestamp: extractTimestamp(os.Args[1]),
		Results:   []BenchmarkResult{},
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		matches := benchRegex.FindStringSubmatch(line)
		if matches != nil {
			result := BenchmarkResult{
				Name:       matches[1],
				Iterations: parseInt(matches[3]),
				NsPerOp:    parseFloat(matches[4]),
			}
			result.MsPerOp = result.NsPerOp / 1e6

			if matches[5] != "" {
				result.BytesPerOp = parseInt64(matches[5])
			}
			if matches[6] != "" {
				result.AllocsPerOp = parseInt64(matches[6])
			}

			report.Results = append(report.Results, result)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		os.Exit(1)
	}

	output, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error marshaling JSON: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(string(output))
}

func extractTimestamp(filename string) string {
	// Extract timestamp from filename like benchmark_20231215_143022.txt
	parts := strings.Split(filename, "_")
	if len(parts) >= 3 {
		return strings.TrimSuffix(parts[len(parts)-2]+"_"+parts[len(parts)-1], ".txt")
	}
	return ""
}

func parseInt(s string) int {
	v, _ := strconv.Atoi(s)
	return v
}

func parseInt64(s string) int64 {
	v, _ := strconv.ParseInt(s, 10, 64)
	return v
}

func parseFloat(s string) float64 {
	v, _ := strconv.ParseFloat(s, 64)
	return v
}
