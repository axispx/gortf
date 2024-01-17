package gortf

import (
	"fmt"
)

type Font struct {
	Name       string
	Charset    CharacterSet
	FontFamily FontFamily
}

type CharacterSet int

const (
	CharacterSetNone CharacterSet = iota
	CharacterSetAnsi
	CharacterSetMac
	CharacterSetPc
	CharacterSetPca
)

func (c CharacterSet) String() string {
	switch c {
	case CharacterSetNone:
		return "None"
	case CharacterSetAnsi:
		return "ANSI"
	case CharacterSetMac:
		return "MAC"
	case CharacterSetPc:
		return "PC"
	case CharacterSetPca:
		return "PCA"
	default:
		return "Unknown"
	}
}

func characterSetFromToken(tkn token) CharacterSet {
	if tkn.tokenType() == tokenTypeControlWord {
		controlWord := tkn.(controlWordToken)

		switch controlWord.name {
		case `\ansi`:
			return CharacterSetAnsi
		case `\mac`:
			return CharacterSetMac
		case `\pc`:
			return CharacterSetPc
		case `\pca`:
			return CharacterSetPca
		}
	}

	return CharacterSetNone
}

type FontFamily int

const (
	FontFamilyNil FontFamily = iota
	FontFamilyRoman
	FontFamilySwiss
	FontFamilyModern
	FontFamilyScript
	FontFamilyDecor
	FontFamilyTech
	FontFamilyBidi
)

func (f FontFamily) String() string {
	switch f {
	case FontFamilyNil:
		return "Nil"
	case FontFamilyRoman:
		return "Roman"
	case FontFamilySwiss:
		return "Swiss"
	case FontFamilyModern:
		return "Modern"
	case FontFamilyScript:
		return "Script"
	case FontFamilyDecor:
		return "Decor"
	case FontFamilyTech:
		return "Tech"
	case FontFamilyBidi:
		return "Bidi"
	default:
		return "Nil"
	}
}

func fontFamilyFromToken(tkn token) FontFamily {
	if tkn.tokenType() == tokenTypeControlWord {
		controlWord := tkn.(controlWordToken)

		switch controlWord.name {
		case `\fnil`:
			return FontFamilyNil
		case `\froman`:
			return FontFamilyRoman
		case `\fswiss`:
			return FontFamilySwiss
		case `\fmodern`:
			return FontFamilyModern
		case `\fscript`:
			return FontFamilyScript
		case `\fdecor`:
			return FontFamilyDecor
		case `\ftech`:
			return FontFamilyTech
		case `\fbidi`:
			return FontFamilyBidi
		}
	}

	return FontFamilyNil
}

type Color struct {
	R int
	G int
	B int
}

func (c Color) valid() bool {
	return c.R >= 0 && c.G >= 0 && c.B >= 0
}

type Style struct {
	Name   string
	Number int
}

type TableRef uint16

type FontTable map[TableRef]Font
type ColorTable map[TableRef]Color
type Stylesheet map[string]Style

type RtfHeader struct {
	Charset    CharacterSet
	FontTable  FontTable
	ColorTable ColorTable
	Stylesheet Stylesheet
}

func (r RtfHeader) String() string {
	return fmt.Sprintf("{Charset: %s, FontTable: %v, ColorTable: %v}", r.Charset, r.FontTable, r.ColorTable)
}
