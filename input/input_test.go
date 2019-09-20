// SPDX-License-Identifier: MIT

package input

import (
	"io/ioutil"
	"log"
	"testing"

	"github.com/issue9/assert"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/simplifiedchinese"

	"github.com/caixw/apidoc/v5/message"
)

func TestParse(t *testing.T) {
	a := assert.New(t)
	erro := log.New(ioutil.Discard, "[ERRO]", 0)
	warn := log.New(ioutil.Discard, "[WARN]", 0)
	info := log.New(ioutil.Discard, "[INFO]", 0)
	h := message.NewHandler(message.NewLogHandlerFunc(erro, warn, info))
	a.NotNil(h)

	php := &Options{
		Lang:      "php",
		Dir:       "./testdata",
		Recursive: true,
		Encoding:  "gbk",
	}
	a.Panic(func() {
		Parse(h, php)
	})
	a.NotError(php.Sanitize())

	c := &Options{
		Lang:      "c++",
		Dir:       "./testdata",
		Recursive: true,
	}
	a.NotError(c.Sanitize())

	doc := Parse(h, php, c)
	a.NotNil(doc).
		Equal(1, len(doc.Apis)).
		Equal(doc.Version, "1.1.1")
	api := doc.Apis[0]
	a.Equal(api.Method, "GET")
}

func TestReadFile(t *testing.T) {
	a := assert.New(t)

	nop, err := readFile("./testdata/gbk.php", encoding.Nop)
	a.NotError(err).
		NotNil(nop).
		NotContains(string(nop), "这是一个 GBK 编码的文件")

	def, err := readFile("./testdata/gbk.php", nil)
	a.NotError(err).
		NotNil(def).
		NotContains(string(def), "这是一个 GBK 编码的文件")
	a.Equal(def, nop)

	data, err := readFile("./testdata/gbk.php", simplifiedchinese.GB18030)
	a.NotError(err).
		NotNil(data).
		Contains(string(data), "这是一个 GBK 编码的文件")
}
