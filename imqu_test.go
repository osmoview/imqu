package imqu

import (
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"testing"
	"time"
)

var testfile = "testdata/test" + Ext

func init() {
	log.SetFlags(log.Lshortfile)
	os.MkdirAll(filepath.Dir(testfile), 0777)
}

func TestImQuWriter(t *testing.T) {
	iq, err := NewImQuWriter(testfile)
	if err != nil {
		t.Fatal(err)
	}

	for i := 0; i < 5; i++ {
		s := fmt.Sprintf("raw%02x", time.Now().UnixNano())
		if err := iq.AddData([]byte(s)); err != nil {
			t.Error(err)
		}
		time.Sleep(100 * time.Millisecond)
	}

	for i := 0; i < 5; i++ {
		s := fmt.Sprintf("time%02x", time.Now().UnixNano())
		it := time.Now().Add(time.Duration(i) * 100 * time.Millisecond)
		if err := iq.AddDataWithTime([]byte(s), it); err != nil {
			t.Error(err)
		}
	}

	if err := iq.Close(); err != nil {
		t.Error(err)
	}

	b, err := os.ReadFile(testfile)
	if err != nil {
		t.Fatal(err)
	}

	count := binary.BigEndian.Uint64(b[len(b)-8:])
	if count != 10 {
		t.Errorf("count %d != %d", count, 10)
	}
}

func TestImQuReader(t *testing.T) {
	iq, err := OpenImQu(testfile)
	if err != nil {
		t.Fatal(err)
	}

	it := time.Now()
	for {
		data, err := iq.NextFrame(1)
		if err != nil {
			if err == io.EOF {
				break
			}
			t.Error(err)
			break
		}

		t.Log(time.Since(it), string(data))
	}

	iq.Close()
}
