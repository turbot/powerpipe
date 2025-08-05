package controldisplay

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"text/template"
	"time"

	"github.com/Masterminds/sprig/v3"
)

// templateFuncs merges desired functions from sprig with custom functions that we
// define in steampipe
func templateFuncs(renderContext TemplateRenderContext) template.FuncMap {
	useFromSprigMap := []string{"upper", "toJson", "quote", "dict", "add", "now", "toPrettyJson"}

	var funcs template.FuncMap = template.FuncMap{}
	sprigMap := sprig.TxtFuncMap()
	for _, use := range useFromSprigMap {
		f, found := sprigMap[use]
		if !found {
			// guarantee that when a function is expected to be present
			// it does not slip through any crack
			panic(fmt.Sprintf("%s not found", use))
		}
		if found {
			funcs[use] = f
		}
	}
	// custom steampipe functions - ones we couldn't find in sprig
	formatterTemplateFuncMap := template.FuncMap{
		"durationInSeconds": durationInSeconds,
		"toCsvCell":         toCSVCellFnFactory(renderContext.Config.Separator),
		"toSafeJson":        toSafeJson,
	}
	for k, v := range formatterTemplateFuncMap {
		funcs[k] = v
	}

	return funcs
}

// toSafeJson safely converts a value to JSON string, handling error cases gracefully
func toSafeJson(v interface{}) string {
	if v == nil {
		return "null"
	}

	// For strings, handle them specially to avoid JSON encoding issues
	if str, ok := v.(string); ok {
		// Clean the string to remove problematic characters that could break JSON
		// Replace newlines with spaces and escape any remaining problematic characters
		cleanedStr := strings.ReplaceAll(str, "\n", " ")
		cleanedStr = strings.ReplaceAll(cleanedStr, "\r", " ")
		cleanedStr = strings.ReplaceAll(cleanedStr, "\t", " ")
		// Remove any null bytes
		cleanedStr = strings.ReplaceAll(cleanedStr, "\x00", "")

		// Use json.Marshal to properly escape the cleaned string
		bytes, err := json.Marshal(cleanedStr)
		if err != nil {
			// If marshaling fails, return a safe fallback
			return `"Error: Unable to serialize error message"`
		}
		return string(bytes)
	}

	// For maps, handle them specially to ensure they're valid
	if m, ok := v.(map[string]interface{}); ok {
		// Create a safe copy of the map, filtering out any problematic values
		safeMap := make(map[string]interface{})
		for k, val := range m {
			if val != nil {
				safeMap[k] = val
			}
		}
		bytes, err := json.Marshal(safeMap)
		if err != nil {
			return "{}"
		}
		return string(bytes)
	}

	// For slices, handle them specially to ensure they're valid
	if slice, ok := v.([]interface{}); ok {
		// Create a safe copy of the slice, filtering out any problematic values
		safeSlice := make([]interface{}, 0, len(slice))
		for _, val := range slice {
			if val != nil {
				safeSlice = append(safeSlice, val)
			}
		}
		bytes, err := json.Marshal(safeSlice)
		if err != nil {
			return "[]"
		}
		return string(bytes)
	}

	// For other types, use the standard approach
	bytes, err := json.Marshal(v)
	if err != nil {
		// Try to convert to string as a fallback
		return fmt.Sprintf("%q", fmt.Sprintf("%v", v))
	}
	return string(bytes)
}

// toCsvCell escapes a value for csv
// we need to do this in a factory function, so that we can
// set the separator for the CSV writer for this render session
func toCSVCellFnFactory(comma string) func(interface{}) string {
	csvWriterBuffer := bytes.NewBufferString("")
	csvBufferLock := sync.Mutex{}

	csvWriter := csv.NewWriter(csvWriterBuffer)
	if len(comma) > 0 {
		csvWriter.Comma = []rune(comma)[0]
	}

	return func(v interface{}) string {
		csvBufferLock.Lock()
		defer csvBufferLock.Unlock()

		csvWriterBuffer.Reset()
		csvWriter.Write([]string{fmt.Sprintf("%v", v)}) //nolint:errcheck // TODO: fix this
		csvWriter.Flush()
		return strings.TrimSpace(csvWriterBuffer.String())
	}
}

// durationInSeconds returns the passed in duration as seconds
func durationInSeconds(t time.Duration) float64 { return t.Seconds() }
