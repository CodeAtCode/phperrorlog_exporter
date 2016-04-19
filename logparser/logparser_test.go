package logparser

import (
	"bytes"
	"testing"
)

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

func TestRead(t *testing.T) {
	const l = `[02-Dec-2015 14:58:30 Europe/Berlin] PHP Warning:  socket_connect(): unable to connect [110]: Connection timed out in /var/www/schild/modules/Foomo.ContentServer/lib/Foomo/ContentServer/Client.php on line 89
[02-Dec-2015 15:00:37 Europe/Berlin] PHP Warning:  socket_connect(): unable to connect [110]: Connection timed out in /var/www/schild/modules/Foomo.ContentServer/lib/Foomo/ContentServer/Client.php on line 89
[02-Dec-2015 15:36:59 Europe/Berlin] PHP Notice:  Foomo\Cache\Persistence\Queryable\PDOPersistor::connect PDO connected at attempt: 1 in /var/www/schild/modules/Foomo/lib/Foomo/Cache/Persistence/Queryable/PDOPersistor.php on line 386
[02-Dec-2015 15:36:59 Europe/Berlin] PHP Notice:  
Foomo\AutoLoader::buildClassMapbuilding a new classmap from :
  /var/www/schild/modules/Foomo/tests
  /var/www/schild/modules/Foomo/lib
  /var/www/schild/modules/Foomo.Backbone/lib
  /var/www/schild/modules/Foomo.Backbone.Demo/tests
  /var/www/schild/modules/Foomo.Backbone.Demo/lib
  /var/www/schild/modules/Foomo.Bundle/tests
  /var/www/schild/modules/Foomo.Bundle/lib
  /var/www/schild/modules/Foomo.CSV/tests
  /var/www/schild/modules/Foomo.CSV/lib
  /var/www/schild/modules/Foomo.ContentServer/lib
  /var/www/schild/modules/Foomo.Docs/lib
  /var/www/schild/modules/Foomo.Go/tests
  /var/www/schild/modules/Foomo.Go/lib
  /var/www/schild/modules/Foomo.JS/lib
  /var/www/schild/modules/Foomo.Jasmine/lib
  /var/www/schild/modules/Foomo.Less/lib
  /var/www/schild/modules/Foomo.Media/tests
  /var/www/schild/modules/Foomo.Media/lib
  /var/www/schild/modules/Foomo.Monolog/tests
  /var/www/schild/modules/Foomo.Monolog/lib
  /var/www/schild/modules/Foomo.Sass/tests
  /var/www/schild/modules/Foomo.Sass/lib
  /var/www/schi in /var/www/schild/modules/Foomo/lib/Foomo/AutoLoader.php on line 296
[02-Dec-2015 15:37:01 Europe/Berlin] PHP Notice:  Foomo\AutoLoader::buildClassMap putting classmap to cache with 1394 classes in /var/www/schild/modules/Foomo/lib/Foomo/AutoLoader.php on line 327
[02-Dec-2015 15:37:01 Europe/Berlin] PHP Notice:  Foomo\Cache\Persistence\Queryable\PDOPersistorFoomo\Cache\Persistence\Queryable\PDOPersistor::storeResourceNameSQLSTATE[42S02]: Base table or view not found: 1146 Table 'schildFoomoCacheTest.$$$CACHED_RESOURCE_NAMES$$$' doesn't exist in /var/www/schild/modules/Foomo/lib/Foomo/Cache/Persistence/Queryable/PDOPersistor.php on line 538
[02-Dec-2015 15:37:01 Europe/Berlin] PHP Fatal error:  Class 'MongoCursor' not found in /var/www/schild/modules/Foomo.SimpleData.MongoDB/lib/Foomo/SimpleData/MongoDB/Cursor.php on line 29
[02-Dec-2015 15:37:40 Europe/Berlin] PHP Fatal error:  Class 'MongoCursor' not found in /var/www/schild/modules/Foomo.SimpleData.MongoDB/lib/Foomo/SimpleData/MongoDB/Cursor.php on line 29
// and here with a stack trace
[10-Feb-2016 00:39:01 UTC] PHP Fatal error:  Call to a member function isAvailable() on a non-object in /home/butkus/public_html/app/code/core/Mage/Adminhtml/Block/Widget/Grid.php on line 572
[10-Feb-2016 00:39:01 UTC] PHP Stack trace:
[10-Feb-2016 00:39:01 UTC] PHP   1. {main}() /home/butkus/public_html/index.php:0
[10-Feb-2016 00:39:01 UTC] PHP   2. Mage::run($code = *uninitialized*, $type = *uninitialized*, $options = *uninitialized*) /home/butkus/public_html/index.php:86
[10-Feb-2016 00:39:01 UTC] PHP   3. Mage_Core_Model_App->run($params = *uninitialized*) /home/butkus/public_html/app/Mage.php:683
[10-Feb-2016 00:39:01 UTC] PHP   4. Mage_Core_Controller_Varien_Front->dispatch() /home/butkus/public_html/app/code/core/Mage/Core/Model/App.php:354
[10-Feb-2016 00:39:01 UTC] PHP   5. Mage_Core_Controller_Varien_Router_Standard->match($request = *uninitialized*) /home/butkus/public_html/app/code/core/Mage/Core/Controller/Varien/Front.php:176
[10-Feb-2016 00:39:01 UTC] PHP   6. Mage_Core_Controller_Varien_Action->dispatch($action = *uninitialized*) /home/butkus/public_html/app/code/core/Mage/Core/Controller/Varien/Router/Standard.php:254
[10-Feb-2016 00:39:01 UTC] PHP   7. Mage_Adminhtml_DashboardController->indexAction() /home/butkus/public_html/app/code/core/Mage/Core/Controller/Varien/Action.php:419
[10-Feb-2016 00:39:01 UTC] PHP   8. Mage_Adminhtml_Controller_Action->loadLayout($ids = *uninitialized*, $generateBlocks = *uninitialized*, $generateXml = *uninitialized*) /home/butkus/public_html/app/code/core/Mage/Adminhtml/controllers/DashboardController.php:40
[10-Feb-2016 00:39:01 UTC] PHP   9. Mage_Core_Controller_Varien_Action->loadLayout($handles = *uninitialized*, $generateBlocks = *uninitialized*, $generateXml = *uninitialized*) /home/butkus/public_html/app/code/core/Mage/Adminhtml/Controller/Action.php:275
[10-Feb-2016 00:39:01 UTC] PHP  10. Mage_Core_Controller_Varien_Action->generateLayoutBlocks() /home/butkus/public_html/app/code/core/Mage/Core/Controller/Varien/Action.php:269
[10-Feb-2016 00:39:01 UTC] PHP  11. Mage_Core_Model_Layout->generateBlocks($parent = *uninitialized*) /home/butkus/public_html/app/code/core/Mage/Core/Controller/Varien/Action.php:344
[10-Feb-2016 00:39:01 UTC] PHP  12. Mage_Core_Model_Layout->generateBlocks($parent = *uninitialized*) /home/butkus/public_html/app/code/core/Mage/Core/Model/Layout.php:210
[10-Feb-2016 00:39:01 UTC] PHP  13. Mage_Core_Model_Layout->_generateBlock($node = *uninitialized*, $parent = *uninitialized*) /home/butkus/public_html/app/code/core/Mage/Core/Model/Layout.php:205
[10-Feb-2016 00:39:01 UTC] PHP  14. Mage_Core_Model_Layout->addBlock($block = *uninitialized*, $blockName = *uninitialized*) /home/butkus/public_html/app/code/core/Mage/Core/Model/Layout.php:239
[10-Feb-2016 00:39:01 UTC] PHP  15. Mage_Core_Model_Layout->createBlock($type = *uninitialized*, $name = *uninitialized*, $attributes = *uninitialized*) /home/butkus/public_html/app/code/core/Mage/Core/Model/Layout.php:472
[10-Feb-2016 00:39:01 UTC] PHP  16. Mage_Core_Block_Abstract->setLayout($layout = *uninitialized*) /home/butkus/public_html/app/code/core/Mage/Core/Model/Layout.php:456
[10-Feb-2016 00:39:01 UTC] PHP  17. Mage_Adminhtml_Block_Dashboard->_prepareLayout() /home/butkus/public_html/app/code/core/Mage/Core/Block/Abstract.php:238
[10-Feb-2016 00:39:01 UTC] PHP  18. Mage_Core_Model_Layout->createBlock($type = *uninitialized*, $name = *uninitialized*, $attributes = *uninitialized*) /home/butkus/public_html/app/code/core/Mage/Adminhtml/Block/Dashboard.php:75
[10-Feb-2016 00:39:01 UTC] PHP  19. Mage_Core_Block_Abstract->setLayout($layout = *uninitialized*) /home/butkus/public_html/app/code/core/Mage/Core/Model/Layout.php:456
[10-Feb-2016 00:39:01 UTC] PHP  20. Mage_Adminhtml_Block_Dashboard_Grids->_prepareLayout() /home/butkus/public_html/app/code/core/Mage/Core/Block/Abstract.php:238
[10-Feb-2016 00:39:01 UTC] PHP  21. Mage_Core_Block_Abstract->toHtml() /home/butkus/public_html/app/code/core/Mage/Adminhtml/Block/Dashboard/Grids.php:64
[10-Feb-2016 00:39:01 UTC] PHP  22. Mage_Adminhtml_Block_Widget_Grid->_beforeToHtml() /home/butkus/public_html/app/code/core/Mage/Core/Block/Abstract.php:862
[10-Feb-2016 00:39:01 UTC] PHP  23. Mage_Adminhtml_Block_Widget_Grid->_prepareGrid() /home/butkus/public_html/app/code/core/Mage/Adminhtml/Block/Widget/Grid.php:632
[10-Feb-2016 00:39:01 UTC] PHP  24. Mage_Adminhtml_Block_Widget_Grid->_prepareMassactionBlock() /home/butkus/public_html/app/code/core/Mage/Adminhtml/Block/Widget/Grid.php:625
[10-Feb-2016 00:44:11 UTC] PHP Fatal error:  Call to a member function isAvailable() on a non-object in /home/butkus/public_html/app/code/core/Mage/Adminhtml/Block/Widget/Grid.php on line 572
[10-Feb-2016 00:44:11 UTC] PHP Stack trace:
[10-Feb-2016 00:44:11 UTC] PHP   1. {main}() /home/butkus/public_html/index.php:0
`
	rd := bytes.NewBuffer([]byte(l))
	stats, _, err := read(rd, int64(0))
	poe(err)
	t.Log(stats)
}

