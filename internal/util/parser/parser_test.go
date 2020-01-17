package parser

import (
	"testing"
)

func TestPrasePrevie(t *testing.T) {
	html := `
	<ul id="abc_list">
	<li><a href="/texts/1056899068.html">Всё идёт по плану</a></li>
	<li><a href="">Unknown</a></li>
	<li><a href="/texts/1056901056.html">Всё как у людей</a></li>
	</ul>
	`
	previews := ParsePreviews(html)
	t.Log(previews)
	if len(previews) != 3 {
		t.Error("Parsing previews failed")
	}
}
