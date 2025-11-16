package common

import (
	"os"
	"unicode/utf8"
)

type RuneReader struct {
	buf      []byte
	off      int
	line     int
	col      int
	filePath string
}

func NewRuneReaderFromFile(path string) (*RuneReader, error) {
	b, err := os.ReadFile(path)

	if err != nil {
		return nil, err
	}

	return &RuneReader{buf: b, line: 1, col: 1, filePath: path}, nil
}

func (r *RuneReader) EOF() bool {
	return r.off >= len(r.buf)
}

func (r *RuneReader) Offset() int {
	return r.off
}

func (r *RuneReader) Pos() (int, int) {
	return r.line, r.col
}

func (r *RuneReader) FilePath() string {
	return r.filePath
}

func (r *RuneReader) Slice(start, end int) string {
	if start < 0 {
		start = 0
	}

	if end > len(r.buf) {
		end = len(r.buf)
	}

	if start > end {
		start, end = end, start
	}

	return string(r.buf[start:end])
}

type Snapshot struct{ off, line, col int } // Buat nyimpen posisi baca saat ini

func (r *RuneReader) Snapshot() Snapshot {
	return Snapshot{r.off, r.line, r.col}
}

func (r *RuneReader) Restore(s Snapshot) {
	r.off, r.line, r.col = s.off, s.line, s.col
}

func (r *RuneReader) Seek(off int) {
	r.off = off
}

func (r *RuneReader) Peek() rune {
	if r.EOF() {
		return 0
	}

	ch, _ := utf8.DecodeRune(r.buf[r.off:])
	return ch
}

func (r *RuneReader) Read() (rune, bool) {
	if r.EOF() {
		return 0, false
	}

	ch, w := utf8.DecodeRune(r.buf[r.off:])
	r.off += w

	switch ch {
	case '\r':
		if !r.EOF() && r.buf[r.off] == '\n' {
			r.off++
		}

		r.line++
		r.col = 1
	case '\n':
		r.line++
		r.col = 1
	default:
		r.col++
	}

	return ch, true
}
