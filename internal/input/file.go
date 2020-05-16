package input

type file struct {
}

func NewFile() Interface {
	return &file{}
}

func (f file) Fetch(path string) (string, error) {
	return path, nil
}
