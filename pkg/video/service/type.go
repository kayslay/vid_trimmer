package service

import "time"

type DownloadStruct struct {
	URL   string
	Start time.Duration
	End   time.Duration
	Type  string
}
