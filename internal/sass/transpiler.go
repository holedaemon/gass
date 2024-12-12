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

// Transpiler is a thin wrapper over a *[godartsass.Transpiler], capable of
// transpiling multiple Sass files from various sources.
type Transpiler struct {
	debug   bool
	binary  string
	timeout time.Duration

	l  *slog.Logger
	tp *godartsass.Transpiler
}

// New creates a new *[Transpiler] ready for use.
func New(opts ...Option) (*Transpiler, error) {
	t := new(Transpiler)
	for _, o := range opts {
		o(t)
	}

	if t.l == nil {
		t.l = logger.New(t.debug)
	}

	tpOpts := godartsass.Options{}
	if t.binary != "" {
		tpOpts.DartSassEmbeddedFilename = t.binary
	}

	if t.timeout != 0 {
		tpOpts.Timeout = t.timeout
	}

	tpOpts.LogEventHandler = t.logHandler

	tp, err := godartsass.Start(tpOpts)
	if err != nil {
		return nil, fmt.Errorf("%w: creating transpiler", err)
	}

	t.tp = tp
	return t, nil
}

// logHandler logs any messages emitted by dart-sass.
func (t *Transpiler) logHandler(e godartsass.LogEvent) {
	switch e.Type {
	case godartsass.LogEventTypeDebug:
		t.l.Debug("@debug", "message", e.Message)
	case godartsass.LogEventTypeDeprecated:
		t.l.Warn("deprecated sass feature", "message", e.Message)
	case godartsass.LogEventTypeWarning:
		t.l.Warn("@warn", "message", e.Message)
	}
}

// Close closes the underlying *[godartsass.Transpiler].
func (t *Transpiler) Close() error {
	return t.tp.Close()
}

// Transpile converts Scss source files into CSS. The path argument should be
// a valid *[gassfile.Gassfile], which will be used to load input sources.
// Consumers may pass options to configure the underlying
// *[godartsass.Transpiler].
func (t *Transpiler) Transpile(path string, opts ...TranspileOption) error {
	gs, err := gassfile.Load(path)
	if err != nil {
		return fmt.Errorf("loading gassfile: %w", err)
	}

	t.l.Debug("iterating sources")

	var buf bytes.Buffer
	for src := range gs.Sources() {
		input := src.Input()
		output := src.Output()

		t.l.Debug("preparing transpiler options")

		relative, err := src.Relative()
		if err != nil {
			return fmt.Errorf("collecting relative directories: %w", err)
		}

		baseDir := src.InputDir()
		ir := &importResolver{baseDir: baseDir}

		var syntax TranspileOption
		switch src.Syntax() {
		case godartsass.SourceSyntaxSASS:
			syntax = Sass()
		case godartsass.SourceSyntaxSCSS:
			syntax = SCSS()
		case godartsass.SourceSyntaxCSS:
			syntax = CSS()
		}

		t.l.Debug("reading input file to string")
		inputFile, err := os.Open(input)
		if err != nil {
			return fmt.Errorf("opening input file: %w", err)
		}

		if _, err := inputFile.WriteTo(&buf); err != nil {
			return fmt.Errorf("copying to input file: %w", err)
		}

		inputFile.Close()

		opts = append(opts,
			IncludePaths(relative...),
			ImportResolver(ir),
			Source(buf.String()),
			syntax,
		)

		to := new(transpileOptions)
		for _, o := range opts {
			o(to)
		}

		args, err := to.Args()
		if err != nil {
			return fmt.Errorf("creating args: %w", err)
		}

		t.l.Debug("transpiling input", "input", input, "output", output)
		res, err := t.tp.Execute(args)
		if err != nil {
			return fmt.Errorf("transpiling input: %w", err)
		}

		t.l.Debug("writing output to file", "output", output)
		outputFile, err := os.Create(output)
		if err != nil {
			return fmt.Errorf("creating/truncating output file: %w", err)
		}

		if _, err := outputFile.WriteString(res.CSS); err != nil {
			return fmt.Errorf("writing CSS to output file: %w", err)
		}

		outputFile.Close()

		if to.sourceMaps {
			t.l.Debug("writing source map for output")
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

	t.l.Debug("finished transpiling sources")
	return nil
}
