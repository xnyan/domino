package main

import (
	//"bufio"
	//"strings"
	"flag"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/op/go-logging"

	"domino/trace/node"
)

var logger = logging.MustGetLogger("Parser")

var filePath string

// mode options
// n : all in nanoseconds (default)
// t : all in time format
// o : order time in ns, lat in ns, time offset in ns
// h : order time format, lat in ns, time offset in ns
var mode string

// order time options
// r : receiving time (default)
// s : sending time
var orderTime string
var timeFormat string
var outputFile string

// TimeZone.Year.Month.Day.Hour.Minute.Second.Nanoseconds
const defaultFormat = "UTC.2006.01.02.15.04.05.000000000"

func main() {
	// Parses command line args
	parseArgs()

	file, err := os.Open(filePath)
	checkErr(err)
	defer file.Close()

	log, err := os.Create(outputFile)
	checkErr(err)
	defer log.Close()

	if mode == "o" || mode == "h" {
		if orderTime == "r" {
			log.WriteString("#RevTime Latency ClockOffset\n")
		} else if orderTime == "s" {
			log.WriteString("#SendTime Latency ClockOffset\n")
		} else {
			logger.Fatalf("Invalid order time type = %s", orderTime)
		}
	} else {
		log.WriteString("#SendTime RevTime ServerClockTime\n")
	}
	b := make([]byte, 24, 24)
	for n, err := file.Read(b); n != 0 && err != io.EOF; n, err = file.Read(b) {
		l := node.ByteToLatInfo(b)
		s := formatL(mode, l)
		log.WriteString(s)
	}
}

func formatL(m string, l *node.L) string {
	switch m {
	case "n":
		return fmt.Sprintf("%d %d %d\n", l.S, l.E, l.C)
	case "t":
		return fmt.Sprintf("%s %s %s\n",
			time.Unix(0, l.S).Format(timeFormat),
			time.Unix(0, l.E).Format(timeFormat),
			time.Unix(0, l.C).Format(timeFormat))
	case "o":
		if orderTime == "r" {
			return fmt.Sprintf("%d %d %d\n", l.E, l.E-l.S, l.C-l.S)
		} else if orderTime == "s" {
			return fmt.Sprintf("%d %d %d\n", l.S, l.E-l.S, l.C-l.S)
		}
	case "h":
		if orderTime == "r" {
			return fmt.Sprintf("%s %d %d\n",
				time.Unix(0, l.E).Format(timeFormat), l.E-l.S, l.C-l.S)
		} else if orderTime == "s" {
			return fmt.Sprintf("%s %d %d\n",
				time.Unix(0, l.S).Format(timeFormat), l.E-l.S, l.C-l.S)
		}
	default:
		logger.Fatalf("Invalid mode %s", m)
	}
	return ""
}

func parseArgs() {
	flag.StringVar(&filePath, "f", "", "probing information log file")

	modeHelp := "parsing mode options:\n" +
		"n : all in nanoseconds\n" +
		"t : all in time format\n" +
		"o : order time in ns, lat in ns, time offset in ns\n" +
		"h : order time format, lat in ns, time offset in ns\n"
	flag.StringVar(&mode, "m", "n", modeHelp)

	orderHelp := "order time:\n" +
		"r : receiving time\n" +
		"s : sending time\n"
	flag.StringVar(&orderTime, "p", "r", orderHelp)

	flag.StringVar(&timeFormat, "t", defaultFormat, "time format, e.g., UTC.YYYY.MM.DD.HH.MM.SS.ns\n")
	flag.StringVar(&outputFile, "o", "", "output file (default: inputFile.txt)")
	flag.Parse()

	if filePath == "" {
		flag.Usage()
		logger.Fatalf("Missing log file.")
	}
	if outputFile == "" {
		outputFile = filePath + ".txt"
	}
}

func checkErr(err error) {
	if err != nil {
		logger.Fatal(err)
	}
}
