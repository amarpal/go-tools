package sa1103

import (
	"go/ast"
	"go/types"

	"honnef.co/go/tools/analysis/lint"

	"golang.org/x/tools/go/analysis"
)

var SCAnalyzer = lint.InitializeAnalyzer(&lint.Analyzer{
	Analyzer: &analysis.Analyzer{
		Name:     "SA1103",
		Run:      run,
		Requires: []*analysis.Analyzer{},
	},
	Doc: &lint.RawDocumentation{
		Title:    "os.Exit should not be called in package not named `main`",
		Text:     `os.Exit is unexpected in packages not named 'main'`,
		Since:    "Unreleased",
		Severity: lint.SeverityWarning,
	},
})

var Analyzer = SCAnalyzer.Analyzer

func run(pass *analysis.Pass) (interface{}, error) {
	// Check if the package name is 'main'. If it is, then do nothing.
	if pass.Pkg.Name() == "main" {
		return nil, nil
	}
	// Inspect the nodes in the AST.
	for _, file := range pass.Files {
		ast.Inspect(file, func(n ast.Node) bool {
			// Look for call expressions (e.g., function calls).
			call, ok := n.(*ast.CallExpr)
			if !ok {
				return true // continue to next node
			}
			// Check if the call is a function call.
			fun, ok := call.Fun.(*ast.SelectorExpr)
			if !ok {
				return true
			}
			// Check if the function is 'Exit' and from the 'os' package.
			if fun.Sel.Name == "Exit" {
				pkgIdent, ok := fun.X.(*ast.Ident)
				if !ok {
					return true
				}
				obj := pass.TypesInfo.Uses[pkgIdent]
				pkgName, ok := obj.(*types.PkgName)
				if ok && pkgName.Imported().Path() == "os" {
					pass.Reportf(call.Pos(), "use of os.Exit in package not named 'main'")
				}
			}
			return true
		})
	}
	return nil, nil
}
