package delphi

import (
	"errors"
	"io/ioutil"
	"path/filepath"
	"strconv"
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
	return ParseCode(path, code)
}

func ParseCode(path string, code []byte) (*File, error) {
	code, err := makeUTF8(code)
	if err != nil {
		return nil, errors.New(
			"delphi.ParseCode: unknown file encoding: " + err.Error(),
		)
	}

	p := parser{
		filePath: path,
		tokens:   newTokenizer([]rune(string(code))),
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

func (t FileType) String() string { return string(t) }

const (
	Program FileType = "program" // .dpr files
	Library FileType = "library" // .dpr files as well
	Unit    FileType = "unit"    // .pas files
	Package FileType = "package" // .dpk files
)

type parser struct {
	filePath string
	tokens   *tokenizer
	file     File
}

func (p *parser) parse() error {
	t := p.nextSolidToken()
	if isWord(t, "program") {
		p.file.Type = Program
	} else if isWord(t, "library") {
		return errors.New("TODO parse library")
	} else if isWord(t, "unit") {
		return errors.New("TODO parse unit")
	} else if isWord(t, "package") {
		return errors.New("TODO parse package")
	} else {
		ext := strings.TrimPrefix(strings.ToLower(filepath.Ext(p.filePath)), ".")
		if ext == "dpr" {
			return errors.New("DPR file must start with 'program' or 'library' keyword")
		}
		return errors.New("Delphi files must start with one of these keywords: 'program', 'library', 'unit', 'package'")
	}

	t = p.nextSolidToken()
	if t.tokenType != tokenWord {
		return errors.New("missing " + p.file.Type.String() + " name")
	}
	p.file.Name = t.text

	t = p.nextSolidToken()
	if t.tokenType != ';' {
		return errors.New("missing ';' after " + p.file.Type.String() + " name, found " + t.String())
	}

	t = p.nextSolidToken()
	if !isWord(t, "begin") {
		return errors.New("missing 'begin' at " + p.file.Type.String() + " start, found " + t.String())
	}

	t = p.nextSolidToken()
	if !isWord(t, "end") {
		return errors.New("missing 'end' at end of " + p.file.Type.String() + ", found " + t.String())
	}

	t = p.nextSolidToken()
	if t.tokenType != '.' {
		return errors.New("missing '.' at end of " + p.file.Type.String() + ", found " + t.String())
	}

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
	case ';', '.':
		haveType = tokenType(r)
		t.nextRune()
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
		} else {
			text = "TODO handle token start " + string(r)
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

func (t token) String() string {
	if 1 <= t.tokenType && t.tokenType <= 127 {
		return "'" + string(t.tokenType) + "'"
	}

	switch t.tokenType {
	case tokenIllegal:
		return "illegal token (" + t.text + ")"
	case tokenEOF:
		return "end of file"
	case tokenComment:
		return "comment"
	case tokenWhiteSpace:
		return "white space"
	case tokenWord:
		return "'" + t.text + "'"
	default:
		return "TODO token.String for " + strconv.Itoa(int(t.tokenType))
	}
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
