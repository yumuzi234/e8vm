package image

import (
	"errors"
	"io"
	"math"
	"os"
)

// Section is a block section.
type Section struct {
	*Header
	Bytes []byte
}

func readHeaders(r io.Reader) ([]*Section, error) {
	var ret []*Section
	for {
		h, err := ReadHeader(r)
		if err != nil {
			return nil, err
		}
		if h.Type == None {
			break
		}
		ret = append(ret, &Section{Header: h})
	}
	return ret, nil
}

// Read reads in an executable file
func Read(r io.ReadSeeker) ([]*Section, error) {
	ret, err := readHeaders(r)
	if err != nil {
		return nil, err
	}

	for _, s := range ret {
		if s.Type == Zeros {
			continue
		}

		s.Bytes = make([]byte, s.Size)
		r.Seek(int64(s.offset), 0)
		if _, err := io.ReadFull(r, s.Bytes); err != nil {
			return nil, err
		}
	}

	return ret, nil
}

// Open opens an executable file from the file system.
func Open(path string) ([]*Section, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return Read(f)
}

// Write writes the sections into a writer.
func Write(w io.Writer, sections []*Section) error {
	offset := int64(sectionLen * (len(sections) + 1))
	if offset > math.MaxUint32 {
		return errors.New("too many headers")
	}

	for _, s := range sections {
		if int64(len(s.Bytes)) > math.MaxUint32 {
			return errors.New("too many bytes in a section")
		}
		if offset > math.MaxUint32-int64(len(s.Bytes)) {
			return errors.New("too many bytes in total")
		}

		var h = *s.Header

		if s.Bytes != nil {
			h.Size = uint32(len(s.Bytes))
			h.offset = uint32(offset)
		} else {
			h.offset = 0
		}

		if _, err := h.WriteTo(w); err != nil {
			return err
		}

		offset += int64(len(s.Bytes))
	}

	// write an empty header for terminating the headers.
	var empty Header
	if _, err := empty.WriteTo(w); err != nil {
		return err
	}

	// now the contents.
	for _, s := range sections {
		if s.Bytes == nil {
			continue
		}
		if _, err := w.Write(s.Bytes); err != nil {
			return err
		}
	}

	return nil
}

// Create creates an executable file.
func Create(path string, sections []*Section) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return Write(f, sections)
}
