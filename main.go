package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/holedaemon/gass/internal/logger"
	"github.com/holedaemon/gass/internal/sass"
	"github.com/holedaemon/gass/internal/version"
)

func main() {
	versionFlag := flag.Bool("v", false, "Display gass' version")
	gassfilePath := flag.String("f", ".gassfile", "Path to Gassfile to use for transpilation")
	debug := flag.Bool("d", false, "Run in debug mode?")
	compressed := flag.Bool("c", false, "Compress CSS output?")
	sourceMaps := flag.Bool("m", true, "Generate source maps?")
	embedSources := flag.Bool("e", false, "Embed sources into source maps?")
	binary := flag.String("b", "", "Path to dart-sass binary to use")
	timeout := flag.Duration("t", time.Second*30, "Maximum length of time alloted for transpilation")

	flag.Parse()

	if *versionFlag {
		fmt.Fprintf(os.Stdout, "gass v%s\n", version.Version)
		os.Exit(0)
	}

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

	l.Debug("preparing to transpile")

	switch *gassfilePath {
	case "-":
		err = r.TranspileFromReader(os.Stdin, tpOpts...)
	default:
		err = r.Transpile(*gassfilePath, tpOpts...)
	}

	if err != nil {
		l.Error("error transpiling sources", slog.Any("error", err.Error()))
		os.Exit(1)
	}
}
