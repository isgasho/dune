package filesystem

import (
	"io"
	"io/ioutil"
	"os"
)

type FS interface {
	Open(name string) (File, error)
	OpenIfExists(name string) (File, error)
	OpenForWrite(name string) (File, error)
	OpenForAppend(name string) (File, error)
	Stat(name string) (os.FileInfo, error)
	Write(name string, data []byte) error
	WritePath(name string, data []byte) error
	Append(name string, data []byte) error
	AppendPath(name string, data []byte) error
	Rename(oldPath, newPath string) error
	RemoveAll(path string) error
	Mkdir(name string) error
	MkdirAll(name string) error
	Chdir(dir string) error
	Getwd() (string, error)
	Abs(name string) (string, error)
	SetHome(name string) error
}

type File interface {
	io.Closer
	io.Reader
	io.ReaderAt
	io.Seeker
	io.WriterAt
	Stat() (os.FileInfo, error)
	Readdir(n int) ([]os.FileInfo, error)
	io.Writer
}

func Exists(fs FS, path string) bool {
	if _, err := fs.Stat(path); err != nil {
		return false
	}
	return true
}

func ReadAll(fs FS, path string) ([]byte, error) {
	f, err := fs.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	b, err := ioutil.ReadAll(f)
	return b, err
}

func ReadDir(fs FS, dirname string) ([]os.FileInfo, error) {
	f, err := fs.Open(dirname)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	list, err := f.Readdir(-1)
	if err != nil {
		return nil, err
	}

	return list, nil
}
