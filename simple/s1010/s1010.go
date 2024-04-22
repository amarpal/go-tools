package s1010

import (
	"go/ast"

	"github.com/amarpal/go-tools/analysis/code"
	"github.com/amarpal/go-tools/analysis/edit"
	"github.com/amarpal/go-tools/analysis/facts/generated"
	"github.com/amarpal/go-tools/analysis/lint"
	"github.com/amarpal/go-tools/analysis/report"
	"github.com/amarpal/go-tools/pattern"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
)

var SCAnalyzer = lint.InitializeAnalyzer(&lint.Analyzer{
	Analyzer: &analysis.Analyzer{
		Name:     "S1010",
		Run:      run,
		Requires: []*analysis.Analyzer{inspect.Analyzer, generated.Analyzer},
	},
	Doc: &lint.RawDocumentation{
		Title: `Omit default slice index`,
		Text: `When slicing, the second index defaults to the length of the value,
making \'s[n:len(s)]\' and \'s[n:]\' equivalent.`,
		Since:   "2017.1",
		MergeIf: lint.MergeIfAny,
	},
})

var Analyzer = SCAnalyzer.Analyzer

var checkSlicingQ = pattern.MustParse(`(SliceExpr x@(Object _) low (CallExpr (Builtin "len") [x]) nil)`)

func run(pass *analysis.Pass) (interface{}, error) {
	fn := func(node ast.Node) {
		if _, ok := code.Match(pass, checkSlicingQ, node); ok {
			expr := node.(*ast.SliceExpr)
			report.Report(pass, expr.High,
				"should omit second index in slice, s[a:len(s)] is identical to s[a:]",
				report.FilterGenerated(),
				report.Fixes(edit.Fix("simplify slice expression", edit.Delete(expr.High))))
		}
	}
	code.Preorder(pass, fn, (*ast.SliceExpr)(nil))
	return nil, nil
}
