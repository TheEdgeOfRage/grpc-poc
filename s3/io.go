package s3

import (
	"io"
	"log"

	"github.com/glycerine/rbuf"
)

type WriteAtReader struct {
	buf  *rbuf.AtomicFixedSizeRingBuf
	done bool
}

func NewWriteAtReader(maxBufSize int) *WriteAtReader {
	buf := rbuf.NewAtomicFixedSizeRingBuf(maxBufSize)
	return &WriteAtReader{
		buf: buf,
	}
}

func (wr *WriteAtReader) WriteAt(p []byte, pos int64) (n int, err error) {
	n, err = wr.buf.Write(p)
	return
}

func (wr *WriteAtReader) Read(p []byte) (n int, err error) {
	if wr.done && wr.buf.Readable() == 0 {
		return 0, io.EOF
	}
	n, err = wr.buf.Read(p)
	if err != nil && err != io.EOF {
		log.Fatalf("failed to read from rbuf: %v\n", err)
	}
	return n, nil
}

func (wr *WriteAtReader) Done() {
	wr.done = true
}
