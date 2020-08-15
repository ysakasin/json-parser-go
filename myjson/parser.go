package myjson

import (
	"errors"
	"fmt"
	"strconv"
)

func Parse(src string) (interface{}, error) {
	srcRunes := []rune(src)

	p := &parser{
		src:    srcRunes,
		srcLen: len(srcRunes),
		index:  0,
	}

	return p.json()
}

type parser struct {
	src    []rune
	srcLen int
	index  int
}

func (p *parser) json() (interface{}, error) {
	value, err := p.element()

	if p.index != p.srcLen {
		return nil, errors.New("source is still continuing")
	}

	return value, err
}

func (p *parser) element() (value interface{}, err error) {
	p.skip_white_space()
	value, err = p.value()
	p.skip_white_space()

	return
}

func (p *parser) value() (interface{}, error) {
	if p.current() == '{' {
		return p.object()
	}

	if p.current() == '[' {
		return p.array()
	}

	if p.current() == '"' {
		return p.string()
	}

	if p.consume("true") {
		return true, nil
	}

	if p.consume("false") {
		return false, nil
	}

	if p.consume("null") {
		return nil, nil
	}

	return p.number()
}

func (p *parser) object() (map[string]interface{}, error) {
	if err := p.expectRune('{'); err != nil {
		return nil, err
	}

	p.skip_white_space()

	if p.consumeRune('}') {
		return make(map[string]interface{}), nil
	}

	obj, err := p.members()
	if err != nil {
		return nil, err
	}

	if err := p.expectRune('}'); err != nil {
		return nil, err
	}

	return obj, nil
}

func (p *parser) members() (map[string]interface{}, error) {
	obj := make(map[string]interface{})

	key, value, err := p.member()
	if err != nil {
		return nil, err
	}
	obj[key] = value

	for p.current() == ',' {
		p.next()

		key, value, err := p.member()
		if err != nil {
			return nil, err
		}
		obj[key] = value
	}

	return obj, nil
}

func (p *parser) member() (string, interface{}, error) {
	p.skip_white_space()

	key, err := p.string()
	if err != nil {
		return "", nil, err
	}

	p.skip_white_space()

	if err := p.expectRune(':'); err != nil {
		return "", nil, err
	}

	value, err := p.element()
	if err != nil {
		return "", nil, err
	}

	return key, value, nil
}

func (p *parser) array() (arr []interface{}, err error) {
	if err := p.expectRune('['); err != nil {
		return nil, err
	}

	p.skip_white_space()

	if p.consumeRune(']') {
		arr = make([]interface{}, 0)
		return
	}

	arr, err = p.elements()
	if err != nil {
		return nil, err
	}

	if err := p.expectRune(']'); err != nil {
		return nil, err
	}

	return
}

func (p *parser) elements() ([]interface{}, error) {
	arr := make([]interface{}, 0)

	elm, err := p.element()
	if err != nil {
		return nil, err
	}
	arr = append(arr, elm)

	for p.current() == ',' {
		p.next()
		elm, err := p.element()
		if err != nil {
			return nil, err
		}

		arr = append(arr, elm)
	}

	return arr, nil
}

func (p *parser) string() (string, error) {
	if err := p.expectRune('"'); err != nil {
		return "", err
	}

	runes := make([]rune, 0)
	for p.current() != '"' && p.index < p.srcLen {
		r, err := p.character()
		if err != nil {
			return "", err
		}

		runes = append(runes, r)
	}

	if err := p.expectRune('"'); err != nil {
		return "", err
	}

	return string(runes), nil
}

func (p *parser) character() (rune, error) {
	cur := p.current()

	if cur < 0x20 {
		return 0, errors.New("hhh")
	}

	if cur == '"' {
		return 0, errors.New("bbb")
	}

	if cur == '\\' {
		p.next()
		return p.escape()
	}

	p.next()
	return cur, nil
}

