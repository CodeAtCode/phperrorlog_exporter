package logparser

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

// Parser parses logs
type Parser struct {
	filename  string
	locations map[string]*time.Location
}

type Stats map[string]int64

type Observation struct {
	Name  string
	Stats Stats
}

func (s Stats) Add(stats Stats) {
	for name, value := range stats {
		s[name] += value
	}
}

func (s Stats) Copy() Stats {
	c := Stats{}
	for k, v := range s {
		c[k] = v
	}
	return c
}

func Read(rd io.Reader) (stats map[string]int64, err error) {
	stats, _, err = read(rd, 0)
	return
}

func Observe(name string, filename string, chanObservation chan Observation, interval time.Duration) {
	currentLine := int64(0)
	stats := Stats{}
	var oldInfo os.FileInfo
	for {
		file, err := os.Open(filename)
		if err == nil {
			newInfo, statErr := file.Stat()
			if statErr == nil {
				if !os.SameFile(newInfo, oldInfo) {
					currentLine = 0
				}
				oldInfo = newInfo
				statsUpdate, numLines, err := read(file, currentLine)
				if err == nil {
					stats.Add(statsUpdate)
					currentLine = numLines
				}
			}
			file.Close()
		} else {
			log.Println("can not open", filename, err)
		}
		chanObservation <- Observation{
			Stats: stats.Copy(),
			Name:  name,
		}
		time.Sleep(interval)
	}
}

func read(rd io.Reader, fromLine int64) (stats map[string]int64, numLines int64, err error) {

	stats = map[string]int64{}
	reader := bufio.NewReader(rd)

	for {
		line, isPrefix, readErr := reader.ReadLine()
		err = readErr
		if err != nil && err != io.EOF {
			return
		}
		numLines++
		if numLines < fromLine {
			continue
		}
		if numLines%1000000 == 0 {
			log.Println("parsed", numLines/1000000, "M lines")
		}

		if err == io.EOF {
			err = nil
			break
		}
		if isPrefix {
			// line is too long
			continue
		}
		_, _, dateErr := extractDateStringFromLine(line)
		if dateErr == nil {
			// new error
			errorName, errorNameErr := extractErrorNameFromLine(string(line))
			if errorNameErr == nil {
				stats[errorName]++
			}
		}
	}
	return
}

// NewLogParser well a fresh log parser, if your log file is not too fishy
func NewLogParser(filename string) (p *Parser, err error) {
	return &Parser{
		filename:  filename,
		locations: map[string]*time.Location{},
	}, nil
}

/*
func reverseBytes(b []byte) (r []byte) {
	r = []byte{}
	for i := len(b) - 1; i > -1; i++ {
		r = append(r, b[i])
	}
	return r
}
*/

func getFileAndLine(lineBytes []byte) (filename string, lineNumber int, err error) {
	// [02-Dec-2015 14:35:00 Europe/Berlin] PHP Fatal error:  hello in /var/path with blanks/to in depp.php on line 3
	lastWords := [][]byte{}
	currentWord := []byte{}
lineByteLoop:
	for i := len(lineBytes) - 1; i > -1; i-- {
		charRune := rune(lineBytes[i])
		switch charRune {
		case ' ':
			lastWords = append(lastWords, currentWord)
			currentWord = []byte{}
		default:
			currentWord = append(currentWord, byte(charRune))
			if len(lastWords) > 3 {
				if bytes.Equal(lastWords[1], []byte("enil")) && bytes.Equal(lastWords[2], []byte("no")) {
					line := string(lineBytes)
					err = errors.New("no found")

					lastIn := strings.LastIndex(line, " in /")
					if lastIn == -1 {
						return
					}

					rest := line[lastIn+4:] // 4 == len(" in ")

					words := strings.Split(rest, " ")
					if len(words) < 4 {
						return
					}
					numWords := len(words)
					lineNumberString := words[numWords-1]
					lineNumber, err = strconv.Atoi(lineNumberString)
					if err != nil {
						err = errors.New("could not parse lineNumber : " + err.Error())
						return
					}
					filename = strings.Join(words[:numWords-3], " ")

					if len(filename) == 0 {
						err = errors.New("filename not found")
						return
					}
					err = nil

					return

				}
				break lineByteLoop
			}
		}
	}
	err = errors.New("file and line not found")
	return
}

func (p *Parser) getDate(line []byte) (t time.Time, err error) {
	dateString, loc, err := extractDateStringFromLine(line)
	if err != nil {
		return
	}
	location, ok := p.locations[loc]
	if !ok {
		location, err = time.LoadLocation(loc)
		if err != nil {
			return
		}
		p.locations[loc] = location
	}
	return getDate(dateString, location)
}

func getDate(dateString string, location *time.Location) (date time.Time, err error) {
	return time.ParseInLocation(timeLayout, dateString, location)
}

func extractErrorNameFromLine(line string) (errorName string, err error) {
	start := strings.Index(line, "] ")
	if start == -1 {
		return "", errors.New("start not found")
	}
	rest := line[start+6:] // 6 len("] PHP")
	end := strings.Index(rest, ":")
	if end == -1 {
		return "", errors.New("end not found")
	}
	errorName = rest[:end]
	switch errorName {
	case "Deprecated":
		errorName = "deprecated"
	case "Parse error":
		errorName = "parse"
	case "Fatal error", "Catchable fatal error":
		errorName = "fatal"
	case "Notice":
		errorName = "notice"
	case "Warning":
		errorName = "warning"
	default:
		errorName = ""
		err = errors.New("unknown error name: " + errorName)
	}
	return errorName, err
}

// Mon Jan 2 15:04:05 MST 2006 (MST is GMT-0700)
//const timeLayout = "02-Jan-2006 15:04:05 America/Denver"
const timeLayout = "02-Jan-2006 15:04:05"

func extractDateStringFromLine(line []byte) (dateString string, location string, err error) {
	err = errors.New("date not found")
	// [02-Dec-2015 14:32:23 Europe/Berlin] PHP
	dateBytes := []byte{}
	if len(line) > 30 {
		part := 0
		for index, char := range line {
			//log.Println(index, char, string(char))
			switch index {
			case 0:
				if char != '[' {
					err = errors.New("invalid beginning of line")
					return
				}
				continue
			case 100:
				err = errors.New("giving up")
				return
			}
			switch char {
			case ' ':
				if part == 0 {
					dateBytes = append(dateBytes, ' ')
				}
				part++
				if part > 2 {
					return
				}
				continue
			case ']':
				err = nil
				minLength := index + 4 // 4 == len(" PHP")
				if len(line) < minLength {
					err = errors.New("missing ] PHP")
					return
				}
				if string(line[index:(index+5)]) != "] PHP" {
					err = errors.New("missing'] PHP '")
					return
				}
				if len(dateString) != 20 && len(location) < 3 {
					err = fmt.Errorf("invalid format %d %d", len(dateString), len(location))
					return
				}
				dateString = string(dateBytes)
				return
			}
			switch part {
			case 0, 1:
				dateBytes = append(dateBytes, char)
			case 2:
				location += string(char)
			}
		}
	}
	return
}
