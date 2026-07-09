package storage

import (
	"errors"
	"io"
	"sync"
)

type bufferedReader struct {
	source ReaderAtCloser
	size   int64

	mu          sync.Mutex
	buffer      []byte
	bufferSize  int
	bufferStart int64
	bufferLen   int
}

func newBufferedReader(source ReaderAtCloser, size int64, bufferSize int) ReaderAtCloser {
	if bufferSize <= 0 {
		return source
	}
	return &bufferedReader{
		source:     source,
		size:       size,
		bufferSize: bufferSize,
	}
}

func (r *bufferedReader) ReadAt(p []byte, offset int64) (int, error) {
	if len(p) == 0 {
		return 0, nil
	}
	if offset < 0 {
		return 0, errors.New("negative read offset")
	}
	if offset >= r.size {
		return 0, io.EOF
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	wanted := len(p)
	remaining := r.size - offset
	if int64(wanted) > remaining {
		wanted = int(remaining)
	}

	copied := 0
	for copied < wanted {
		currentOff := offset + int64(copied)

		n := r.copyFromBuffer(p[copied:wanted], currentOff)
		if n > 0 {
			copied += n
			continue
		}

		n, err := r.fillBuffer(currentOff, wanted-copied)
		if err != nil && n == 0 {
			if copied > 0 {
				return copied, err
			}
			return 0, err
		}

		if n == 0 {
			if copied > 0 {
				return copied, io.EOF
			}
			return 0, io.EOF
		}
	}

	if copied < len(p) {
		return copied, io.EOF
	}
	return copied, nil
}

func (r *bufferedReader) Close() error {
	return r.source.Close()
}

func (r *bufferedReader) copyFromBuffer(p []byte, offset int64) int {
	bufferOffset := offset - r.bufferStart
	if bufferOffset < 0 || bufferOffset >= int64(r.bufferLen) {
		return 0
	}

	start := int(bufferOffset)
	return copy(p, r.buffer[start:r.bufferLen])
}

func (r *bufferedReader) fillBuffer(off int64, need int) (int, error) {
	if off >= r.size {
		r.bufferStart = off
		r.bufferLen = 0
		r.buffer = r.buffer[:0]
		return 0, io.EOF
	}

	readSize := max(need, r.bufferSize)
	remaining := r.size - off
	if int64(readSize) > remaining {
		readSize = int(remaining)
	}

	if cap(r.buffer) < readSize {
		r.buffer = make([]byte, readSize)
	} else {
		r.buffer = r.buffer[:readSize]
	}

	n, err := r.source.ReadAt(r.buffer, off)
	r.bufferStart = off
	r.bufferLen = n
	r.buffer = r.buffer[:n]

	if err != nil && n > 0 {
		return n, nil
	}
	return n, err
}
