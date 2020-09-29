package common

import (
	"bufio"
	"os"
	"strconv"
	"strings"
	"time"
)

type Properties struct {
	prop map[string]string
}

func NewProperties() *Properties {
	p := &Properties{
		prop: make(map[string]string, 64),
	}
	return p
}

func (p *Properties) Load(filePath string) {
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
		if idx := strings.Index(line, "="); idx >= 0 {
			key := strings.TrimSpace(line[:idx])
			val := strings.TrimSpace(line[idx+1:])
			p.prop[key] = val
		}
	}
}

func (p *Properties) GetPropMap() map[string]string {
	return p.prop
}

func (p *Properties) Get(k string) (string, bool) {
	v, ok := p.prop[k]
	return v, ok
}

func (p *Properties) GetWithDefault(k, d string) string {
	if v, ok := p.prop[k]; ok {
		return v
	}
	return d
}

//Specific type functions
func (p *Properties) GetStr(k string) string {
	s, ok := p.Get(k)
	if !ok {
		logger.Fatalf("No value for %s", k)
	}
	return s
}

func (p *Properties) GetStrList(k, regex string) []string {
	s, ok := p.Get(k)
	if !ok {
		return make([]string, 0)
	}
	rList := strings.Split(s, regex)
	return rList
}

func (p *Properties) GetInt64WithDefault(k, def string) int64 {
	s := p.GetWithDefault(k, def)
	v, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		logger.Fatalf("Invalid %s = %s", k, s)
	}
	return v
}

func (p *Properties) GetInt64(k string) int64 {
	s, ok := p.Get(k)
	if !ok {
		logger.Fatalf("No value for %s", k)
	}
	v, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		logger.Fatalf("Invalid %s = %s", k, s)
	}
	return v
}

func (p *Properties) GetIntWithDefault(k, def string) int {
	s := p.GetWithDefault(k, def)
	v, err := strconv.Atoi(s)
	if err != nil {
		logger.Fatalf("Invalid %s = %s", k, s)
	}
	return v
}

func (p *Properties) GetInt(k string) int {
	s, ok := p.Get(k)
	if !ok {
		logger.Fatalf("No value for %s", k)
	}
	v, err := strconv.Atoi(s)
	if err != nil {
		logger.Fatalf("Invalid %s = %s", k, s)
	}
	return v
}

func (p *Properties) GetFloat64WithDefault(k, def string) float64 {
	s := p.GetWithDefault(k, def)
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		logger.Fatalf("Invalid %s = %s", k, s)
	}
	return v
}

func (p *Properties) GetBoolWithDefault(k, def string) bool {
	s := p.GetWithDefault(k, def)
	v, err := strconv.ParseBool(s)
	if err != nil {
		logger.Fatalf("Invalid %s = %s", k, s)
	}
	return v
}

func (p *Properties) GetTimeDurationWithDefault(k, def string) time.Duration {
	s := p.GetWithDefault(k, def)
	d, err := time.ParseDuration(s)
	if err != nil {
		logger.Fatalf("Invalid %s = %s", k, s)
	}
	return d
}

func (p *Properties) GetTimeDuration(k string) time.Duration {
	s := p.GetStr(k)
	d, err := time.ParseDuration(s)
	if err != nil {
		logger.Fatalf("Invalid %s = %s", k, s)
	}
	return d
}
