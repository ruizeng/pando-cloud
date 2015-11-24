package mqtt

import (
	"io"
)

type Payload interface {
	Size() int
	WritePayload(w io.Writer) error
	Readpayload(r io.Reader) error
}

type BytesPayload []byte

func (p BytesPayload) Size() int {
	return len(p)
}

func (p BytesPayload) WritePayload(w io.Writer) error {
	_, err := w.Write(p)

	return err
}

func (p BytesPayload) Readpayload(r io.Reader) error {
	_, err := io.ReadFull(r, p)

	return err
}
