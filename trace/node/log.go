package node

import (
	"os"
)

type Log interface {
	Write(b []byte)
	Flush()
	Close()
}

func NewLog(fileName string) Log {
	return newByteLog(fileName)
}

type ByteLog struct {
	fileName string
	log      *os.File
}

func newByteLog(fileName string) Log {
	l := &ByteLog{
		fileName: fileName,
	}
	if log, err := os.Create(fileName); err == nil {
		l.log = log
	} else {
		logger.Fatalf("Fails to create file %s error: %v", fileName, err)
	}

	return l
}

func (l *ByteLog) Write(b []byte) {
	n, err := l.log.Write(b)
	if err != nil {
		logger.Fatalf("Fails to write bytes %v", b)
	}
	if n != len(b) {
		logger.Fatalf("Expected to write %d bytes but only %d", len(b), n)
	}
	//logger.Debugf("Wrote %d bytes to file %s", n, l.fileName)

	//_, err = l.log.WriteString("\n")
	//if err != nil {
	//	logger.Fatalf("Faist to write a line space")
	//}
}

func (l *ByteLog) Flush() {
	l.log.Sync()
}

func (l *ByteLog) Close() {
	l.log.Close()
}
