package delphi

import (
	"errors"
	"io/ioutil"
	"strings"
	"unicode"
)

func ParseFile(path string) (*File, error) {
	code, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, errors.New(
			"delphi.ParseFile: cannot read file '" + path + "': " + err.Error(),
		)
	}
	return ParseCode(code)
}

func ParseCode(code []byte) (*File, error) {
	code, err := makeUTF8(code)
	if err != nil {
		return nil, errors.New(
			"delphi.ParseCode: unknown file encoding: " + err.Error(),
		)
	}

	p := parser{
		tokens: newTokenizer([]rune(string(code))),
	}
	err = p.parse()
	if err != nil {
		return nil, err
	}
	return &p.file, nil
}

// TODO makeUTF8 should handle different file encodings and make them UTF-8.
func makeUTF8(code []byte) ([]byte, error) {
	return code, nil // for now assume it is UTF-8 already
}

type File struct {
	Type FileType
	Name string
}

type FileType string

const (
	Program FileType = "program"
	Unit    FileType = "unit"
	Package FileType = "package"
	Library FileType = "library"
)

type parser struct {
	tokens *tokenizer
	file   File
}

func (p *parser) parse() error {
	t := p.nextSolidToken()
	if isWord(t, "program") {
		p.file.Type = Program
	} else {
		return errors.New("TODO keyword 'program' expected")
	}

	t = p.nextSolidToken()
	if t.tokenType != tokenWord {
		return errors.New("name expected")
	}
	p.file.Name = t.text

	return nil
}

func (p *parser) nextSolidToken() token {
	t := p.tokens.next()
	for t.tokenType == tokenComment || t.tokenType == tokenWhiteSpace {
		t = p.tokens.next()
	}
	return t
}

func isWord(t token, word string) bool {
	return t.tokenType == tokenWord && strings.ToLower(t.text) == word
}

func newTokenizer(code []rune) *tokenizer {
	return &tokenizer{code: code}
}

type tokenizer struct {
	code []rune
	cur  int
}

func (t *tokenizer) next() token {
	haveType := tokenIllegal
	var text string
	start := t.cur

	r := t.currentRune()
	switch r {
	case 0:
		return token{tokenType: tokenEOF, text: "end of file"}
	default:
		if unicode.IsSpace(r) {
			for unicode.IsSpace(t.nextRune()) {
			}
			haveType = tokenWhiteSpace
		} else if r == '_' || unicode.IsLetter(r) {
			word := func(r rune) bool {
				return r == '_' || unicode.IsLetter(r) || unicode.IsDigit(r)
			}
			for word(t.nextRune()) {
			}
			haveType = tokenWord
		}
	}

	if haveType != tokenIllegal {
		text = string(t.code[start:t.cur])
	}
	return token{tokenType: haveType, text: text}
}

func (t *tokenizer) currentRune() rune {
	if t.cur >= len(t.code) {
		return 0
	}
	return t.code[t.cur]
}

func (t *tokenizer) nextRune() rune {
	if t.cur < len(t.code) {
		t.cur++
	}
	return t.currentRune()
}

type token struct {
	tokenType tokenType
	text      string
}

type tokenType rune

const (
	tokenIllegal tokenType = -1
	tokenEOF     tokenType = 0
	// Reserve ASCII characters 1..255 for their literal representations.
	tokenComment    tokenType = 256
	tokenWhiteSpace tokenType = 257
	tokenWord       tokenType = 258
)
