package gortf

import (
	"fmt"
	"os"
	"strings"
)

type Painter struct {
	FontRef   TableRef
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
	doc.InformationGroup = r.parseInformationGroup()

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
				currentPainter.FontRef = TableRef(controlWord.parameter)
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

			headerTableFound := false
			if controlWord.controlWordType == controlWordTypeFontTable {
				headerTableFound = true
				fontTableTokens := r.consumeTokensUntilMatchingBracket()
				header.FontTable = r.parseFontTable(fontTableTokens)
			} else if controlWord.controlWordType == controlWordTypeColorTable {
				headerTableFound = true
				colorTableTokens := r.consumeTokensUntilMatchingBracket()
				header.ColorTable = r.parseColorTable(colorTableTokens)
			} else if controlWord.controlWordType == controlWordTypeStylesheet {
				headerTableFound = true
				stylesheetTokens := r.consumeTokensUntilMatchingBracket()
				header.Stylesheet = r.parseStylesheet(stylesheetTokens)
			}

			if headerTableFound {
				if r.areMoreHeaderTablePresent() {
					continue
				} else {
					break
				}
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
	var currentKey TableRef = 0
	currentFont := Font{}

	for _, tkn := range fontTableTokens {
		switch tkn.tokenType() {
		case tokenTypeControlWord:
			controlWord := tkn.(controlWordToken)

			switch controlWord.controlWordType {
			case controlWordTypeFontNumber:
				table[currentKey] = currentFont
				currentKey = TableRef(controlWord.parameter)

			case controlWordTypeFontFamily:
				fontFamily := fontFamilyFromToken(tkn)
				currentFont.FontFamily = fontFamily
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

func (r *RtfParser) parseColorTable(colorTableTokens []token) ColorTable {
	table := make(ColorTable)
	var currentKey TableRef = 1
	var currentColor = Color{-1, -1, -1}

	for _, tkn := range colorTableTokens {
		switch tkn.tokenType() {
		case tokenTypeControlWord:
			controlWord := tkn.(controlWordToken)

			switch controlWord.controlWordType {
			case controlWordTypeColorRed:
				currentColor.R = controlWord.parameter
			case controlWordTypeColorGreen:
				currentColor.G = controlWord.parameter
			case controlWordTypeColorBlue:
				currentColor.B = controlWord.parameter
			}

			if currentColor.valid() {
				table[currentKey] = currentColor
				currentKey += 1
				currentColor = Color{-1, -1, -1}
			}
		}
	}

	return table
}

func (r *RtfParser) parseStylesheet(stylesheetTokens []token) Stylesheet {
	stylesheet := make(Stylesheet)

	currentStyle := Style{}

	for _, tkn := range stylesheetTokens {
		switch tkn.tokenType() {
		case tokenTypeText:
			tt := tkn.(textToken)
			currentStyle.Name = strings.TrimSuffix(tt.value, ";")
		}
	}

	return stylesheet
}

func (r *RtfParser) areMoreHeaderTablePresent() bool {
	nextToken := r.peek()
	nextToNextToken := r.peekN(1)

	if nextToken.tokenType() == tokenTypeGroup && nextToNextToken.tokenType() == tokenTypeControlWord {
		controlWord := nextToNextToken.(controlWordToken)

		switch controlWord.controlWordType {
		case controlWordTypeFontTable, controlWordTypeColorTable:
			return true
		default:
			return false
		}
	}

	return false
}

func (r *RtfParser) parseInformationGroup() RtfInformationGroup {
	r.cursor = 0

	// don't use advance here because we don't know if there is an information group yet
	currentToken := r.peek()
	nextToken := r.peekN(1)

	if currentToken == nil || nextToken == nil {
		return RtfInformationGroup{}
	}

	if currentToken.tokenType() == tokenTypeGroup && nextToken.tokenType() == tokenTypeControlWord {
		controlWord := nextToken.(controlWordToken)

		if controlWord.controlWordType != controlWordTypeInfo {
			return RtfInformationGroup{}
		}

		// advance into the information group
		r.advance()
	} else {
		return RtfInformationGroup{}
	}

	informationGroup := RtfInformationGroup{}
	for !r.isAtEnd() {
		currentToken := r.advance()
		nextToken := r.peek()

		if currentToken.tokenType() == tokenTypeGroup && nextToken.tokenType() == tokenTypeControlWord {
			controlWord := nextToken.(controlWordToken)
			tokens := r.consumeTokensUntilMatchingBracket()

			if controlWord.controlWordType == controlWordTypeInfoVersion {
				informationGroup.Version = controlWord.parameter
			} else if len(tokens) == 3 {
				textToken := tokens[1].(textToken)

				switch controlWord.controlWordType {
				case controlWordTypeInfoTitle:
					informationGroup.Title = textToken.value
				case controlWordTypeInfoSubject:
					informationGroup.Subject = textToken.value
				case controlWordTypeInfoAuthor:
					informationGroup.Author = textToken.value
				case controlWordTypeInfoManager:
					informationGroup.Manager = textToken.value
				case controlWordTypeInfoCompany:
					informationGroup.Company = textToken.value
				case controlWordTypeInfoOperator:
					informationGroup.Operator = textToken.value
				case controlWordTypeInfoCategory:
					informationGroup.Category = textToken.value
				case controlWordTypeInfoKeywords:
					informationGroup.Keywords = textToken.value
				case controlWordTypeInfoComment:
					informationGroup.Comment = textToken.value
				case controlWordTypeInfoDoccom:
					informationGroup.DocumentComment = textToken.value
				case controlWordTypeInfoHlinkBase:
					informationGroup.BaseAddress = textToken.value
				}
			}

		}

		if currentToken.tokenType() == tokenTypeGroupEnd {
			break
		}
	}

	return informationGroup
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

// advance returns the token at the cursor location
// and removes the token from the list of tokens
func (r *RtfParser) advance() token {
	if len(r.tokens) == 0 {
		panic("no tokens")
	}

	t := r.tokens[r.cursor]

	r.tokens = append(r.tokens[:r.cursor], r.tokens[r.cursor+1:]...)

	return t
}

func (r *RtfParser) advanceN(n int) {
	if len(r.tokens) < n {
		panic(fmt.Sprintf("only %d tokens left", len(r.tokens)))
	}

	r.tokens = r.tokens[r.cursor+n:]
}

func (r *RtfParser) peek() token {
	if r.isAtEnd() {
		return nil
	}

	return r.tokens[r.cursor]
}

func (r *RtfParser) peekN(n int) token {
	if r.cursor >= len(r.tokens)-n {
		return nil
	}

	return r.tokens[r.cursor+n]
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

func (r *RtfParser) readInformationGroup() {

}
