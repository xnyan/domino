package common

import (
	"os"
)

type Writer interface {
	Write(text string)
}

type FileWriter struct {
	filePath string
	io       *os.File
}

func NewFileWriter(filePath string) *FileWriter {
	w := &FileWriter{
		filePath: filePath,
	}

	var err error
	w.io, err = os.Create(filePath)
	if err != nil {
		logger.Errorf("Error: %v", err)
		logger.Fatalf("Fails to create file %s", filePath)
	}

	return w
}

func (w *FileWriter) Write(text string) {
	_, err := w.io.WriteString(text)
	if err != nil {
		logger.Errorf("Error: %v", err)
		logger.Fatalf("Fails to write file %s (text = %s)", w.filePath, text)
	}
}

func (w *FileWriter) Sync() {
	err := w.io.Sync()
	if err != nil {
		logger.Errorf("Error: %v", err)
		logger.Fatalf("Fails to sync file %s", w.filePath)
	}
}

func (w *FileWriter) Close() {
	err := w.io.Close()
	if err != nil {
		logger.Errorf("Error: %v", err)
		logger.Fatalf("Fails to close file %s", w.filePath)
	}
}
