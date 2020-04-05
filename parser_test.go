package delphi_test

import (
	"github.com/gonutz/check"
	"github.com/gonutz/delphi"
	"testing"
)

func TestParseEmptyProgram(t *testing.T) {
	ast, err := delphi.ParseFile("testdata/empty.dpr")
	check.Eq(t, err, nil)
	check.Eq(t, ast, &delphi.File{
		Type: delphi.Program,
		Name: "Empty",
	})
}
