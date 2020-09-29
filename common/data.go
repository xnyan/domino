package common

import (
	"bufio"
	"os"
)

func LoadKey(keyFile string) []string {
	file, err := os.Open(keyFile)
	if err != nil {
		logger.Fatal(err)
	}
	defer file.Close()

	keyList := make([]string, 0)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		keyList = append(keyList, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		logger.Fatal(err)
	}

	return keyList
}

func GenVal(size int) string {
	val := ""
	for i := 0; i < size; i++ {
		val += "i"
	}
	return val
}
