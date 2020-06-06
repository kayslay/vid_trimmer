package filestore

import (
	"io"
	"time"
)

type FileStorer interface {
	Write(path string, r io.ReadCloser)
	Exists(path string) bool
	Get(path string) (io.ReadCloser, FileStat, error)
	GeneratePath(key string) string
}

type FileStat struct {
	Size int64
}

type DownloadStruct struct {
	URL   string
	Start time.Duration
	End   time.Duration
	Type  string
}
