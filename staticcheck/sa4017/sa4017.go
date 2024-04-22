package sa4017

import (
	"fmt"
	"go/types"

	"github.com/amarpal/go-tools/analysis/code"
	"github.com/amarpal/go-tools/analysis/facts/purity"
	"github.com/amarpal/go-tools/analysis/lint"
	"github.com/amarpal/go-tools/analysis/report"
	"github.com/amarpal/go-tools/go/ir"
	"github.com/amarpal/go-tools/go/ir/irutil"
	"github.com/amarpal/go-tools/go/types/typeutil"
	"github.com/amarpal/go-tools/internal/passes/buildir"

	"golang.org/x/tools/go/analysis"
)

var SCAnalyzer = lint.InitializeAnalyzer(&lint.Analyzer{
	Analyzer: &analysis.Analyzer{
		Name:     "SA4017",
		Run:      run,
		Requires: []*analysis.Analyzer{buildir.Analyzer, purity.Analyzer},
	},
	Doc: &lint.Documentation{
		Title:    `Discarding the return values of a function without side effects, making the call pointless`,
		Since:    "2017.1",
		Severity: lint.SeverityWarning,
		MergeIf:  lint.MergeIfAll,
	},
})

var Analyzer = SCAnalyzer.Analyzer

func run(pass *analysis.Pass) (interface{}, error) {
	pure := pass.ResultOf[purity.Analyzer].(purity.Result)

fnLoop:
	for _, fn := range pass.ResultOf[buildir.Analyzer].(*buildir.IR).SrcFuncs {
		if code.IsInTest(pass, fn) {
			params := fn.Signature.Params()
			for i := 0; i < params.Len(); i++ {
				param := params.At(i)
				if typeutil.IsType(param.Type(), "*testing.B") {
					// Ignore discarded pure functions in code related
					// to benchmarks. Instead of matching BenchmarkFoo
					// functions, we match any function accepting a
					// *testing.B. Benchmarks sometimes call generic
					// functions for doing the actual work, and
					// checking for the parameter is a lot easier and
					// faster than analyzing call trees.
					continue fnLoop
				}
			}
		}

		for _, b := range fn.Blocks {
			for _, ins := range b.Instrs {
				ins, ok := ins.(*ir.Call)
				if !ok {
					continue
				}
				refs := ins.Referrers()
				if refs == nil || len(irutil.FilterDebug(*refs)) > 0 {
					continue
				}

				callee := ins.Common().StaticCallee()
				if callee == nil {
					continue
				}
				if callee.Object() == nil {
					// TODO(dh): support anonymous functions
					continue
				}
				if _, ok := pure[callee.Object().(*types.Func)]; ok {
					if pass.Pkg.Path() == "fmt_test" && callee.Object().(*types.Func).FullName() == "fmt.Sprintf" {
						// special case for benchmarks in the fmt package
						continue
					}
					report.Report(pass, ins, fmt.Sprintf("%s doesn't have side effects and its return value is ignored", callee.Object().Name()))
				}
			}
		}
	}
	return nil, nil
}
