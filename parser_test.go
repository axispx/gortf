package gortf

import (
	"fmt"
	"testing"
)

func TestRTFParser(t *testing.T) {
	parser := NewRtfParser()

	content := `{\rtf-1\ansi{\fonttbl\f0\fswiss Helvetica;}\f0\pard` + "\n"
	content += `This is some {\b bold} text.\par` + "\n"
	content += `}`

	doc, _ := parser.ParseContent(content)

	fmt.Println(content)
	fmt.Println(parser.tokens)
	fmt.Println(doc)
}
