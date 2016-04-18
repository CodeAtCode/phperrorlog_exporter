package logparser

import (
	"testing"
	"time"
)

const (
	lineNotice  = `[02-Dec-2015 14:32:23 Europe/Berlin] PHP Notice:  Foomo\Cache\Persistence\Queryable\PDOPersistor::connect PDO connected at attempt: 1 in /var/www/schild/modules/Foomo/lib/Foomo/Cache/Persistence/Queryable/PDOPersistor.php on line 386`
	lineWarning = `[02-Dec-2015 14:35:00 Europe/Berlin] PHP Warning:  mkdir(): File exists in /var/www/schild/modules/Foomo/lib/Foomo/Setup.php on line 160`
)

func TestDate(t *testing.T) {
	location, _ := time.LoadLocation("Europe/Berlin")
	date, err := time.ParseInLocation(timeLayout, "02-Dec-2015 14:32:23", location)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(date)
}

func TestExtractDateStringFromLine(t *testing.T) {
	te, z, e := extractDateStringFromLine(lineNotice)
	t.Log("te", te, "z", z, "e", e)

}
