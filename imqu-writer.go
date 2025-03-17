package imqu

import (
	"encoding/binary"
	"os"
	"time"
)

const Ext = ".imqu"

// ImQuWriter is a structure required to write data to the file
type ImQuWriter struct {
	fn string
	f  *os.File

	prevImgTime time.Time
	count       uint64
}

// NewImQuWriter create file and initialize ImQuWriter instance
func NewImQuWriter(fn string) (*ImQuWriter, error) {
	f, err := os.Create(fn)
	if err != nil {
		return nil, err
	}

	return &ImQuWriter{fn: fn, f: f}, nil
}

// AddData write any bytes to the file with time offset until previous write
func (w *ImQuWriter) AddData(b []byte) error {
	return w.AddDataWithTime(b, time.Now())
}

// AddDataWithTime write bytes to the file with specific time
func (w *ImQuWriter) AddDataWithTime(b []byte, t time.Time) error {
	w.writeWait(t)
	w.writeUint64(uint64(len(b)))

	if _, err := w.f.Write(b); err != nil {
		return err
	}

	w.prevImgTime = t
	w.count++

	return nil
}

func (w *ImQuWriter) writeWait(t time.Time) error {
	if w.prevImgTime.IsZero() {
		return w.writeUint64(0)
	}

	return w.writeUint64(uint64(t.Sub(w.prevImgTime)))
}

func (w *ImQuWriter) writeUint64(n uint64) error {
	b := make([]byte, 8)
	if n > 0 {
		binary.BigEndian.PutUint64(b, n)
	}

	_, err := w.f.Write(b)
	return err
}

// Close the writer, write count to tail of the file and close the file
func (w *ImQuWriter) Close() error {
	binary.Write(w.f, binary.BigEndian, w.count)
	return w.f.Close()
}
