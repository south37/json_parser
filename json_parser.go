package json_parser

import (
	"errors"
	"fmt"
	"regexp"
	"unicode"
)

type Parser struct {
	buf []byte
}

func NewParser(s string) *Parser {
	return &Parser{buf: []byte(s)}
}

func (p *Parser) parse() (err error) {
	p.skipSpace()
	if p.buf[0] == '[' {
		return p.parseArray()
	} else if p.buf[0] == '{' {
		return p.parseObject()
	} else {
		return errors.New(fmt.Sprintf("json_parser: parse encountered unknown: %s", string(p.buf)))
	}
}

func (p *Parser) parseArray() (err error) {
	err = p.scan('[')
	if err != nil {
		return
	}
	p.skipSpace()
	if p.buf[0] == ']' {
		return p.scan(']')
	}
	err = p.parseExpression()
	if err != nil {
		return
	}
	p.skipSpace()
	for p.buf[0] != ']' {
		err = p.scan(',')
		if err != nil {
			return
		}
		p.skipSpace()
		err = p.parseExpression()
		if err != nil {
			return
		}
		p.skipSpace()
	}
	return p.scan(']')
}

func (p *Parser) parseObject() (err error) {
	err = p.scan('{')
	if err != nil {
		return
	}
	p.skipSpace()
	if p.buf[0] == '}' {
		return p.scan('}')
	}
	err = p.parseKeyValuePair()
	if err != nil {
		return
	}
	p.skipSpace()
	for p.buf[0] != '}' {
		err = p.scan(',')
		if err != nil {
			return
		}
		p.skipSpace()
		err = p.parseKeyValuePair()
		if err != nil {
			return
		}
		p.skipSpace()
	}
	return p.scan('}')
}

// Consider only ascii.
func (p *Parser) scan(r byte) (err error) {
	if p.buf[0] == r {
		p.buf = p.buf[1:]
		return
	} else {
		return errors.New(fmt.Sprintf("json_parser: scan expected %s got %s", string(r), string(p.buf[0])))
	}
}

func (p *Parser) parseKeyValuePair() (err error) {
	err = p.parseString()
	if err != nil {
		return
	}
	p.skipSpace()
	err = p.scan(':')
	if err != nil {
		return
	}
	p.skipSpace()
	return p.parseExpression()
}

func isNumber(r byte) bool {
	return r >= '0' && r <= '9'
}

func (p *Parser) parseExpression() (err error) {
	if p.buf[0] == '[' {
		return p.parseArray()
	} else if p.buf[0] == '{' {
		return p.parseObject()
	} else if p.buf[0] == 'n' {
		return p.parseNull()
	} else if p.buf[0] == '"' {
		return p.parseString()
	} else if p.buf[0] == '-' || isNumber(p.buf[0]) {
		return p.parseNumber()
	} else {
		return errors.New(fmt.Sprintf("json_parser: parseExpression encountered unknown: %s", string(p.buf)))
	}

}

func (p *Parser) parseNull() (err error) {
	if string(p.buf[0:4]) == `null` {
		p.buf = p.buf[4:]
		return
	} else {
		return errors.New(fmt.Sprintf("json_parser: parseNull expected null got %s", string(p.buf[0:4])))
	}
}

var strRe, _ = regexp.Compile(`^"(?:[^"]|\\")*"`)

func (p *Parser) parseString() (err error) {
	var loc = strRe.FindIndex(p.buf)
	if loc == nil {
		return errors.New(fmt.Sprintf("json_parser: parseString cannot parsed: %s", string(p.buf)))
	}
	p.buf = p.buf[loc[1]:]
	return
}

func matchString(s string) bool {
	return strRe.MatchString(s)
}

var numRe, _ = regexp.Compile(`^(?:0|-?[1-9][0-9]*)`)

func (p *Parser) parseNumber() (err error) {
	var loc = numRe.FindIndex(p.buf)
	if loc == nil {
		return errors.New(fmt.Sprintf("json_parser: parseNumber cannot parse %s", string(p.buf)))
	}
	p.buf = p.buf[loc[1]:]
	return
}

func matchNumber(s string) bool {
	return numRe.MatchString(s)
}

func (p *Parser) skipSpace() {
	for unicode.IsSpace(rune(p.buf[0])) {
		p.buf = p.buf[1:]
	}
}
