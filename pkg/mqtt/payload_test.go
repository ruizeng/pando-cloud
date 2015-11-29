package mqtt

import (
        "bytes"
        "errors"
        "reflect"
        "testing"
)

type TestReadWriter struct {
        buf []byte
        r   int
        w   int
}

func (rw *TestReadWriter) Read(p []byte) (int, error) {
        n := len(p)
        if rw.w == rw.r {
                return 0, errors.New("no data")
        }
        if rw.w-rw.r < n {
                n = rw.w - rw.r
        }

        copy(p[0:n], rw.buf[rw.r:rw.r+n])
        rw.r += n

        return n, nil
}

func (rw *TestReadWriter) Write(p []byte) (int, error) {
        n := len(p)
        copy(rw.buf[rw.w:rw.w+n], p[0:n])
        rw.w += n

        return n, nil
}

func (rw *TestReadWriter) Size() int {
        return rw.w-rw.r
}

func TestPayload(t *testing.T) {
	bytesValue := []byte{'a', 'c', 'd', 'g'}
	payload := BytesPayload(bytesValue)

	rw := &TestReadWriter{
                w:   0,
                r:   0,
                buf: make([]byte, 10000),
        }	

	buf := new(bytes.Buffer)
        err := payload.WritePayload(buf)
        if err != nil {
                t.Error(err)
        }
	_, err = rw.Write(buf.Bytes())
        if err != nil {
                t.Error(err)
        }

        if payload.Size() != rw.Size() {
                t.Error("size is not correct", payload.Size(), rw.Size())
        }

	newPayload := make(BytesPayload, payload.Size())
	err = newPayload.ReadPayload(rw, payload.Size())
        if err != nil {
                t.Error(err)
        }

	if !reflect.DeepEqual(payload, newPayload) {
                t.Errorf("the origin:\n%x\n, now:\n%x\n", payload, newPayload)
        }
}
