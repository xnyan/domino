package common

import (
	"bufio"
	"os"
	"strings"
)

// Loads a property file
func LoadProperties(filePath string) map[string]string {
	file, err := os.Open(filePath)
	if err != nil {
		logger.Fatal(err)
	}
	defer file.Close()

	config := make(map[string]string)

	// Parses the property file
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		// Skips invalid lines
		if len(line) == 0 || strings.HasPrefix(line, "#") {
			continue
		}
		// Parses a valid line
		if idx := strings.Index(line, "="); idx >= 0 {
			key := strings.TrimSpace(line[:idx])
			val := strings.TrimSpace(line[idx+1:])
			config[key] = val
		}
	}
	return config
}

func GetProperty(config map[string]string, key string, defVal string) string {
	if val, exists := config[key]; exists {
		return val
	}
	return defVal
}

func GetValue(config map[string]string, key string) string {
	return config[key]
}
