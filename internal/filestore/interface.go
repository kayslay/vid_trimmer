package filestore

import (
	"io"
	"time"
)

const (
	StateNull    = "NULL"
	StatePending = "PENDING"
	StateDone    = "DONE"
	StateError   = "ERROR"
)

type FileStorer interface {
	//Write write the content of the temp file to a store.
	//remove temp file if no longer needed
	Write(key, tempOutputPath string)
	//FileState the state of the file
	FileState(outputPath string) string
	//Get get the file with the key
	Get(key string) (io.ReadCloser, FileStat, error)
	//GeneratePath generates the path in which the file will be saved
	GeneratePath(key string) string
	//ChangeState change the the state
	ChangeState(key, state string)
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