func (p *parser) escape() (r rune, err error) {
	switch p.current() {
	case '"':
		r = 0x0022 // \" quotation mark
	case '\\':
		r = 0x005C // \\ reverse solidus
	case '/':
		r = 0x003F // \/ solidus
	case 'b':
		r = 0x0008 // \b backspace
	case 'f':
		r = 0x000C // \f form feed
	case 'n':
		r = 0x000A // \n line feed
	case 'r':
		r = 0x000D // \r carriage return
	case 't':
		r = 0x0009 // \t character tabulation
	case 'u':
		p.next()
		r, err = p.codePoint()
	default:
		err = errors.New("unexpected escape")
	}

	return
}

func (p *parser) codePoint() (rune, error) {
	begin := p.index

	if err := p.hex(); err != nil {
		return 0, err
	}
	if err := p.hex(); err != nil {
		return 0, err
	}
	if err := p.hex(); err != nil {
		return 0, err
	}
	if err := p.hex(); err != nil {
		return 0, err
	}

	ri, err := strconv.ParseUint(string(p.src[begin:p.index]), 16, 32)
	return rune(ri), err
}

func (p *parser) hex() error {
	if !isHex(p.current()) {
		return errors.New("expected hex")
	}

	p.next()
	return nil
}

func (p *parser) number() (interface{}, error) {
	begin := p.index

	if err := p.integer(); err != nil {
		return nil, err
	}

	end_integer := p.index

	if err := p.fraction(); err != nil {
		return nil, err
	}
	if err := p.exponent(); err != nil {
		return nil, err
	}

	end := p.index
	number_literal := string(p.src[begin:end])

	if end_integer == end {
		i64, err := strconv.ParseInt(number_literal, 10, 32)
		return int(i64), err
	} else {
		return strconv.ParseFloat(number_literal, 64)
	}
}

func (p *parser) integer() error {
	if p.current() == '-' {
		p.next()
	}

	if p.current() == '0' {
		p.next()
		return nil
	}

	return p.digits()
}

func (p *parser) digits() error {
	if !is_digit(p.current()) {
		msg := fmt.Sprintf("expected digit, but %q", p.current())
		return errors.New(msg)
	}

	p.next()

	for is_digit(p.current()) {
		p.next()
	}

	return nil
}

func (p *parser) fraction() error {
	if p.current() == '.' {
		p.next()
		return p.digits()
	}

	return nil
}

func (p *parser) exponent() error {
	cur := p.current()
	if cur == 'E' || cur == 'e' {
		p.next()
		p.sign()
		return p.digits()
	}

	return nil
}

func (p *parser) sign() {
	cur := p.current()
	if cur == '+' || cur == '-' {
		p.next()
	}
}

func (p *parser) skip_white_space() {
	for is_white_space(p.current()) {
		p.next()
	}
}

func (p *parser) expectRune(r rune) error {
	if p.current() == r {
		p.next()
		return nil
	}

	msg := fmt.Sprintf("expected %q, but %q", r, p.current())
	return errors.New(msg)
}

func (p *parser) consume(str string) bool {
	end := p.index + len(str)
	if end > p.srcLen {
		end = p.srcLen
	}

	actual := string(p.src[p.index:end])
	if actual == str {
		p.index = end
		return true
	}

	return false
}

func (p *parser) consumeRune(r rune) bool {
	if p.current() == r {
		p.next()
		return true
	}

	return false
}

func (p *parser) current() rune {
	if p.index >= p.srcLen {
		return 0
	}

	return p.src[p.index]
}

func (p *parser) next() {
	if p.index < p.srcLen {
		p.index += 1
	}
}

func is_white_space(r rune) bool {
	return r == ' ' || r == '\n' || r == '\r' || r == '\t'
}

func is_digit(r rune) bool {
	return '0' <= r && r <= '9'
}

func is_onenine(r rune) bool {
	return '1' <= r && r <= '9'
}

func isHex(r rune) bool {
	return is_digit(r) || ('a' <= r && r <= 'f') || ('A' <= r && r <= 'F')
}
