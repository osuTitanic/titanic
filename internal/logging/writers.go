package logging

import (
	"io"
	"os"
)

func GetConsoleWriter() io.Writer {
	return os.Stdout
}

func GetFileWriter(path string) (io.Writer, error) {
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}
	return file, nil
}
