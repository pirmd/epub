package epub

import (
	"archive/zip"
	"io"
)

// ReadAtSeeker groups a io.ReaderAt and a io.Seeker.
type ReadAtSeeker interface {
	io.ReaderAt
	io.Seeker
}

func newZipReader(r ReadAtSeeker) (*zip.Reader, error) {
	size, err := getSize(r)
	if err != nil {
		return nil, err
	}

	return zip.NewReader(r, size)
}

func getSize(f io.Seeker) (int64, error) {
	sz, err := f.Seek(0, io.SeekEnd)
	if err != nil {
		return 0, err
	}

	if _, err := f.Seek(0, io.SeekStart); err != nil {
		return 0, err
	}

	return sz, nil
}
