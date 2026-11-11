package controldisplay

import (
	"strings"
	"testing"
	"time"

	"github.com/turbot/go-kit/helpers"
	"github.com/turbot/powerpipe/internal/controlexecute"
)

func renderDimensionsWithTimeout(t *testing.T, dimensions []controlexecute.Dimension, width int) string {
	t.Helper()
	ControlColors, _ = NewControlColorScheme(ColorSchemes["plain"])
	colorGenerator, err := controlexecute.NewDimensionColorGenerator(4, 27)
	if err != nil {
		t.Fatal(err)
	}

	done := make(chan string, 1)
	go func() {
		done <- NewDimensionsRenderer(dimensions, colorGenerator, width).Render()
	}()
	select {
	case output := <-done:
		return output
	case <-time.After(5 * time.Second):
		t.Fatalf("DimensionsRenderer.Render did not return for width %d", width)
		return ""
	}
}

func TestDimensionsRenderMultiLineValue(t *testing.T) {
	// two lines that are each shorter than the combined printable length (#1146)
	value := strings.Repeat("b ", 24) + "b" + "\n" + strings.Repeat("c ", 33) + "c"
	dimensions := []controlexecute.Dimension{{Key: "note", Value: value}}

	for _, width := range []int{1, 10, 40, 80, 116, 200} {
		output := renderDimensionsWithTimeout(t, dimensions, width)
		if strings.ContainsAny(output, "\r\n") {
			t.Errorf("width %d: output contains a line break: %q", width, output)
		}
		if length := helpers.PrintableLength(output); length > width {
			t.Errorf("width %d: output is %d characters wide: %q", width, length, output)
		}
	}
}

func TestDimensionsRenderLineBreaksAsSpaces(t *testing.T) {
	dimensions := []controlexecute.Dimension{{Key: "note", Value: "first line\r\nsecond\nthird\rfourth"}}

	output := renderDimensionsWithTimeout(t, dimensions, 100)

	if want := "first line second third fourth"; helpers.Clean(output) != want {
		t.Errorf("expected %q, got %q", want, helpers.Clean(output))
	}
}
