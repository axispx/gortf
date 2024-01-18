package gortf

import (
	"encoding/json"
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
	tokenTypeIgnorable
)

type controlWordType int

const (
	controlWordTypeUnknown controlWordType = iota

	// prolog
	controlWordTypeRtf

	// character set
	controlWordTypeCharacterSet

	// font table
	controlWordTypeFontTable
	controlWordTypeFontNumber
	controlWordTypeFontSize
	controlWordTypeFontFamily
	controlWordTypeFontCharset
	controlWordTypeFontAlternative
	controlWordTypeFontPitch
	controlWordTypeFontPanose
	controlWordTypeFontName
	controlWordTypeFontBias

	// file table
	controlWordTypeFileTable
	controlWordTypeFile
	controlWordTypeFileID
	controlWordTypeFileRelative
	controlWordTypeFileOSNumber
	controlWordTypeFileValidMac
	controlWordTypeFileValidDOS
	controlWordTypeFileValidNTFS
	controlWordTypeFileValidHPFS
	controlWordTypeFileValidNetwork

	// color table
	controlWordTypeColorTable
	controlWordTypeColorRed
	controlWordTypeColorGreen
	controlWordTypeColorBlue

	// stylesheet
	controlWordTypeStylesheet
	controlWordTypeStyleCharacter
	controlWordTypeStyleParagraph
	controlWordTypeStyleSection
	controlWordTypeStyleAdditive
	controlWordTypeStyleBasedOn
	controlWordTypeStyleNext
	controlWordTypeStyleAutoUpdate
	controlWordTypeStyleHidden
	controlWordTypeStylePersonalEmail
	controlWordTypeStyleEmailCompose
	controlWordTypeStyleEmailReply
	controlWordTypeStyleKeycode
	controlWordTypeStyleAltModifierKey
	controlWordTypeStyleShiftModifierKey
	controlWordTypeStyleControlModifierKey
	controlWordTypeStyleFunctionKey

	// list table
	controlWordTypeListTable
	controlWordTypeList
	controlWordTypeListID
	controlWordTypeListTemplateID
	controlWordTypeListSimple
	controlWordTypeListHybrid
	controlWordTypeListRestartSection
	controlWordTypeListName
	controlWordTypeListLevel
	controlWordTypeListLevelStartAt
	controlWordTypeListLevelNfc
	controlWordTypeListLevelJc
	controlWordTypeListLevelNfcn
	controlWordTypeListLevelJcn
	controlWordTypeListLevelOld
	controlWordTypeListLevelPrev
	controlWordTypeListPrevSpace
	controlWordTypeListLevelIndent
	controlWordTypeListLevelSpace
	controlWordTypeListLevelText
	controlWordTypeListLevelNumbers
	controlWordTypeListLevelFollow
	controlWordTypeListLevelNoRestart

	// list override table
	controlWordListOverrideTable
	controlWordTypeListOverride
	controlWordTypeListOverrideListID
	controlWordTypeListOverrideCount
	controlWordTypeListOverrideLs
	controlWordTypeListOverrideLevel
	controlWordTypeListOverrideLevelStartAt
	controlWordTypeListOverrideLevelFormat

	// information group
	controlWordTypeInfo
	controlWordTypeInfoTitle
	controlWordTypeInfoSubject
	controlWordTypeInfoAuthor
	controlWordTypeInfoManager
	controlWordTypeInfoCompany
	controlWordTypeInfoOperator
	controlWordTypeInfoCategory
	controlWordTypeInfoKeywords
	controlWordTypeInfoComment
	controlWordTypeInfoVersion
	controlWordTypeInfoDoccom
	controlWordTypeInfoHlinkBase

	// character formatting
	controlWordTypeItalic
	controlWordTypeBold
	controlWordTypeUnderline
	controlWordTypeUnderlineNone
	controlWordTypeSuperscript
	controlWordTypeSubscript
	controlWordTypeSmallcaps
	controlWordTypeStrikethrough
)