func TestDate(t *testing.T) {

	p, err := NewLogParser("foo")
	poe(err)

	lt, e := p.getDate([]byte(lineNotice))
	t.Log(lt, e)
	poe(e)
	lt, e = p.getDate([]byte(lineWarning))
	t.Log(lt, e)
	poe(e)
	lt, e = p.getDate([]byte(lineCrap))
	t.Log(lt, e)
	pone(e)
	lt, e = p.getDate([]byte(lineLongCrap))
	t.Log(lt, e)
	pone(e)
	lt, e = p.getDate([]byte(lineAlmost))
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
		lineAlmost: expectedFileAndLine{
			err: true,
		},
		lineMultiUnknown: expectedFileAndLine{
			line: 0,
			file: "/usr/local/opt/php56-yaml/yaml.so in Unknown",
			err:  false,
		},
		lineInBlanks: expectedFileAndLine{
			line: 3,
			file: "/Users/jan/go/src/github.com/foomo/phperrorlog_exporter/sepp in depp.php",
			err:  false,
		},
	}
	for line, expectation := range expectations {
		actualFile, actualLine, err := getFileAndLine([]byte(line))
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
		lineNotice:   "notice",
		lineWarning:  "warning",
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
	te, z, e := extractDateStringFromLine([]byte(lineNotice))
	t.Log("te", te, "z", z, "e", e)
}
