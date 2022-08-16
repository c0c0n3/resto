package bytez

import (
	"bytes"
	"io"
	"reflect"
	"testing"
)

const dataSize = 4 * 1024

func makeData() []byte {
	data := make([]byte, dataSize)
	for k := 0; k < dataSize; k++ {
		data[k] = byte(k % 256)
	}
	return data
}

func checkData(t *testing.T, data []byte) {
	if len(data) != dataSize {
		t.Errorf("want size: %d; got: %d", dataSize, len(data))
	}
	for k := 0; k < len(data); k++ {
		want := byte(k % 256)
		if data[k] != want {
			t.Errorf("[%d] want: %d; got: %d", k, want, data[k])
		}
	}
}

func writeAll(dest io.WriteCloser) {
	defer dest.Close()
	src := bytes.NewBuffer(makeData())
	io.Copy(dest, src)
}

func readAll(src io.ReadCloser) []byte {
	defer src.Close()
	data, _ := io.ReadAll(src)
	return data
}

func TestWriteThenRead(t *testing.T) {
	buf := NewBuffer()
	writeAll(buf)
	got := readAll(buf)
	checkData(t, got)
}

func TestBytesBeforeAnyRead(t *testing.T) {
	buf := NewBuffer()
	buf.Write([]byte{1, 2, 3})

	got := buf.Bytes()
	want := []byte{1, 2, 3}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("want: %v; got: %v", want, got)
	}

	got, _ = io.ReadAll(buf)
	if !reflect.DeepEqual(got, want) {
		t.Errorf("want: %v; got: %v", want, got)
	}
}

func TestBytesAfterRead(t *testing.T) {
	buf := NewBuffer()
	buf.Write([]byte{1, 2, 3})

	firstTwo := make([]byte, 2)
	buf.Read(firstTwo)

	want := []byte{1, 2}
	if !reflect.DeepEqual(firstTwo, want) {
		t.Errorf("want: %v; got: %v", want, firstTwo)
	}

	got := buf.Bytes()
	want = []byte{3}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("want: %v; got: %v", want, got)
	}
}

func TestNewBufferFromNil(t *testing.T) {
	got := NewBufferFrom(nil)
	if got == nil {
		t.Fatalf("want: buffer; got: nil")
	}

	gotLen := len(got.Bytes())
	if gotLen != 0 {
		t.Errorf("want: 0; got: %d", gotLen)
	}
}

func TestNewBufferCopy(t *testing.T) {
	original := []byte{1}
	got := NewBufferFrom(original)
	if got == nil {
		t.Fatalf("want: buffer; got: nil")
	}

	original[0] = 2
	gotCopy := got.Bytes()[0]
	if gotCopy != 1 {
		t.Errorf("want: 1; got: %d", gotCopy)
	}
}

func TestReaderFromNil(t *testing.T) {
	got := Reader(nil)
	if got == nil {
		t.Fatalf("want: buffer; got: nil")
	}

	gotBuf, err := io.ReadAll(got)
	if err != nil {
		t.Fatalf("want: buffer; got: %v", err)
	}
	gotLen := len(gotBuf)
	if gotLen != 0 {
		t.Errorf("want: 0; got: %d", gotLen)
	}
}

func TestReaderBufferSharing(t *testing.T) {
	original := []byte{1}
	got := Reader(original)
	if got == nil {
		t.Fatalf("want: buffer; got: nil")
	}

	original[0] = 2
	gotBuf, err := io.ReadAll(got)
	if err != nil {
		t.Fatalf("want: buffer; got: %v", err)
	}
	if gotBuf[0] != 2 {
		t.Errorf("want: 2; got: %d", gotBuf[0])
	}
}
