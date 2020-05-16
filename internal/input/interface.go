package input

import "os"

type Interface interface {
	//Fetch takes a path and returns the new path for the file
	Fetch(path string) (string, error)
}

func Remove(path string) error {
	return os.Remove(path)
}
