package logparser

import "testing"

const (
	lineCrap     = `böa`
	lineLongCrap = `[¢jkjkljas fjksajfksajfklas jfklsjfklsaf jfklajkljafkljfkl asjklasjfklas jfkalsjfklsf jasklfj askl fjasklfjsklaf jaklfjkflas djklfasj fklasjdfkl j fklas jfakslöböa`
	lineNotice   = `[02-Dec-2015 14:32:23 Europe/Berlin] PHP Notice:  Foomo\Cache\Persistence\Queryable\PDOPersistor::connect PDO connected at attempt: 1 in /var/www/schild/modules/Foomo/lib/Foomo/Cache/Persistence/Queryable/PDOPersistor.php on line 386`
	lineWarning  = `[02-Dec-2015 14:35:00 Europe/Berlin] PHP Warning:  mkdir(): File exists in /var/www/schild/modules/Foomo/lib/Foomo/Setup.php on line 160`
	lineAlmost   = `[02-Dec-2015 14:35:00 Europe/Berlin] PHp hjkh`

	lineMultiUnknown = `[02-Dec-2015 14:35:00 Europe/Berlin] PHP Warning:  PHP Startup: Unable to load dynamic library '/usr/local/opt/php56-yaml/yaml.so' - dlopen(/usr/local/opt/php56-yaml/yaml.so, 9): Symbol not found: _basic_globals
  Referenced from: /usr/local/opt/php56-yaml/yaml.so
  Expected in: flat namespace
 in /usr/local/opt/php56-yaml/yaml.so in Unknown on line 0`
	lineInBlanks = `[02-Dec-2015 14:35:00 Europe/Berlin] PHP Fatal error:  hello in /Users/jan/go/src/github.com/foomo/phperrorlog_exporter/sepp in depp.php on line 3`
)

func poe(err error) {
	if err != nil {
		panic(err)
	}
}

func pone(err error) {
	if err == nil {
		panic("no err")
	}
}

func TestDate(t *testing.T) {
	p, err := NewLogParser("foo")
	poe(err)
	lt, e := p.getDate(lineNotice)
	t.Log(lt, e)
	poe(e)
	lt, e = p.getDate(lineWarning)
	t.Log(lt, e)
	poe(e)
	lt, e = p.getDate(lineCrap)
	t.Log(lt, e)
	pone(e)
	lt, e = p.getDate(lineLongCrap)
	t.Log(lt, e)
	pone(e)
	lt, e = p.getDate(lineAlmost)
	t.Log(lt, e)
	pone(e)

}

func TestGetFileAndLine(t *testing.T) {
	type expectedFileAndLine struct {
		line int
		file string
		err  bool
	}
	expectations := map[string]expectedFileAndLine{
		lineNotice: expectedFileAndLine{
			line: 386,
			file: "/var/www/schild/modules/Foomo/lib/Foomo/Cache/Persistence/Queryable/PDOPersistor.php",
			err:  false,
		},
		lineCrap: expectedFileAndLine{
			err: true,
		},
	}
	for line, expectation := range expectations {
		actualFile, actualLine, err := getFileAndLine(line)
		t.Log("actualFile, actualLine, err", actualFile, actualLine, err, "from", line)
		if expectation.err {
			if err == nil {
				t.Fatal("where is my error")
			}
			continue
		}
		if actualFile != expectation.file {
			t.Fatal("file missmatch", actualFile, "!=", expectation.file)
		}
		if actualLine != expectation.line {
			t.Fatal("line missmatch", actualLine, "!=", expectation.line)
		}
	}
}

func TestExtractErrorNameFromLine(t *testing.T) {
	expectations := map[string]string{
		lineNotice:   "Notice",
		lineWarning:  "Warning",
		lineAlmost:   "",
		lineCrap:     "",
		lineLongCrap: "",
	}
	for line, expectedName := range expectations {
		errorName, err := extractErrorNameFromLine(line)
		t.Log("expected", expectedName, "got", errorName, "in", line)
		if expectedName == "" {
			if err == nil || errorName != "" {
				t.Fatal("there should have been an error", errorName)
			}
		} else {
			if err != nil {
				t.Fatal("unexpected error", err, "in", line)
			}
			if expectedName != errorName {
				t.Fatal("unexpected", expectedName, "!=", errorName)
			}
		}
	}
}

func TestExtractDateStringFromLine(t *testing.T) {
	te, z, e := extractDateStringFromLine(lineNotice)
	t.Log("te", te, "z", z, "e", e)
}
