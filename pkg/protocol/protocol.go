package protocol

import (
	"bytes"
	"encoding/binary"
	"github.com/PandoCloud/pando-cloud/pkg/tlv"
)

type Payload interface {
	Marshal() ([]byte, error)
	UnMarshal([]byte) error
}

func (c *Command) Marshal() ([]byte, error) {
	buffer := new(bytes.Buffer)
	err := binary.Write(buffer, binary.BigEndian, c.Head)
	if err != nil {
		return nil, err
	}

	for _, param := range c.Params {
		err = binary.Write(buffer, binary.BigEndian, param.ToBinary())
		if err != nil {
			return nil, err
		}
	}

	return buffer.Bytes(), nil
}

func (c *Command) UnMarshal(buf []byte) error {
	n := len(buf)
	r := bytes.NewReader(buf)
	err := binary.Read(r, binary.BigEndian, &c.Head)
	if err != nil {
		return err
	}
	c.Params = []tlv.TLV{}
	for i := binary.Size(c.Head); i < n; {
		tlv := tlv.TLV{}
		tlv.FromBinary(r)
		i += int(tlv.Length())
		c.Params = append(c.Params, tlv)
	}

	return nil
}

func (e *Event) Marshal() ([]byte, error) {
	buffer := new(bytes.Buffer)
	err := binary.Write(buffer, binary.BigEndian, e.Head)
	if err != nil {
		return nil, err
	}

	for _, param := range e.Params {
		err = binary.Write(buffer, binary.BigEndian, param.ToBinary())
		if err != nil {
			return nil, err
		}
	}

	return buffer.Bytes(), nil
}

func (e *Event) UnMarshal(buf []byte) error {
	n := len(buf)
	r := bytes.NewReader(buf)
	err := binary.Read(r, binary.BigEndian, &e.Head)
	if err != nil {
		return err
	}
	e.Params = []tlv.TLV{}
	for i := binary.Size(e.Head); i < n; {
		tlv := tlv.TLV{}
		tlv.FromBinary(r)
		i += int(tlv.Length())
		e.Params = append(e.Params, tlv)
	}

	return nil
}

func (d *Data) Marshal() ([]byte, error) {
	buffer := new(bytes.Buffer)
	err := binary.Write(buffer, binary.BigEndian, d.Head)
	if err != nil {
		return nil, err
	}

	for _, sub := range d.SubData {
		err = binary.Write(buffer, binary.BigEndian, sub.Head)
		if err != nil {
			return nil, err
		}
		for _, param := range sub.Params {
			err = binary.Write(buffer, binary.BigEndian, param.ToBinary())
			if err != nil {
				return nil, err
			}
		}
	}

	return buffer.Bytes(), nil
}

func (d *Data) UnMarshal(buf []byte) error {
	n := len(buf)
	r := bytes.NewReader(buf)
	err := binary.Read(r, binary.BigEndian, &d.Head)
	if err != nil {
		return err
	}
	d.SubData = []SubData{}
	for i := binary.Size(d.Head); i < n; {
		sub := SubData{}
		err = binary.Read(r, binary.BigEndian, &sub.Head)
		if err != nil {
			return err
		}
		i += int(binary.Size(sub.Head))
		sub.Params = []tlv.TLV{}
		for j := 0; j < int(sub.Head.ParamsCount); j++ {
			param := tlv.TLV{}
			param.FromBinary(r)
			i += int(param.Length())
			sub.Params = append(sub.Params, param)
		}
		d.SubData = append(d.SubData, sub)
	}

	return nil
}
