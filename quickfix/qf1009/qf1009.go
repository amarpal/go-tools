package qf1009

import (
	"go/ast"
	"go/token"

	"github.com/amarpal/go-tools/analysis/code"
	"github.com/amarpal/go-tools/analysis/edit"
	"github.com/amarpal/go-tools/analysis/lint"
	"github.com/amarpal/go-tools/analysis/report"
	"github.com/amarpal/go-tools/pattern"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
)

var SCAnalyzer = lint.InitializeAnalyzer(&lint.Analyzer{
	Analyzer: &analysis.Analyzer{
		Name:     "QF1009",
		Run:      run,
		Requires: []*analysis.Analyzer{inspect.Analyzer},
	},
	Doc: &lint.Documentation{
		Title:    `Use \'time.Time.Equal\' instead of \'==\' operator`,
		Since:    "2021.1",
		Severity: lint.SeverityInfo,
	},
})

var Analyzer = SCAnalyzer.Analyzer

var timeEqualR = pattern.MustParse(`(CallExpr (SelectorExpr lhs (Ident "Equal")) rhs)`)

func run(pass *analysis.Pass) (interface{}, error) {
	// FIXME(dh): create proper suggested fix for renamed import

	fn := func(node ast.Node) {
		expr := node.(*ast.BinaryExpr)
		if expr.Op != token.EQL {
			return
		}
		if !code.IsOfType(pass, expr.X, "time.Time") || !code.IsOfType(pass, expr.Y, "time.Time") {
			return
		}
		report.Report(pass, node, "probably want to use time.Time.Equal instead",
			report.Fixes(edit.Fix("Use time.Time.Equal method",
				edit.ReplaceWithPattern(pass.Fset, node, timeEqualR, pattern.State{"lhs": expr.X, "rhs": expr.Y}))))
	}
	code.Preorder(pass, fn, (*ast.BinaryExpr)(nil))
	return nil, nil
}
