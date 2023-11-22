package gortf

func splitAtFirstWhitespace(text string) (string, string) {
	var head string
	var tail string

	whitespaceIndex := -1
	for idx := range text {
		if isWhitespace(text[idx]) {
			whitespaceIndex = idx
			break
		}
	}

	if whitespaceIndex > -1 {
		head = text[:whitespaceIndex]
		tail = text[whitespaceIndex+1:]
	} else {
		head = text
	}

	return head, tail
}

func isWhitespace(ch byte) bool {
	return ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r' || ch == '\v' || ch == '\f'
}

func isAlphaLower(c byte) bool {
	return c >= 'a' && c <= 'z'
}

func isNumber(c byte) bool {
	return c >= '0' && c <= '9'
}
