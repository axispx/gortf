package gortf

import (
	"reflect"
	"testing"
)

func TestSimpleTokenization(t *testing.T) {
	scanner := newScanner(`{\b Words in bold}`)
	scanner.scanTokens()

	expected := []token{
		groupToken{},
		controlWordToken{`\b`, controlWordTypeBold, -1},
		textToken{"Words in bold"},
		groupEndToken{},
	}

	if !reflect.DeepEqual(scanner.tokens, expected) {
		t.Errorf("expected: %v\nactual: %v", expected, scanner.tokens)
	}
}

func TestScanEntireFile(t *testing.T) {
	scanner := newScanner(`{ \rtf1\ansi{\fonttbl\f0\fswiss Helvetica;}\f0\pard Voici du texte en {\b gras}.\par }`)
	scanner.scanTokens()

	expected := []token{
		groupToken{},
		controlWordToken{`\rtf`, controlWordTypeRtf, 1},
		controlWordToken{`\ansi`, controlWordTypeAnsi, -1},
		groupToken{},
		controlWordToken{`\fonttbl`, controlWordTypeFontTable, -1},
		controlWordToken{`\f`, controlWordTypeFontNumber, 0},
		controlWordToken{`\fswiss`, controlWordTypeUnknown, -1},
		textToken{"Helvetica;"},
		groupEndToken{},
		controlWordToken{`\f`, controlWordTypeFontNumber, 0},
		controlWordToken{`\pard`, controlWordTypeUnknown, -1},
		textToken{"Voici du texte en "},
		groupToken{},
		controlWordToken{`\b`, controlWordTypeBold, -1},
		textToken{"gras"},
		groupEndToken{},
		textToken{"."},
		controlWordToken{`\par`, controlWordTypeUnknown, -1},
		groupEndToken{},
	}

	if !reflect.DeepEqual(scanner.tokens, expected) {
		t.Errorf("\n\nexpected: %v\n\nactual\t: %v", expected, scanner.tokens)
	}
}

func TestEscapedText(t *testing.T) {
	content := `{\f0\fs24 \cf0 test de code \` + "\n"
	content += `if (a == b) \{\` + "\n"
	content += `    test();\` + "\n"
	content += `\} else \{\` + "\n"
	content += `    return;\` + "\n"
	content += `\}}`

	scanner := newScanner(content)
	scanner.scanTokens()

	expected := []token{
		groupToken{},
		controlWordToken{`\f`, controlWordTypeFontNumber, 0},
		controlWordToken{`\fs`, controlWordTypeFontSize, 24},
		controlWordToken{`\cf`, controlWordTypeUnknown, 0},
		textToken{"test de code "},
		crlfToken{},
		textToken{"if (a == b) "},
		textToken{"{"},
		crlfToken{},
		textToken{"    test();"},
		crlfToken{},
		textToken{"} else "},
		textToken{"{"},
		crlfToken{},
		textToken{"    return;"},
		crlfToken{},
		textToken{"}"},
		groupEndToken{},
	}

	if !reflect.DeepEqual(scanner.tokens, expected) {
		t.Errorf("\n\nexpected: %v\n\nactual\t: %v", expected, scanner.tokens)
	}
}

func TestIgnorableDestination(t *testing.T) {
	content := `{\*\expandedcolortbl;;}`

	scanner := newScanner(content)
	scanner.scanTokens()

	expected := []token{
		groupToken{},
		ignorableDestination{},
		controlWordToken{`\expandedcolortbl;`, controlWordTypeUnknown, -1},
		groupEndToken{},
	}

	if !reflect.DeepEqual(scanner.tokens, expected) {
		t.Errorf("\n\nexpected: %v\n\nactual\t: %v", expected, scanner.tokens)
	}
}

func TestShouldParseControlSymbolEndingSemicolon(t *testing.T) {
	content := `{\red255\blue255;}`

	scanner := newScanner(content)
	scanner.scanTokens()

	expected := []token{
		groupToken{},
		controlWordToken{`\red`, controlWordTypeUnknown, 255},
		controlWordToken{`\blue`, controlWordTypeUnknown, 255},
		groupEndToken{},
	}

	if !reflect.DeepEqual(scanner.tokens, expected) {
		t.Errorf("\n\nexpected: %v\n\nactual\t: %v", expected, scanner.tokens)
	}
}
