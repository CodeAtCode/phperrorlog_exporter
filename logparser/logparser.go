package logparser

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type Level int

const (
	LevelEError            Level = 1
	LevelEWarning                = 2
	LevelEParse                  = 4
	LevelENotice                 = 8
	LevelECoreError              = 16
	LevelECoreWarning            = 32
	LevelECompileError           = 64
	LevelECompileWarning         = 128
	LevelEUserError              = 256
	LevelEUserWarning            = 512
	LevelEUserNotice             = 1024
	LevelEStrict                 = 2048
	LevelERecoverableError       = 4096
	LevelEDeprecated             = 8192
	LevelEUserDeprecated         = 16384
)

type LogLine struct {
	Level Level
	Time  time.Time
	Error string
	Line  int
}

// LogParser parses logs
type Parser struct {
	filename  string
	locations map[string]*time.Location
}

// NewLogParser well a fresh log parser, if your log file is not too fishy
func NewLogParser(filename string) (p *Parser, err error) {
	return &Parser{
		filename:  filename,
		locations: map[string]*time.Location{},
	}, nil
}

func getFileAndLine(line string) (filename string, lineNumber int, err error) {
	// [02-Dec-2015 14:35:00 Europe/Berlin] PHP Fatal error:  hello in /var/path with blanks/to in depp.php on line 3
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
	if words[numWords-2] != "line" || words[numWords-3] != "on" {
		return
	}

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

func (p *Parser) getDate(line string) (t time.Time, err error) {
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
	start := strings.Index(line, "] PHP ")
	if start == -1 {
		return "", errors.New("start not found")
	}
	rest := line[start+6:] // 6 len("] PHP ")
	end := strings.Index(rest, ":")
	if end == -1 {
		return "", errors.New("end not found")
	}
	return rest[:end], nil
}

// Mon Jan 2 15:04:05 MST 2006 (MST is GMT-0700)
//const timeLayout = "02-Jan-2006 15:04:05 America/Denver"
const timeLayout = "02-Jan-2006 15:04:05"

func extractDateStringFromLine(line string) (dateString string, location string, err error) {
	err = errors.New("date not found")
	// [02-Dec-2015 14:32:23 Europe/Berlin] PHP
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
					dateString += " "
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
				if line[index:(index+5)] != "] PHP" {
					err = errors.New("missing'] PHP '")
					return
				}
				if len(dateString) != 20 && len(location) < 3 {
					err = fmt.Errorf("invalid format %d %d", len(dateString), len(location))
					return
				}
				return
			}
			switch part {
			case 0, 1:
				dateString += string(char)
			case 2:
				location += string(char)
			}
		}
	}
	return
}
