// staticcheck analyses Go code and makes it better.
package main

import (
	"log"
	"os"

	"github.com/amarpal/go-tools/lintcmd"
	"github.com/amarpal/go-tools/lintcmd/version"
	"github.com/amarpal/go-tools/quickfix"
	"github.com/amarpal/go-tools/simple"
	"github.com/amarpal/go-tools/staticcheck"
	"github.com/amarpal/go-tools/stylecheck"
	"github.com/amarpal/go-tools/unused"
)

func main() {
	cmd := lintcmd.NewCommand("staticcheck")
	cmd.SetVersion(version.Version, version.MachineVersion)

	fs := cmd.FlagSet()
	debug := fs.String("debug.unused-graph", "", "Write unused's object graph to `file`")
	qf := fs.Bool("debug.run-quickfix-analyzers", false, "Run quickfix analyzers")

	cmd.ParseFlags(os.Args[1:])

	cmd.AddAnalyzers(simple.Analyzers...)
	cmd.AddAnalyzers(staticcheck.Analyzers...)
	cmd.AddAnalyzers(stylecheck.Analyzers...)
	cmd.AddAnalyzers(unused.Analyzer)

	if *qf {
		cmd.AddAnalyzers(quickfix.Analyzers...)
	}

	if *debug != "" {
		f, err := os.OpenFile(*debug, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
		if err != nil {
			log.Fatal(err)
		}
		unused.Debug = f
	}

	cmd.Run()
}
