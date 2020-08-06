package filestore

import (
	"fmt"
	"github.com/patrickmn/go-cache"
	"io"
	"os"
	"path"
	"time"
)

const (
	cacheDuration = 30 * time.Minute
)

type tempFile struct {
	cache *cache.Cache
}

func NewTempFile() FileStorer {
	return &tempFile{cache: cache.New(cacheDuration, 10*time.Minute)}
}

func (tf tempFile) Write(key, tempOutputPath string) {

	//	this does nothing since the file is already in the temp folder
	tf.cache.Set(key, StateDone, cacheDuration)
}

func (tf tempFile) FileState(key string) string {
	s, exists := tf.cache.Get(key)
	if !exists {
		return StateNull
	}

	state := fmt.Sprint(s)
	if state != StateDone {
		return state
	}

	if _, err := os.Open(tf.GeneratePath(key)); err != nil {
		return StateNull
	}

	return fmt.Sprint(s)

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

func (tf tempFile) ChangeState(key, state string) {
	tf.cache.Set(key, state, cacheDuration)
}
