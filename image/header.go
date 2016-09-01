package image

import (
	"encoding/binary"
	"io"
)

// Header is a section header in the executable file.
type Header struct {
	Type uint8
	Flag uint8
	_    [2]uint8
	Addr uint32
	Size uint32

	offset uint32
}

const sectionLen = 16

// ReadHeader reads a new section header from the reader.
func ReadHeader(r io.Reader) (*Header, error) {
	h := new(Header)
	_, err := h.ReadFrom(r)
	if err != nil {
		return nil, err
	}
	return h, nil
}

// WriteTo writes the section to a writer.
func (h *Header) WriteTo(w io.Writer) (int64, error) {
	buf := make([]byte, sectionLen)
	enc := binary.LittleEndian
	buf[0] = h.Type
	buf[1] = h.Flag
	enc.PutUint32(buf[4:8], h.Addr)
	enc.PutUint32(buf[8:12], h.Size)
	enc.PutUint32(buf[12:16], h.offset)

	n, err := w.Write(buf)
	return int64(n), err
}

// ReadFrom read in the section from a reader.
func (h *Header) ReadFrom(r io.Reader) (int64, error) {
	buf := make([]byte, sectionLen)
	n, err := io.ReadFull(r, buf)
	if err != nil {
		return int64(n), err
	}

	enc := binary.LittleEndian
	h.Type = buf[0]
	h.Flag = buf[1]
	h.Addr = enc.Uint32(buf[4:8])
	h.Size = enc.Uint32(buf[8:12])
	h.offset = enc.Uint32(buf[12:16])
	return int64(n), nil
}
