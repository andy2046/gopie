// Package fileutil implements some file utils.
package fileutil

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
)

// ReadDirOperation represents read-directory operation.
type ReadDirOperation struct {
	ext string
}

const (
	// PrivateFileMode represents file read/write mode to owner.
	PrivateFileMode = 0600
	// PrivateDirMode represents file make/remove mode in the directory to owner.
	PrivateDirMode = 0700
)

// WithExt filters file names by extension.
func WithExt(ext string) func(*ReadDirOperation) {
	return func(op *ReadDirOperation) { op.ext = ext }
}

func (op *ReadDirOperation) applyOpts(opts []func(*ReadDirOperation)) {
	for _, opt := range opts {
		opt(op)
	}
}

// ReadDir returns the file names in the provided directory in order.
func ReadDir(d string, opts ...func(*ReadDirOperation)) ([]string, error) {
	var err error
	op := &ReadDirOperation{}
	op.applyOpts(opts)

	dir, err := os.Open(d)
	if err != nil {
		return nil, err
	}
	defer dir.Close()

	names, err := dir.Readdirnames(-1)
	if err != nil {
		return nil, err
	}
	sort.Strings(names)

	if op.ext != "" {
		temp := make([]string, 0)
		for _, v := range names {
			if filepath.Ext(v) == op.ext {
				temp = append(temp, v)
			}
		}
		names = temp
	}
	return names, nil
}

// IsDirWriteable checks if dir is writable by writing and removing a file to dir.
func IsDirWriteable(dir string) error {
	f := filepath.Join(dir, ".touch")
	if err := ioutil.WriteFile(f, []byte(""), PrivateFileMode); err != nil {
		return err
	}
	return os.Remove(f)
}

// TouchDirAll creates directories with 0700 permission if any directory
// does not exist and ensures the provided directory is writable.
func TouchDirAll(dir string) error {
	err := os.MkdirAll(dir, PrivateDirMode)
	if err != nil {
		return err
	}
	return IsDirWriteable(dir)
}

// CreateDirAll wraps TouchDirAll but returns error
// if the deepest directory is not empty.
func CreateDirAll(dir string) error {
	err := TouchDirAll(dir)
	if err == nil {
		var ns []string
		if ns, err = ReadDir(dir); err != nil {
			return err
		}
		if len(ns) != 0 {
			err = fmt.Errorf("expected %q to be empty, got %q", dir, ns)
		}
	}
	return err
}

// Exist returns true if a file or directory exists.
func Exist(name string) bool {
	if _, err := os.Stat(name); os.IsNotExist(err) {
		return false
	}
	return true
}
