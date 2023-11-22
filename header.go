package gortf

import (
	"fmt"
	"strings"
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

		// TODO: implement the rest
		switch controlWord.controlWordType {
		case controlWordTypeAnsi:
			return CharacterSetAnsi
		}
	}

	return CharacterSetNone
}

type FontFamily int

const (
	FontFamilyNone FontFamily = iota
	FontFamilyNil
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
	case FontFamilyNone:
		return "None"
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
		return "Unknown"
	}
}

func fontFamilyFromName(name string) FontFamily {
	switch {
	case strings.HasPrefix(name, `\fnil`):
		return FontFamilyNil
	case strings.HasPrefix(name, `\froman`):
		return FontFamilyRoman
	case strings.HasPrefix(name, `\fswiss`):
		return FontFamilySwiss
	default:
		return FontFamilyNone
	}
}

type FontRef uint16
type FontTable map[FontRef]Font

type RtfHeader struct {
	Charset   CharacterSet
	FontTable FontTable
}

func (r RtfHeader) String() string {
	return fmt.Sprintf("{Charset: %s, FontTable: %v}", r.Charset, r.FontTable)
}
