package st1023

import (
	"github.com/amarpal/go-tools/analysis/lint"
	"github.com/amarpal/go-tools/internal/sharedcheck"
)

func init() {
	SCAnalyzer.Analyzer.Name = "ST1023"
}

var SCAnalyzer = lint.InitializeAnalyzer(&lint.Analyzer{
	Analyzer: sharedcheck.RedundantTypeInDeclarationChecker("should", false),
	Doc: &lint.RawDocumentation{
		Title:      "Redundant type in variable declaration",
		Since:      "2021.1",
		NonDefault: true,
		MergeIf:    lint.MergeIfAll,
	},
})

var Analyzer = SCAnalyzer.Analyzer
