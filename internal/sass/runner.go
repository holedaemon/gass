// Package sass provides a Sass transpiler capable of transpiling multiple
// input files at once.
package sass

import (
	"bytes"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/bep/godartsass/v2"
	"github.com/holedaemon/gass/internal/gassfile"
	"github.com/holedaemon/gass/internal/logger"
)

type Runner struct {
	debug   bool
	binary  string
	timeout time.Duration

	l  *slog.Logger
	tp *godartsass.Transpiler
}

func NewRunner(opts ...Option) (*Runner, error) {
	r := new(Runner)
	for _, o := range opts {
		o(r)
	}

	if r.l == nil {
		r.l = logger.New(r.debug)
	}

	tpOpts := godartsass.Options{}
	if r.binary != "" {
		tpOpts.DartSassEmbeddedFilename = r.binary
	}

	if r.timeout != 0 {
		tpOpts.Timeout = r.timeout
	}

	tpOpts.LogEventHandler = r.logHandler

	tp, err := godartsass.Start(tpOpts)
	if err != nil {
		return nil, fmt.Errorf("%w: creating transpiler", err)
	}

	r.tp = tp
	return r, nil
}

func (r *Runner) logHandler(e godartsass.LogEvent) {
	switch e.Type {
	case godartsass.LogEventTypeDebug:
		r.l.Debug("debug from dart-sass", "message", e.Message)
	case godartsass.LogEventTypeDeprecated:
		r.l.Warn("deprecated sass feature", "message", e.Message)
	case godartsass.LogEventTypeWarning:
		r.l.Warn("warning from dart-sass", "message", e.Message)
	}
}

func (r *Runner) Close() error {
	return r.tp.Close()
}

func (r *Runner) Transpile(path string, opts ...TranspileOption) error {
	gs, err := gassfile.Load(path)
	if err != nil {
		return fmt.Errorf("loading gassfile: %w", err)
	}

	var buf bytes.Buffer

	r.l.Debug("iterating sources")
	for src := range gs.Sources() {
		input := src.Input()
		output := src.Output()

		to := new(transpileOptions)
		for _, o := range opts {
			o(to)
		}

		r.l.Debug("reading input file to string")
		inputFile, err := os.Open(input)
		if err != nil {
			return fmt.Errorf("opening input file: %w", err)
		}

		if _, err := inputFile.WriteTo(&buf); err != nil {
			return fmt.Errorf("copying to input file: %w", err)
		}

		inputFile.Close()

		r.l.Debug("converting options into arguments")
		args, err := to.Args(buf.String())
		if err != nil {
			return fmt.Errorf("creating args: %w", err)
		}

		r.l.Debug("transpiling input", "input", input, "output", output)
		res, err := r.tp.Execute(args)
		if err != nil {
			return fmt.Errorf("transpiling input: %w", err)
		}

		fmt.Println(output)
		r.l.Debug("writing output to file", "output", output)
		outputFile, err := os.Create(output)
		if err != nil {
			return fmt.Errorf("creating/truncating output file: %w", err)
		}

		if _, err := outputFile.WriteString(res.CSS); err != nil {
			return fmt.Errorf("writing CSS to output file: %w", err)
		}

		outputFile.Close()

		if to.sourceMaps {
			r.l.Debug("writing source map for output")
			outputMap := src.OutputMap()
			outputMapFile, err := os.Create(outputMap)
			if err != nil {
				return fmt.Errorf("creating/trunating output map file: %w", err)
			}

			if _, err := outputMapFile.WriteString(res.SourceMap); err != nil {
				return fmt.Errorf("writing source map to output file: %w", err)
			}

			outputMapFile.Close()
		}
	}

	return nil
}
