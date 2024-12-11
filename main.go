package main

import (
	"flag"
	"log/slog"
	"os"
	"time"

	"github.com/bep/godartsass/v2"
	"github.com/holedaemon/gass/internal/logger"
	"github.com/holedaemon/gass/internal/sass"
)

func main() {
	gassfilePath := flag.String("f", ".gassfile", "Path to Gassfile to use for transpilation")
	debug := flag.Bool("d", false, "Run in debug mode?")
	compressed := flag.Bool("c", false, "Compress CSS output?")
	sourceMaps := flag.Bool("m", true, "Generate source maps?")
	embedSources := flag.Bool("e", false, "Embed sources into source maps?")
	syntaxOpt := flag.String("s", "SCSS", "The type of syntax used by inputs")
	binary := flag.String("b", "", "Path to dart-sass binary to use")
	timeout := flag.Duration("t", time.Second*30, "Maximum length of time alloted for transpilation")

	flag.Parse()

	l := logger.New(*debug)
	opts := make([]sass.Option, 0)
	if *debug {
		opts = append(opts, sass.Debug())
	}

	opts = append(opts,
		sass.SassBinary(*binary),
		sass.Logger(l),
		sass.Timeout(*timeout),
	)

	r, err := sass.New(opts...)
	if err != nil {
		l.Error("error creating new transpiler", slog.Any("error", err.Error()))
		os.Exit(1)
	}

	l.Debug("created new transpiler")

	tpOpts := make([]sass.TranspileOption, 0)
	if *compressed {
		tpOpts = append(tpOpts, sass.Compressed())
	} else {
		tpOpts = append(tpOpts, sass.Expanded())
	}

	if *sourceMaps {
		tpOpts = append(tpOpts, sass.SourceMaps())
	}

	if *embedSources {
		tpOpts = append(tpOpts, sass.EmbedSources())
	}

	syntax := godartsass.ParseSourceSyntax(*syntaxOpt)
	switch syntax {
	case godartsass.SourceSyntaxSASS:
		tpOpts = append(tpOpts, sass.Sass())
	case godartsass.SourceSyntaxSCSS:
		tpOpts = append(tpOpts, sass.SCSS())
	case godartsass.SourceSyntaxCSS:
		tpOpts = append(tpOpts, sass.CSS())
	}

	l.Debug("preparing to transpile")
	if err := r.Transpile(*gassfilePath, tpOpts...); err != nil {
		l.Error("error transpiling sources", slog.Any("error", err.Error()))
		os.Exit(1)
	}
}
