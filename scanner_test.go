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
	content := `{\rtf1\ansi{\fonttbl\f0\fswiss Helvetica;}{\colortbl;\red0\green0\blue0;}`
	content += `{\stylesheet{\s0\snext0\ql\nowidctlpar\hyphpar0\ltrpar\cf17\dbch\af9\langfe2052\dbch\af13\afs24\alang1081\kerning0\loch\f3\fs24\lang1033 Normal;}}`
	content += `\f0\pard Voici du texte en {\b gras}.\par }`

	scanner := newScanner(content)
	scanner.scanTokens()

	expected := []token{
		groupToken{},
		controlWordToken{`\rtf`, controlWordTypeRtf, 1},
		controlWordToken{`\ansi`, controlWordTypeCharacterSet, -1},
		groupToken{},
		controlWordToken{`\fonttbl`, controlWordTypeFontTable, -1},
		controlWordToken{`\f`, controlWordTypeFontNumber, 0},
		controlWordToken{`\fswiss`, controlWordTypeFontFamily, -1},
		textToken{"Helvetica;"},
		groupEndToken{},
		groupToken{},
		controlWordToken{`\colortbl`, controlWordTypeColorTable, -1},
		controlWordToken{`\red`, controlWordTypeColorRed, 0},
		controlWordToken{`\green`, controlWordTypeColorGreen, 0},
		controlWordToken{`\blue`, controlWordTypeColorBlue, 0},
		groupEndToken{},
		groupToken{},
		controlWordToken{`\stylesheet`, controlWordTypeStylesheet, -1},
		groupToken{},
		controlWordToken{`\s`, controlWordTypeStyleParagraph, 0},
		controlWordToken{`\snext`, controlWordTypeStyleNext, 0},
		controlWordToken{`\ql`, controlWordTypeUnknown, -1},
		controlWordToken{`\nowidctlpar`, controlWordTypeUnknown, -1},
		controlWordToken{`\hyphpar`, controlWordTypeUnknown, 0},
		controlWordToken{`\ltrpar`, controlWordTypeUnknown, -1},
		controlWordToken{`\cf`, controlWordTypeUnknown, 17},
		controlWordToken{`\dbch`, controlWordTypeUnknown, -1},
		controlWordToken{`\af`, controlWordTypeUnknown, 9},
		controlWordToken{`\langfe`, controlWordTypeUnknown, 2052},
		controlWordToken{`\dbch`, controlWordTypeUnknown, -1},
		controlWordToken{`\af`, controlWordTypeUnknown, 13},
		controlWordToken{`\afs`, controlWordTypeUnknown, 24},
		controlWordToken{`\alang`, controlWordTypeUnknown, 1081},
		controlWordToken{`\kerning`, controlWordTypeUnknown, 0},
		controlWordToken{`\loch`, controlWordTypeUnknown, -1},
		controlWordToken{`\f`, controlWordTypeFontNumber, 3},
		controlWordToken{`\fs`, controlWordTypeFontSize, 24},
		controlWordToken{`\lang`, controlWordTypeUnknown, 1033},
		textToken{"Normal;"},
		groupEndToken{},
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

	expected := []token{}

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
		controlWordToken{`\red`, controlWordTypeColorRed, 255},
		controlWordToken{`\blue`, controlWordTypeColorBlue, 255},
		groupEndToken{},
	}

	if !reflect.DeepEqual(scanner.tokens, expected) {
		t.Errorf("\n\nexpected: %v\n\nactual\t: %v", expected, scanner.tokens)
	}
}
