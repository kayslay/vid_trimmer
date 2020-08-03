package input

import (
	"context"
	"os"
)

type Interface interface {
	//Fetch takes a path/url and returns the new tmp path for the file
	Fetch(ctx context.Context, path string) (string, error)
}

func Remove(path string) error {
	return os.Remove(path)
}
