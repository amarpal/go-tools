package sa1102

import (
	"go/ast"
	"go/types"

	"honnef.co/go/tools/analysis/lint"

	"golang.org/x/tools/go/analysis"
)

var SCAnalyzer = lint.InitializeAnalyzer(&lint.Analyzer{
	Analyzer: &analysis.Analyzer{
		Name:     "SA1102",
		Run:      run,
		Requires: []*analysis.Analyzer{},
	},
	Doc: &lint.RawDocumentation{
		Title:    "log.Fatalln should not be called in package not named `main`",
		Text:     `log.Fatalln calls os.Exit(1) which is unexpected in packages not named 'main'`,
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
			// Check if the function is 'Fatalln' and from the 'log' package.
			if fun.Sel.Name == "Fatalln" {
				pkgIdent, ok := fun.X.(*ast.Ident)
				if !ok {
					return true
				}
				obj := pass.TypesInfo.Uses[pkgIdent]
				pkgName, ok := obj.(*types.PkgName)
				if ok && pkgName.Imported().Path() == "log" {
					pass.Reportf(call.Pos(), "use of log.Fatalln in package not named 'main'")
				}
			}
			return true
		})
	}
	return nil, nil
}
