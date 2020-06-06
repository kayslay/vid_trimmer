package filestore

import (
	"io"
	"os"
	"path"
)

type tempFile struct{}

func NewTempFile() FileStorer {
	return &tempFile{}
}

func (tf tempFile) Write(key string, r io.ReadCloser) {

	//	this does nothing since the file is already in the temp folder
}

func (tf tempFile) Exists(key string) bool {
	if _, err := os.Open(tf.GeneratePath(key)); err != nil {
		return false
	}
	return true
}

func (tf tempFile) Get(key string) (io.ReadCloser, FileStat, error) {
	file, err := os.Open(tf.GeneratePath(key))
	if err != nil {
		return nil, FileStat{}, err
	}

	fileInfo, _ := file.Stat()

	return file, FileStat{
		Size: fileInfo.Size(),
	}, nil
}

func (tf tempFile) GeneratePath(key string) string {
	return path.Join(os.TempDir(), key)
}
