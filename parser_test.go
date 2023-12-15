package gortf

import (
	"reflect"
	"testing"
)

func TestRTFParser(t *testing.T) {
	parser := NewRtfParser()

	content := `{\rtf1\ansi{\fonttbl\f0\fswiss Helvetica;}{\colortbl;\red0\green0\blue0;\red255\green255\blue255;}\f0\pard` + "\n"
	content += `This is some {\b bold} text.\par` + "\n"
	content += `}`

	doc, _ := parser.ParseContent(content)

	expected := RtfDocument{
		Header: RtfHeader{
			Charset: CharacterSetAnsi,
			FontTable: map[FontRef]Font{
				0: Font{
					Name:       "Helvetica",
					Charset:    CharacterSetNone,
					FontFamily: FontFamilySwiss,
				},
			},
			ColorTable: map[ColorRef]Color{
				1: Color{0, 0, 0},
				2: Color{255, 255, 255},
			},
		},
		Body: []StyleBlock{
			StyleBlock{
				Painter: Painter{},
				Text:    "This is some ",
			},
			StyleBlock{
				Painter: Painter{
					Bold: true,
				},
				Text: "bold",
			},
			StyleBlock{
				Painter: Painter{},
				Text:    " text.",
			},
		},
	}

	if !reflect.DeepEqual(doc, expected) {
		t.Errorf("\n\nexpected: %v\n\nactual\t: %v", expected, doc)
	}
}

func TestParseIgnoreGroup(t *testing.T) {
	content := `{\*\expandedcolortbl;;}`

	parser := NewRtfParser()
	parser.ParseContent(content)

	if len(parser.tokens) != 0 {
		t.Errorf("\n\nexpected: %v\n\nactual\t: %v", 0, len(parser.tokens))
	}
}

func TestRTFToText(t *testing.T) {
	content := `{\rtf1\ansi{\fonttbl\f0\fswiss Helvetica;}\f0\pard Voici du texte en {\b gras}.\par}`

	parser := NewRtfParser()
	doc, _ := parser.ParseContent(content)
	text, _ := doc.ToText()

	expected := "Voici du texte en gras."

	if text != expected {
		t.Errorf("\n\nexpected: %v\n\nactual\t: %v", expected, text)
	}
}
