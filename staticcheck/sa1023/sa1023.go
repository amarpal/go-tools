package sa1023

import (
	"go/types"

	"github.com/amarpal/go-tools/analysis/lint"
	"github.com/amarpal/go-tools/analysis/report"
	"github.com/amarpal/go-tools/go/ir"
	"github.com/amarpal/go-tools/go/ir/irutil"
	"github.com/amarpal/go-tools/internal/passes/buildir"
	"github.com/amarpal/go-tools/knowledge"

	"golang.org/x/tools/go/analysis"
)

var SCAnalyzer = lint.InitializeAnalyzer(&lint.Analyzer{
	Analyzer: &analysis.Analyzer{
		Name:     "SA1023",
		Run:      run,
		Requires: []*analysis.Analyzer{buildir.Analyzer},
	},
	Doc: &lint.RawDocumentation{
		Title:    `Modifying the buffer in an \'io.Writer\' implementation`,
		Text:     `\'Write\' must not modify the slice data, even temporarily.`,
		Since:    "2017.1",
		Severity: lint.SeverityError,
		MergeIf:  lint.MergeIfAny,
	},
})

var Analyzer = SCAnalyzer.Analyzer

func run(pass *analysis.Pass) (interface{}, error) {
	// TODO(dh): this might be a good candidate for taint analysis.
	// Taint the argument as MUST_NOT_MODIFY, then propagate that
	// through functions like bytes.Split

	for _, fn := range pass.ResultOf[buildir.Analyzer].(*buildir.IR).SrcFuncs {
		sig := fn.Signature
		if fn.Name() != "Write" || sig.Recv() == nil {
			continue
		}
		if !types.Identical(sig, knowledge.Signatures["(io.Writer).Write"]) {
			continue
		}

		for _, block := range fn.Blocks {
			for _, ins := range block.Instrs {
				switch ins := ins.(type) {
				case *ir.Store:
					addr, ok := ins.Addr.(*ir.IndexAddr)
					if !ok {
						continue
					}
					if addr.X != fn.Params[1] {
						continue
					}
					report.Report(pass, ins, "io.Writer.Write must not modify the provided buffer, not even temporarily")
				case *ir.Call:
					if !irutil.IsCallTo(ins.Common(), "append") {
						continue
					}
					if ins.Common().Args[0] != fn.Params[1] {
						continue
					}
					report.Report(pass, ins, "io.Writer.Write must not modify the provided buffer, not even temporarily")
				}
			}
		}
	}
	return nil, nil
}
