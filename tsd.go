package tsd

import (
	"bytes"
	"encoding/binary"
	"io"
)

type ChunkID = uint64

// TSDWriter wraps an io.Writer and facilitates writing TSD Chunks
type TSDWriter struct {
	w   io.Writer
	buf []byte
}

func NewTSDWriter(w io.Writer) *TSDWriter {
	return &TSDWriter{
		w: w,
		// random pre-allocated buffer for copies, speeds up allocations
		buf: make([]byte, 32*1024),
	}
}

// Write a small chunk with
func (t TSDWriter) Write(id ChunkID, data []byte) error {
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

// getting a little clever here to save space with the streamer
// 16383 is the largest number that can be represented by 2 bytes as a varint
const streamingChunkSize int = 16383
const continuationChunkID = ChunkID(0)

// Stream in a long chunk with undefined length
func (t TSDWriter) Stream(id ChunkID, r io.Reader) error {
	chunk := make([]byte, streamingChunkSize)
	for i := 0; ; i++ {
		n, err := io.ReadFull(r, chunk)
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		if i == 0 {
			err = t.Write(id, chunk)
		} else {
			err = t.Write(continuationChunkID, chunk)
		}
		if err != nil {
			return err
		}
		if n != streamingChunkSize {
			return nil
		}
	}
}
