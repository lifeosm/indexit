package output

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type Writer struct {
	file   *os.File
	writer *bufio.Writer
}

func New(path string, stdout io.Writer) (*Writer, error) {
	if path == "" || path == "-" {
		return &Writer{writer: bufio.NewWriter(stdout)}, nil
	}
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return nil, fmt.Errorf("create output directory: %w", err)
	}
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, fmt.Errorf("open output: %w", err)
	}
	return &Writer{file: file, writer: bufio.NewWriter(file)}, nil
}

func (w *Writer) Write(v any) error {
	line, err := json.Marshal(v)
	if err != nil {
		return fmt.Errorf("encode jsonl record: %w", err)
	}
	line = append(line, '\n')
	if _, err := w.writer.Write(line); err != nil {
		return fmt.Errorf("write jsonl record: %w", err)
	}
	return w.writer.Flush()
}

func (w *Writer) Close() error {
	var err error
	if w.writer != nil {
		err = w.writer.Flush()
	}
	if w.file != nil {
		if closeErr := w.file.Close(); err == nil {
			err = closeErr
		}
	}
	return err
}
