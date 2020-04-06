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

func TestParseEmptyProgramAsCode(t *testing.T) {
	ast, err := delphi.ParseCode("empty.dpr", []byte(`
program Empty;
begin
end.
`))
	check.Eq(t, err, nil)
	check.Eq(t, ast, &delphi.File{
		Type: delphi.Program,
		Name: "Empty",
	})
}

func parseError(t *testing.T, path, code, wantErrorMsg string) {
	t.Helper()
	_, err := delphi.ParseCode(path, []byte(code))
	if err == nil {
		t.Fatal("error expected")
	}
	check.Eq(t, err.Error(), wantErrorMsg)
}

func TestErrorsInEmptyProgram(t *testing.T) {
	parseError(t, "empty.dpr", `
Empty;
begin
end.
`,
		"DPR file must start with 'program' or 'library' keyword")

	parseError(t, "empty.dpr", `
program ;
begin
end.
	`,
		"missing program name")

	parseError(t, "empty.dpr", `
program Empty
begin
end.
	`,
		"missing ';' after program name, found 'begin'")

	parseError(t, "empty.dpr", `
program Empty;
end.
	`,
		"missing 'begin' at program start, found 'end'")

	parseError(t, "empty.dpr", `
program Empty;
begin
.
	`,
		"missing 'end' at end of program, found '.'")

	parseError(t, "empty.dpr", `
program Empty;
begin
end
`,
		"missing '.' at end of program, found end of file")
}
