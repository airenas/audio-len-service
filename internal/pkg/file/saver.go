package file

import (
	"fmt"
	"io"
	"os"
	"path"
)

// Saver saves file
type Saver struct {
	dir string
}

// NewSaver creates temporary file saver dir
func NewSaver(dir string) (*Saver, error) {
	res := Saver{}
	res.dir = dir
	if dir == "" {
		return nil, fmt.Errorf("no temp dir")
	}
	err := os.MkdirAll(dir, 0700)
	return &res, err
}

// Save saves file to temp dir
func (s *Saver) Save(name string, reader io.Reader) (string, error) {
	fn := path.Join(s.dir, name)

	file, err := os.Create(fn)
	if err != nil {
		return "", fmt.Errorf("can't create file %s: %w ", fn, err)
	}
	defer file.Close()

	if _, err := io.Copy(file, reader); err != nil {
		return "", fmt.Errorf("can't write file %s: %w ", fn, err)
	}

	return fn, nil
}
