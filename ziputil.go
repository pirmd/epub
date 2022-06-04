package epub

import (
	"archive/zip"
	"io"
	"os"
)

// ReadAtSeeker groups a io.ReaderAt and a io.Seeker.
type ReadAtSeeker interface {
	io.ReaderAt
	io.Seeker
}

func openFromZip(r ReadAtSeeker, path string) (io.ReadCloser, error) {
	zr, err := newZipReader(r)
	if err != nil {
		return nil, err
	}

	for _, f := range zr.File {
		if f.Name == path {
			return f.Open()
		}
	}
	return nil, os.ErrNotExist
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
