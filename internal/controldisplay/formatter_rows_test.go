package controldisplay

import "testing"

// only formatters that read control rows one at a time may run with the rows on disk:
// every other formatter must have them loaded first
func TestStreamsRows(t *testing.T) {
	streaming := map[string]Formatter{"none": &NullFormatter{}, "snapshot": &SnapshotFormatter{}}
	for name, f := range streaming {
		if !StreamsRows(f) {
			t.Errorf("expected the %s formatter to stream rows", name)
		}
	}
	notStreaming := map[string]Formatter{"text": &TextFormatter{}, "template": &TemplateFormatter{}}
	for name, f := range notStreaming {
		if StreamsRows(f) {
			t.Errorf("expected the %s formatter to need all rows in memory", name)
		}
	}
}
