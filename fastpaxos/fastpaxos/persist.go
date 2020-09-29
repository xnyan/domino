package fastpaxos

import (
	//"github.com/op/go-logging"
	"os"
)

type Logger struct {
	logFilePath string
	log         *os.File
}

func NewLogger(filePath string) *Logger {
	l := &Logger{
		logFilePath: filePath,
	}

	var err error
	l.log, err = os.Create(filePath)
	if err != nil {
		logger.Fatalf("Fails to create file %s", filePath)
	}

	return l
}

func (l *Logger) Write(text string) {
	_, err := l.log.WriteString(text)
	if err != nil {
		logger.Fatalf("Fails to write to file %s with text = %s", l.logFilePath, text)
	}
}

func (l *Logger) Sync() {
	err := l.log.Sync()
	if err != nil {
		logger.Fatalf("Fails to sync file %s", l.logFilePath)
	}
}

func (l *Logger) Close() {
	err := l.log.Close()
	if err != nil {
		logger.Fatalf("Error in closing file %s", l.logFilePath)
	}
}