func (c controlWordType) String() string {
	switch c {
	// prolog
	case controlWordTypeRtf:
		return "rtf"

	// character set
	case controlWordTypeCharacterSet:
		return "characterset"

	// font table
	case controlWordTypeFontTable:
		return "fonttbl"
	case controlWordTypeFontNumber:
		return "f"
	case controlWordTypeFontFamily:
		return "fontfamily"
	case controlWordTypeFontCharset:
		return "fcharset"
	case controlWordTypeFontPitch:
		return "fprq"
	case controlWordTypeFontPanose:
		return "panose"
	case controlWordTypeFontBias:
		return "fbias"
	case controlWordTypeFontName:
		return "fname"
	case controlWordTypeFontAlternative:
		return "falt"

	// color table
	case controlWordTypeColorTable:
		return "colortbl"
	case controlWordTypeColorRed:
		return "red"
	case controlWordTypeColorGreen:
		return "green"
	case controlWordTypeColorBlue:
		return "blue"

	// stylesheet
	case controlWordTypeStylesheet:
		return "stylesheet"
	case controlWordTypeStyleCharacter:
		return "cs"
	case controlWordTypeStyleParagraph:
		return "s"
	case controlWordTypeStyleSection:
		return "ds"
	case controlWordTypeStyleAdditive:
		return "additive"
	case controlWordTypeStyleBasedOn:
		return "sbasedon"
	case controlWordTypeStyleNext:
		return "snext"
	case controlWordTypeStyleAutoUpdate:
		return "sautoupd"
	case controlWordTypeStyleHidden:
		return "shidden"
	case controlWordTypeStylePersonalEmail:
		return "spersonal"
	case controlWordTypeStyleEmailCompose:
		return "scompose"
	case controlWordTypeStyleEmailReply:
		return "sreply"
	case controlWordTypeStyleKeycode:
		return "keycode"
	case controlWordTypeStyleAltModifierKey:
		return "alt"
	case controlWordTypeStyleShiftModifierKey:
		return "shift"
	case controlWordTypeStyleControlModifierKey:
		return "ctrl"
	case controlWordTypeStyleFunctionKey:
		return "fn"

	// information group
	case controlWordTypeInfo:
		return "info"
	case controlWordTypeInfoTitle:
		return "title"
	case controlWordTypeInfoSubject:
		return "subject"
	case controlWordTypeInfoAuthor:
		return "author"
	case controlWordTypeInfoManager:
		return "manager"
	case controlWordTypeInfoCompany:
		return "company"
	case controlWordTypeInfoOperator:
		return "operator"
	case controlWordTypeInfoCategory:
		return "category"
	case controlWordTypeInfoKeywords:
		return "keywords"
	case controlWordTypeInfoComment:
		return "comment"
	case controlWordTypeInfoVersion:
		return "version"
	case controlWordTypeInfoDoccom:
		return "doccom"
	case controlWordTypeInfoHlinkBase:
		return "hlinkbase"

	// character formatting
	case controlWordTypeItalic:
		return "i"
	case controlWordTypeBold:
		return "b"
	case controlWordTypeUnderline:
		return "ul"
	case controlWordTypeUnderlineNone:
		return "ulnone"
	case controlWordTypeSuperscript:
		return "super"
	case controlWordTypeSubscript:
		return "sub"
	case controlWordTypeSmallcaps:
		return "scaps"
	case controlWordTypeStrikethrough:
		return "strike"

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

type controlWordToken struct {
	name            string
	controlWordType controlWordType
	parameter       int
}

func (c controlWordToken) tokenType() tokenType {
	return tokenTypeControlWord
}

func (c controlWordToken) String() string {
	b, _ := json.Marshal(c)
	return string(b)
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

	controlType := getControlWordTypeFromPrefix(prefix)

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

func getControlWordTypeFromPrefix(prefix string) controlWordType {
	switch prefix {
	// prolog
	case `\rtf`:
		return controlWordTypeRtf

	// character set
	case `\ansi`, `\mac`, `\pc`, `\pca`:
		return controlWordTypeCharacterSet

	// font table
	case `\fonttbl`:
		return controlWordTypeFontTable
	case `\f`:
		return controlWordTypeFontNumber
	case `\fs`:
		return controlWordTypeFontSize
	case `\fnil`, `\froman`, `\fswiss`, `\fmodern`, `\fscript`, `\fdecor`, `\ftech`, `\fbidi`:
		return controlWordTypeFontFamily
	case `\fcharset`:
		return controlWordTypeFontCharset
	case `\falt`:
		return controlWordTypeFontAlternative
	case `\fprq`:
		return controlWordTypeFontPitch
	case `\*\panose`:
		return controlWordTypeFontPanose
	case `\*\fname`:
		return controlWordTypeFontName
	case `\fbias`:
		return controlWordTypeFontBias

	case `\colortbl`:
		return controlWordTypeColorTable
	case `\red`:
		return controlWordTypeColorRed
	case `\green`:
		return controlWordTypeColorGreen
	case `\blue`:
		return controlWordTypeColorBlue

	case `\stylesheet`:
		return controlWordTypeStylesheet
	case `\*\cs`:
		return controlWordTypeStyleCharacter
	case `\s`:
		return controlWordTypeStyleParagraph
	case `\ds`:
		return controlWordTypeStyleSection
	case `\additive`:
		return controlWordTypeStyleAdditive
	case `\sbasedon`:
		return controlWordTypeStyleBasedOn
	case `\snext`:
		return controlWordTypeStyleNext
	case `\sautoupd`:
		return controlWordTypeStyleAutoUpdate
	case `\shidden`:
		return controlWordTypeStyleHidden
	case `\spersonal`:
		return controlWordTypeStylePersonalEmail
	case `\scompose`:
		return controlWordTypeStyleEmailCompose
	case `\sreply`:
		return controlWordTypeStyleEmailReply
	case `\keycode`:
		return controlWordTypeStyleKeycode
	case `\alt`:
		return controlWordTypeStyleAltModifierKey
	case `\shift`:
		return controlWordTypeStyleShiftModifierKey
	case `\ctrl`:
		return controlWordTypeStyleControlModifierKey
	case `\fn`:
		return controlWordTypeStyleFunctionKey

		// information group
	case `\info`:
		return controlWordTypeInfo
	case `\title`:
		return controlWordTypeInfoTitle
	case `\subject`:
		return controlWordTypeInfoSubject
	case `\author`:
		return controlWordTypeInfoAuthor
	case `\manager`:
		return controlWordTypeInfoManager
	case `\company`:
		return controlWordTypeInfoCompany
	case `\operator`:
		return controlWordTypeInfoOperator
	case `\category`:
		return controlWordTypeInfoCategory
	case `\keywords`:
		return controlWordTypeInfoKeywords
	case `\comment`:
		return controlWordTypeInfoComment
	case `\version`:
		return controlWordTypeInfoVersion
	case `\doccom`:
		return controlWordTypeInfoDoccom
	case `\hlinkbase`:
		return controlWordTypeInfoHlinkBase

	// character formatting
	case `\i`:
		return controlWordTypeItalic
	case `\b`:
		return controlWordTypeBold
	case `\ul`:
		return controlWordTypeUnderline
	case `\ulnone`:
		return controlWordTypeUnderlineNone
	case `\super`:
		return controlWordTypeSuperscript
	case `\sub`:
		return controlWordTypeSubscript
	case `\scaps`:
		return controlWordTypeSmallcaps
	case `\strike`:
		return controlWordTypeStrikethrough

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
