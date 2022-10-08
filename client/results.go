package client

import (
	// "compress/gzip"
	"io"
	"log"
	"os"

	"github.com/glycerine/rbuf"
	// "github.com/klauspost/compress/zstd"

	"grpc-test/constants"
	api "grpc-test/gen/proto/go/results/api/v1"
)

type ResultsReader struct {
	buf  *rbuf.AtomicFixedSizeRingBuf
	done bool
}

func (r *ResultsReader) GetChunks(stream api.TestService_GetResultsClient) {
	for {
		chunk, err := stream.Recv()
		if err == io.EOF {
			r.done = true
			break
		}
		if err != nil {
			log.Fatalf("client.GetResults failed: %v", err)
		}
		r.buf.Write(chunk.Data[:chunk.Size])
	}
}

func NewResultsReader() *ResultsReader {
	buf := rbuf.NewAtomicFixedSizeRingBuf(constants.RingBufferSize)
	return &ResultsReader{
		buf: buf,
	}
}

func (r *ResultsReader) Read(p []byte) (n int, err error) {
	if r.done && r.buf.Readable() == 0 {
		return 0, io.EOF
	}
	n, err = r.buf.Read(p)
	if err != nil && err != io.EOF {
		log.Fatalf("failed to read from rbuf: %v\n", err)
	}
	return n, nil
}

func WriteReaderToFile(r io.Reader) {
	buf := make([]byte, constants.ReaderToFileBufferSize)
	file, err := os.OpenFile(
		"files/output",
		os.O_WRONLY|os.O_TRUNC|os.O_CREATE,
		0644,
	)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	for {
		readBytes, err := r.Read(buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("failed to read: %v", err)
			return
		}

		writtenBytes, err := file.Write(buf[:readBytes])
		if err != nil {
			log.Fatal(err)
		}
		if writtenBytes != readBytes {
			log.Fatalf("mismatch in written and read bytes: %d != %d", writtenBytes, readBytes)
		}
	}
}

func (c *Client) GetResults() io.Reader {
	stream, err := c.client.GetResults(c.ctx, &api.GetResultsRequest{})
	if err != nil {
		log.Fatalf("grpc.GetResults failed: %v", err)
	}

	reader := NewResultsReader()
	go reader.GetChunks(stream)

	// decompressedReader, err := zstd.NewReader(reader)
	// decompressedReader, err := gzip.NewReader(reader)
	// if err != nil {
	//     log.Fatalf("failed to create decompressed reader: %v", err)
	// }
	return reader
}
