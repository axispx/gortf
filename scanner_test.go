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
		controlWordToken{
			name:            `\b`,
			controlWordType: controlWordTypeBold,
			parameter:       -1,
		},
		textToken{
			value: "Words in bold",
		},
		groupEndToken{},
	}

	if !reflect.DeepEqual(scanner.tokens, expected) {
		t.Errorf("expected: %v\nactual: %v", expected, scanner.tokens)
	}
}
