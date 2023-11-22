package gortf

import (
	"fmt"
	"strconv"
)

type tokenType int

const (
	tokenTypeNone tokenType = iota
	tokenTypeGroup
	tokenTypeGroupEnd
	tokenTypeText
	tokenTypeControlWord
	tokenTypeCRLF
	tokenTypeIgnorableDestination
)

type controlWordType int

const (
	controlWordTypeRtf controlWordType = iota + 1
	controlWordTypeAnsi
	controlWordTypeFontTable
	controlWordTypeFontNumber
	controlWordTypeFontSize
	controlWordTypeItalic
	controlWordTypeBold
	controlWordTypeUnderline
	controlWordTypeUnknown
)

func (c controlWordType) String() string {
	switch c {
	case controlWordTypeRtf:
		return "rtf"
	case controlWordTypeAnsi:
		return "ansi"
	case controlWordTypeFontTable:
		return "fonttbl"
	case controlWordTypeFontNumber:
		return "f"
	case controlWordTypeFontSize:
		return "fs"
	case controlWordTypeItalic:
		return "i"
	case controlWordTypeBold:
		return "b"
	case controlWordTypeUnderline:
		return "i"
	default:
		return "unknown"
	}
}

type token interface {
	tokenType() tokenType
}

type binaryToken struct {
	value []byte
}

func (b binaryToken) tokenType() tokenType {
	return tokenTypeText
}

func NewBinaryToken(value []byte) binaryToken {
	return binaryToken{
		value: value,
	}
}

type groupToken struct {
}

func (g groupToken) tokenType() tokenType {
	return tokenTypeGroup
}

func (g groupToken) String() string {
	return "{Group}"
}

func newGroupToken() groupToken {
	return groupToken{}
}

type groupEndToken struct {
}

func (g groupEndToken) tokenType() tokenType {
	return tokenTypeGroupEnd
}

func (g groupEndToken) String() string {
	return "{GroupEnd}"
}

func newGroupEndToken() groupEndToken {
	return groupEndToken{}
}

type crlfToken struct {
}

func (c crlfToken) tokenType() tokenType {
	return tokenTypeCRLF
}

func (c crlfToken) String() string {
	return "{CRLF}"
}

func newCrlfToken() crlfToken {
	return crlfToken{}
}

type ignorableDestination struct {
}

func (i ignorableDestination) tokenType() tokenType {
	return tokenTypeIgnorableDestination
}

func (i ignorableDestination) String() string {
	return "{IgnorableDestination}"
}

func newIgnorableDestination() ignorableDestination {
	return ignorableDestination{}
}

type controlWordToken struct {
	name            string
	controlWordType controlWordType
	parameter       int
}

func (c controlWordToken) tokenType() tokenType {
	return tokenTypeControlWord
}

func (c controlWordToken) String() string {
	return fmt.Sprintf("{ControlWord %s %s %d}", c.name, c.controlWordType, c.parameter)
}

func newControlWordToken(input string) (controlWordToken, error) {
	prefix := ""
	suffix := ""
	param := -1

	suffixIndex := getSuffixIndex(input)
	if suffixIndex == -1 {
		prefix = input
	} else {
		prefix = input[:suffixIndex]
		suffix = input[suffixIndex:]

		p, err := strconv.Atoi(suffix)
		if err != nil {
			return controlWordToken{}, err
		}

		param = p
	}

	controlType := getControlTypeFromPrefix(prefix)

	return controlWordToken{
		name:            prefix,
		controlWordType: controlType,
		parameter:       param,
	}, nil
}

func getSuffixIndex(text string) int {
	suffixIndex := -1
	for idx := range text {
		if isNumber(text[idx]) || text[idx] == '-' {
			suffixIndex = idx
			break
		}
	}

	return suffixIndex
}

func getControlTypeFromPrefix(prefix string) controlWordType {
	switch prefix {
	case `\rtf`:
		return controlWordTypeRtf
	case `\ansi`:
		return controlWordTypeAnsi
	case `\fonttbl`:
		return controlWordTypeFontTable
	case `\f`:
		return controlWordTypeFontNumber
	case `\fs`:
		return controlWordTypeFontSize
	case `\i`:
		return controlWordTypeItalic
	case `\b`:
		return controlWordTypeBold
	case `\u`:
		return controlWordTypeUnderline
	default:
		return controlWordTypeUnknown
	}
}

type textToken struct {
	value string
}

func (t textToken) tokenType() tokenType {
	return tokenTypeText
}

func (t textToken) String() string {
	return fmt.Sprintf("{Text %s}", t.value)
}

func newTextToken(value string) textToken {
	return textToken{
		value: value,
	}
}
