package gortf

import (
	"fmt"
	"os"
	"strings"
)

type Painter struct {
	FontRef   FontRef
	FontSize  int
	Bold      bool
	Italic    bool
	Underline bool
}

func (p Painter) String() string {
	return fmt.Sprintf("Painter: {fr: %d, fs: %d, b: %v, i: %v, u: %v}", p.FontRef, p.FontSize, p.Bold, p.Italic, p.Underline)
}

type StyleBlock struct {
	Painter Painter
	Text    string
}

func (s StyleBlock) String() string {
	return fmt.Sprintf("StyleBlock: {Painter: {%v}, Text: %s}", s.Painter, s.Text)
}

type RtfDocument struct {
	Header RtfHeader
	Body   []StyleBlock
}

func (r RtfDocument) String() string {
	return fmt.Sprintf("{Header: %v, Body: %v}", r.Header, r.Body)
}

func (r *RtfDocument) pushToBody(sb StyleBlock) {
	r.Body = append(r.Body, sb)
}

func (r *RtfDocument) popFromBody() StyleBlock {
	if len(r.Body) == 0 {
		panic("too many group endings")
	}

	index := len(r.Body) - 1
	element := r.Body[index]
	r.Body = r.Body[:index]

	return element
}

type RtfParser struct {
	tokens       []token
	painterStack []*Painter
	cursor       int
}

func NewRtfParser() RtfParser {
	return RtfParser{
		tokens:       []token{},
		painterStack: []*Painter{},
		cursor:       0,
	}
}

func (r *RtfParser) ParseFile(filePath string) (RtfDocument, error) {
	buf, err := os.ReadFile(filePath)
	if err != nil {
		return RtfDocument{}, err
	}

	doc, err := r.ParseContent(string(buf))
	if err != nil {
		return RtfDocument{}, err
	}

	return doc, nil
}

func (r *RtfParser) ParseContent(content string) (RtfDocument, error) {
	scanner := newScanner(content)
	scanner.scanTokens()
	r.tokens = scanner.tokens

	doc, err := r.parse()
	if err != nil {
		return RtfDocument{}, err
	}

	return doc, nil
}

func (r *RtfParser) parse() (RtfDocument, error) {
	doc := RtfDocument{}
	doc.Header = r.parseHeader()

	r.pushPainter(Painter{})
	for _, tkn := range r.tokens {
		switch tkn.tokenType() {
		case tokenTypeGroup:
			r.pushPainter(Painter{})

		case tokenTypeGroupEnd:
			r.popPainter()

		case tokenTypeControlWord:
			currentPainter := r.lastPainter()
			controlWord := tkn.(controlWordToken)

			switch controlWord.controlWordType {
			case controlWordTypeFontNumber:
				currentPainter.FontRef = FontRef(controlWord.parameter)
			case controlWordTypeBold:
				currentPainter.Bold = true
			case controlWordTypeItalic:
				currentPainter.Italic = true
			case controlWordTypeUnderline:
				currentPainter.Underline = true
			}

		case tokenTypeText:
			currentPainter := r.lastPainter()
			tt := tkn.(textToken)

			doc.pushToBody(StyleBlock{
				Painter: *currentPainter,
				Text:    tt.value,
			})
		}

	}

	return doc, nil
}

func (r *RtfParser) parseHeader() RtfHeader {
	r.cursor = 0
	header := RtfHeader{Charset: CharacterSetAnsi}

	for !r.isAtEnd() {
		currentToken := r.advance()
		nextToken := r.peek()

		if currentToken.tokenType() == tokenTypeGroup && nextToken.tokenType() == tokenTypeControlWord {
			controlWord := nextToken.(controlWordToken)
			if controlWord.controlWordType == controlWordTypeFontTable {
				fontTableTokens := r.consumeTokensUntilMatchingBracket()
				header.FontTable = r.parseFontTable(fontTableTokens)
				break
			}
		}

		if currentToken != nil {
			charset := characterSetFromToken(currentToken)
			if charset != CharacterSetNone {
				header.Charset = charset
			}
		}

		if currentToken == nil && nextToken == nil {
			break
		}
	}

	return header
}

func (r *RtfParser) parseFontTable(fontTableTokens []token) FontTable {
	table := make(FontTable)
	var currentKey FontRef = 0
	currentFont := Font{FontFamily: FontFamilyNil}

	for _, tkn := range fontTableTokens {
		switch tkn.tokenType() {
		case tokenTypeControlWord:
			controlWord := tkn.(controlWordToken)

			switch controlWord.controlWordType {
			case controlWordTypeFontNumber:
				table[currentKey] = currentFont
				currentKey = FontRef(controlWord.parameter)
			case controlWordTypeUnknown:
				fontFamily := fontFamilyFromName(controlWord.name)
				if fontFamily != FontFamilyNone {
					currentFont.FontFamily = fontFamily
				}
			}
		case tokenTypeText:
			tt := tkn.(textToken)
			currentFont.Name = strings.TrimSuffix(tt.value, ";")
		case tokenTypeGroupEnd:
			table[currentKey] = currentFont
		}
	}

	return table
}

func (r *RtfParser) consumeTokensUntilMatchingBracket() []token {
	tokens := []token{}
	count := 0

	for !r.isAtEnd() {
		currentToken := r.advance()

		switch currentToken.tokenType() {
		case tokenTypeGroup:
			count += 1
		case tokenTypeGroupEnd:
			count -= 1
		}

		tokens = append(tokens, currentToken)

		if count < 0 {
			break
		}
	}

	return tokens
}

func (r *RtfParser) advance() token {
	if len(r.tokens) == 0 {
		panic("no tokens")
	}

	t := r.tokens[r.cursor]

	r.tokens = append(r.tokens[:r.cursor], r.tokens[r.cursor+1:]...)

	return t
}

func (r *RtfParser) peek() token {
	if r.isAtEnd() {
		return nil
	}

	return r.tokens[r.cursor]
}

func (r *RtfParser) isAtEnd() bool {
	return r.cursor >= len(r.tokens)
}

func (r *RtfParser) popToken() token {
	if len(r.tokens) == 0 {
		panic("too many group endings")
	}

	index := len(r.tokens) - 1
	element := r.tokens[index]
	r.tokens = r.tokens[:index]

	return element
}

func (r *RtfParser) pushPainter(p Painter) {
	r.painterStack = append(r.painterStack, &p)
}

func (r *RtfParser) popPainter() Painter {
	if len(r.painterStack) == 0 {
		panic("too many group endings")
	}

	index := len(r.painterStack) - 1
	element := r.painterStack[index]
	r.painterStack = r.painterStack[:index]

	return *element
}

func (r *RtfParser) lastPainter() *Painter {
	topIndex := len(r.painterStack) - 1

	if topIndex < 0 {
		panic("malformed painter stack")
	}

	return r.painterStack[topIndex]
}
