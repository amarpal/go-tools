package sa6006

import (
	"go/ast"

	"github.com/amarpal/go-tools/analysis/code"
	"github.com/amarpal/go-tools/analysis/lint"
	"github.com/amarpal/go-tools/analysis/report"
	"github.com/amarpal/go-tools/pattern"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
)

var SCAnalyzer = lint.InitializeAnalyzer(&lint.Analyzer{
	Analyzer: &analysis.Analyzer{
		Name:     "SA6006",
		Run:      run,
		Requires: []*analysis.Analyzer{inspect.Analyzer},
	},
	Doc: &lint.RawDocumentation{
		Title: `Using io.WriteString to write \'[]byte\'`,
		Text: `Using io.WriteString to write a slice of bytes, as in

    io.WriteString(w, string(b))

is both unnecessary and inefficient. Converting from \'[]byte\' to \'string\'
has to allocate and copy the data, and we could simply use \'w.Write(b)\'
instead.`,

		Since: "2024.1",
	},
})

var Analyzer = SCAnalyzer.Analyzer

var ioWriteStringConversion = pattern.MustParse(`(CallExpr (Symbol "io.WriteString") [_ (CallExpr (Builtin "string") [arg])])`)

func run(pass *analysis.Pass) (interface{}, error) {
	fn := func(node ast.Node) {
		m, ok := code.Match(pass, ioWriteStringConversion, node)
		if !ok {
			return
		}
		if !code.IsOfStringConvertibleByteSlice(pass, m.State["arg"].(ast.Expr)) {
			return
		}
		report.Report(pass, node, "use io.Writer.Write instead of converting from []byte to string to use io.WriteString")
	}
	code.Preorder(pass, fn, (*ast.CallExpr)(nil))

	return nil, nil
}
