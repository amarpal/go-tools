package sa1007

import (
	"fmt"
	"go/constant"
	"net/url"

	"github.com/amarpal/go-tools/analysis/callcheck"
	"github.com/amarpal/go-tools/analysis/lint"
	"github.com/amarpal/go-tools/internal/passes/buildir"
	"github.com/amarpal/go-tools/knowledge"

	"golang.org/x/tools/go/analysis"
)

var SCAnalyzer = lint.InitializeAnalyzer(&lint.Analyzer{
	Analyzer: &analysis.Analyzer{
		Name:     "SA1007",
		Requires: []*analysis.Analyzer{buildir.Analyzer},
		Run:      callcheck.Analyzer(rules),
	},
	Doc: &lint.RawDocumentation{
		Title:    `Invalid URL in \'net/url.Parse\'`,
		Since:    "2017.1",
		Severity: lint.SeverityError,
		MergeIf:  lint.MergeIfAny,
	},
})

var Analyzer = SCAnalyzer.Analyzer

var rules = map[string]callcheck.Check{
	"net/url.Parse": func(call *callcheck.Call) {
		arg := call.Args[knowledge.Arg("net/url.Parse.rawurl")]
		if c := callcheck.ExtractConstExpectKind(arg.Value, constant.String); c != nil {
			s := constant.StringVal(c.Value)
			_, err := url.Parse(s)
			if err != nil {
				arg.Invalid(fmt.Sprintf("%q is not a valid URL: %s", s, err))
			}
		}
	},
}
