package tsd

import (
	"bytes"
	"encoding/binary"
	"io"
	"io/ioutil"
)

type ChunkID = uint64

// TSDWriter wraps an io.Writer and facilitates writing TSD Chunks
type TSDWriter struct {
	w   io.Writer
	buf []byte
}

// Create a new TSD writer to start writing in TSD format
func NewTSDWriter(w io.Writer) *TSDWriter {
	return &TSDWriter{
		w: w,
		// pre-allocate buffer to speed up copies
		buf: make([]byte, 32*1024),
	}
}

// Write a small chunk with
func (t *TSDWriter) Write(id ChunkID, data []byte) error {
	buf := make([]byte, binary.MaxVarintLen64+binary.MaxVarintLen32)
	n := binary.PutUvarint(buf, id)
	n += binary.PutUvarint(buf[n:], uint64(len(data)))
	_, err := io.CopyBuffer(t.w, bytes.NewBuffer(buf[:n]), t.buf)
	if err != nil {
		return err
	}
	_, err = io.CopyBuffer(t.w, bytes.NewBuffer(data), t.buf)
	return err
}


type TSDReader struct {
	r ByteReaderReader
	prev *io.LimitedReader
}

type ByteReaderReader interface {
	io.ByteReader
	io.Reader
}

func NewTSDReader(r ByteReaderReader) *TSDReader {
	return &TSDReader{r: r}
}

// Next gets the next chunk in the TSD stream. Client can get the
func (t *TSDReader) Next() (ChunkID, io.Reader, error) {
	if t.prev != nil && t.prev.N > 0 {
		_, err := io.Copy(ioutil.Discard, t.prev)
		if err != nil {
			return 0, nil, err
		}
	}
	id, err := binary.ReadUvarint(t.r)
	if err != nil {
		return 0, nil, err
	}
	ln, err := binary.ReadUvarint(t.r)
	if err != nil {
		return 0, nil, err
	}
	t.prev = &io.LimitedReader{R: t.r, N: int64(ln)}
	return id, t.prev, nil
}
