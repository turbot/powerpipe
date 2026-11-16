package controlexecute

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"os"
	"sync"
)

// rowStore holds the result rows of control runs in a temporary file, so that they are not held in memory.
// The control runs write their rows concurrently, each in segments.
type rowStore struct {
	mu   sync.Mutex
	file *os.File
	size int64
}

// rowSegment is a range of the row store file
type rowSegment struct {
	offset, length int64
}

func newRowStore() (*rowStore, error) {
	f, err := os.CreateTemp("", "powerpipe-rows-*")
	if err != nil {
		return nil, err
	}
	// remove the name now: the open file stays usable, and nothing is left behind if the process stops
	if err := os.Remove(f.Name()); err != nil {
		return nil, errors.Join(err, f.Close())
	}
	return &rowStore{file: f}, nil
}

// write appends b to the file and returns its segment
func (s *rowStore) write(b []byte) (rowSegment, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, err := s.file.WriteAt(b, s.size); err != nil {
		return rowSegment{}, err
	}
	segment := rowSegment{offset: s.size, length: int64(len(b))}
	s.size += segment.length
	return segment, nil
}

func (s *rowStore) reader(segment rowSegment) io.Reader {
	return io.NewSectionReader(s.file, segment.offset, segment.length)
}

func (s *rowStore) close() error {
	return s.file.Close()
}

// rowSegmentSize is how many bytes of rows a control run buffers for each status before it writes them to the store
const rowSegmentSize = 64 << 10

// storedRows are the result rows of one control run in a rowStore.
// Each status has its own segments, so that the rows can be read back in status order.
type storedRows struct {
	store    *rowStore
	pending  map[string]*bytes.Buffer
	segments map[string][]rowSegment
	// complete is set when all rows have been added. Until then there are no rows to read.
	complete bool
	// err is the first error writing or reading rows
	err error
}

// storedRow is a ResultRow as it is written to the store, one JSON object per line
type storedRow struct {
	Reason     string      `json:"r"`
	Resource   string      `json:"s"`
	Dimensions [][3]string `json:"d,omitempty"`
}

func newStoredRows(store *rowStore) *storedRows {
	return &storedRows{
		store:    store,
		pending:  make(map[string]*bytes.Buffer),
		segments: make(map[string][]rowSegment),
	}
}

func (s *storedRows) add(row *ResultRow) {
	if s.err != nil {
		return
	}
	stored := storedRow{Reason: row.Reason, Resource: row.Resource}
	for _, d := range row.Dimensions {
		stored.Dimensions = append(stored.Dimensions, [3]string{d.Key, d.Value, d.SqlType})
	}
	line, err := json.Marshal(stored)
	if err != nil {
		s.err = err
		return
	}
	buf, ok := s.pending[row.Status]
	if !ok {
		buf = &bytes.Buffer{}
		s.pending[row.Status] = buf
	}
	buf.Write(line)
	buf.WriteByte('\n')
	if buf.Len() >= rowSegmentSize {
		s.flush(row.Status)
	}
}

// flush writes the pending rows of the status to the store
func (s *storedRows) flush(status string) {
	buf := s.pending[status]
	if s.err != nil || buf == nil || buf.Len() == 0 {
		return
	}
	segment, err := s.store.write(buf.Bytes())
	if err != nil {
		s.err = err
		return
	}
	s.segments[status] = append(s.segments[status], segment)
	buf.Reset()
}

// setComplete writes the remaining rows to the store, and makes the rows readable
func (s *storedRows) setComplete() {
	for status := range s.pending {
		s.flush(status)
	}
	s.pending = nil
	s.complete = true
}

// each calls f with each row, in the status order of resultStatusOrder, then in the order the rows were added
func (s *storedRows) each(run *ControlRun, f func(*ResultRow) error) error {
	if s.err != nil || !s.complete {
		return s.err
	}
	for _, status := range resultStatusOrder {
		for _, segment := range s.segments[status] {
			if err := s.eachInSegment(segment, run, status, f); err != nil {
				return err
			}
		}
	}
	return nil
}

func (s *storedRows) eachInSegment(segment rowSegment, run *ControlRun, status string, f func(*ResultRow) error) error {
	r := bufio.NewReader(s.store.reader(segment))
	for {
		line, err := r.ReadBytes('\n')
		if err == io.EOF && len(line) == 0 {
			return nil
		}
		if err != nil {
			s.err = err
			return err
		}
		var stored storedRow
		if err := json.Unmarshal(line, &stored); err != nil {
			s.err = err
			return err
		}
		row := &ResultRow{Reason: stored.Reason, Resource: stored.Resource, Status: status, Run: run, Control: run.Control}
		for _, d := range stored.Dimensions {
			row.Dimensions = append(row.Dimensions, Dimension{Key: d[0], Value: d[1], SqlType: d[2]})
		}
		if err := f(row); err != nil {
			return err
		}
	}
}
