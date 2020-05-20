package input

import "context"

type file struct {
}

func NewFile() Interface {
	return &file{}
}

func (f file) Fetch(ctx context.Context, path string) (string, error) {
	return path, nil
}
