package gortf

type scanner struct {
	start   int
	current int
	source  string
	tokens  []token
}

func newScanner(source string) scanner {
	return scanner{
		start:   0,
		current: 0,
		source:  source,
		tokens:  []token{},
	}
}

func (s *scanner) scanTokens() {
	for !s.isAtEnd() {
		s.start = s.current

		s.scanToken()
	}
}

func (s *scanner) scanToken() {
	c := s.advance()

	switch c {
	case '{':
		s.addToken(newGroupToken())
	case '}':
		s.addToken(newGroupEndToken())
	case '\\':
		pc := s.peek()

		if pc == '\\' || pc == '{' || pc == '}' { // escaped characters

		} else if pc == '\n' { // CRLF

		} else if isAlphaLower(pc) { // control word
			for s.peek() != '\\' && s.peek() != '{' && s.peek() != '}' {
				s.advance()
			}

			slice := s.source[s.start:s.current]
			head, tail := splitAtFirstWhitespace(slice)

			s.addToken(newControlWordToken(head))

			if len(tail) > 0 {
				s.addToken(newTextToken(tail))
			}
		}
	}
}

func (s *scanner) addToken(token token) {
	s.tokens = append(s.tokens, token)
}

func (s *scanner) peek() byte {
	if s.isAtEnd() {
		return 0
	}

	return s.source[s.current]
}

func (s *scanner) peekNext() byte {
	if (s.current + 1) >= len(s.source) {
		return 0
	}

	return s.source[s.current+1]
}

func (s *scanner) isAtEnd() bool {
	return s.current >= len(s.source)
}

func (s *scanner) advance() byte {
	currentChar := s.source[s.current]

	s.current += 1

	return currentChar
}
