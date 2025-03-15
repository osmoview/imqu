package imqu

import (
	"encoding/binary"
	"io"
	"os"
	"time"
)

// ImQuReader is a structure required to read data from the file
type ImQuReader struct {
	f *os.File

	time  time.Time
	count int
}

// OpenImQu open file and initialize ImQuReader instance
func OpenImQu(fn string) (*ImQuReader, error) {
	f, err := os.Open(fn)
	if err != nil {
		return nil, err
	}

	// read tail of the file
	if _, err := f.Seek(-8, io.SeekEnd); err != nil {
		return nil, err
	}

	tail := make([]byte, 8)
	if _, err := f.Read(tail); err != nil {
		return nil, err
	}

	if _, err := f.Seek(0, io.SeekStart); err != nil {
		return nil, err
	}

	count := binary.BigEndian.Uint64(tail)

	return &ImQuReader{f: f, count: int(count)}, nil
}

// GetCount returens a frames contains in this file
func (iq *ImQuReader) GetCount() int {
	return iq.count
}

// SeekFrame seek and returns a frame data nearest to specified duration offset
func (iq *ImQuReader) SeekFrame(dur time.Duration) ([]byte, error) {
	start := time.Now()

	for {
		p, err := iq.readPart()
		if err != nil {
			return nil, err
		}

		if time.Since(start) > dur {
			return p.data, nil
		}
	}
}

// NextFrame read, wait and returns next frame data, where speed is coeficient of the delay delimiter
func (iq *ImQuReader) NextFrame(speed int) ([]byte, error) {
	p, err := iq.readPart()
	if err != nil {
		return nil, err
	}
	if speed <= 0 {
		speed = 1
	}

	time.Sleep(p.wait/time.Duration(speed) - time.Since(iq.time))

	iq.time = time.Now()

	return p.data, nil
}

type frame struct {
	wait time.Duration // ms
	size uint64        // size of the data
	data []byte        // data
}

func (iq *ImQuReader) readPart() (frame, error) {
	wait, err := iq.readUint64()
	if err != nil {
		return frame{}, err
	}

	size, err := iq.readUint64()
	if err != nil {
		return frame{}, err
	}

	data := make([]byte, size)
	_, err = iq.f.Read(data)

	return frame{
		wait: time.Duration(wait),
		size: size,
		data: data,
	}, err
}

func (iq *ImQuReader) readUint64() (uint64, error) {
	b := make([]byte, 8)
	if _, err := iq.f.Read(b); err != nil {
		return 0, err
	}

	return binary.BigEndian.Uint64(b), nil
}

// Close ImQuReader file
func (iq *ImQuReader) Close() error {
	return iq.f.Close()
}
