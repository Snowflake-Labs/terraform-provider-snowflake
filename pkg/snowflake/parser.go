package snowflake

import (
	"strings"
)

// ViewSelectStatementExtractor is a simplistic parser that only exists to extract the select statement from a
// create view statement
type ViewSelectStatementExtractor struct {
	input string
	pos   int
}

func NewViewSelectStatementExtractor(input string) *ViewSelectStatementExtractor {
	return &ViewSelectStatementExtractor{
		input: input,
	}
}

func (e *ViewSelectStatementExtractor) Extract() (string, error) {
	e.consumeToken("create")
	// e.consumeSpace()

	return "", nil
}

func (e *ViewSelectStatementExtractor) remainingString() string {
	if e.pos >= len(e.input) {
		return ""
	}
	return e.input[e.pos:]
}

func max(x, y int) int {
	if x > y {
		return x
	}
	return y
}

func (e *ViewSelectStatementExtractor) nextN(n int) string {
	if e.pos >= len(e.input) {
		return ""
	}

	if e.pos+n > len(e.input) {
		n = len(e.input) - e.pos
	}
	return e.input[e.pos : e.pos+n]
}

func (e *ViewSelectStatementExtractor) consumeToken(t string) {
	lenT := len(t)

	if len(e.remainingString()) < lenT {
		return
	}

	if strings.EqualFold(e.nextN(lenT), t) {
		e.pos += len(t)
	}
}

func (e *ViewSelectStatementExtractor) done() bool {
	return e.pos >= len(e.input)
}

// only safe to call if you know done() is not true
func (e *ViewSelectStatementExtractor) peek() byte {
	return e.input[e.pos]
}

// func (e *ViewSelectStatementExtractor) consumeSpace() {
// 	for {
// 		if e.done() || !unicode.IsSpace(e.peek()) {
// 			return
// 		}

// 		e.pos += 1
// 	}
// }
