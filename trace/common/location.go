package common

import (
	"bufio"
	"os"
	"strings"
)

type MapParser struct {
	dir map[string]string
}

func NewMapParser() *MapParser {
	return &MapParser{
		dir: make(map[string]string),
	}
}

func (m *MapParser) Load(filePath string) {
	file, err := os.Open(filePath)
	if err != nil {
		logger.Fatal(err)
	}
	defer file.Close()

	// Parses the property file
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		// Skips invalid lines
		if len(line) == 0 || strings.HasPrefix(line, "#") {
			continue
		}
		// Parses a valid line
		w := strings.Fields(line)
		if len(w) != 2 {
			logger.Fatalf("Invalid line %s in directory file %s", line, filePath)
		}
		m.dir[w[0]] = w[1]
	}
}

func (m *MapParser) Get(key string) (string, bool) {
	val, ok := m.dir[key]
	return val, ok
}

func (m *MapParser) GetMap() map[string]string {
	return m.dir
}
