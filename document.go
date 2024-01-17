package gortf

import (
	"encoding/json"
	"strings"
	"time"
)

type RtfInformationGroup struct {
	Title           string
	Subject         string
	Author          string
	Manager         string
	Company         string
	Operator        string
	Category        string
	Keywords        string
	Comment         string
	Version         int
	DocumentComment string
	BaseAddress     string
	CreationTime    *time.Time
	RevisionTime    *time.Time
	LastPrintTime   *time.Time
	BackupTime      *time.Time
}

type RtfDocument struct {
	Header           RtfHeader
	InformationGroup RtfInformationGroup
	Body             []StyleBlock
}

func (r RtfDocument) String() string {
	b, _ := json.Marshal(r)
	return string(b)
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

func (r *RtfDocument) ToText() (string, error) {
	var sb strings.Builder

	for _, b := range r.Body {
		_, err := sb.WriteString(b.Text)
		if err != nil {
			return "", err
		}
	}

	return sb.String(), nil
}

func (r *RtfDocument) ToHTML() (string, error) {
	return RTFToHTML(r)
}
