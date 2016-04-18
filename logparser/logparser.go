package logparser

import (
	"errors"
	"fmt"
	"log"
)

// Mon Jan 2 15:04:05 MST 2006 (MST is GMT-0700)
//const timeLayout = "02-Jan-2006 15:04:05 America/Denver"
const timeLayout = "02-Jan-2006 15:04:05"

func extractDateStringFromLine(line string) (dateString string, location string, err error) {
	err = errors.New("date not found")
	// [02-Dec-2015 14:32:23 Europe/Berlin] PHP
	if len(line) > 30 {
		part := 0
		for index, char := range line {
			log.Println(index, char, string(char))
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
				part++
				continue
			case ']':
				err = nil
				if len(dateString) != 20 && len(location) < 3 {
					err = fmt.Errorf("invalid format %d %d", len(dateString), len(location))
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
