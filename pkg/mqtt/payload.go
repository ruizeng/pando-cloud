package mqtt

import (
	"bytes"
	"io"
)

// Payload is the interface for Publish payloads. Typically the BytesPayload
// implementation will be sufficient for small payloads whose full contents
// will exist in memory. However, other implementations can read or write
// payloads requiring them holding their complete contents in memory.
type Payload interface {
	// Size returns the number of bytes that WritePayload will write.
	Size() int

	// WritePayload writes the payload data to w. Implementations must write
	// Size() bytes of data, but it is *not* required to do so prior to
	// returning. Size() bytes must have been written to w prior to another
	// message being encoded to the underlying connection.
	WritePayload(b *bytes.Buffer) error

	// ReadPayload reads the payload data from r (r will EOF at the end of the
	// payload). It is *not* required for r to have been consumed prior to this
	// returning. r must have been consumed completely prior to another message
	// being decoded from the underlying connection.
	ReadPayload(r io.Reader, n int) error
}

type BytesPayload []byte

func (p BytesPayload) Size() int {
	return len(p)
}

func (p BytesPayload) WritePayload(b *bytes.Buffer) error {
	_, err := b.Write(p)
	return err
}

func (p BytesPayload) ReadPayload(r io.Reader, n int) error {
	_, err := io.ReadFull(r, p)
	return err
}
