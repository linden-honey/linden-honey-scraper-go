package parser

import (
	"testing"

	log "github.com/sirupsen/logrus"
)

var (
	testParser = NewDefaultParser(log.StandardLogger())
)

func TestParsePreviews(t *testing.T) {
	html := `
	<ul id="abc_list">
	<li><a href="/texts/1056899068.html">Всё идёт по плану</a></li>
	<li><a href=""></a></li>
	<li><a href="">Unknown</a></li>
	<li><a href="/texts/1056901056.html">Всё как у людей</a></li>
	</ul>
	`
	previews, _ := testParser.ParsePreviews(html)
	t.Log(previews)
	if len(previews) != 2 {
		t.Error("Parsing previews failed")
	}
}
