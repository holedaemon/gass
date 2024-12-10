package sass

import (
	"errors"

	"github.com/bep/godartsass/v2"
)

type transpileOptions struct {
	syntax       godartsass.SourceSyntax
	style        godartsass.OutputStyle
	sourceMaps   bool
	embedSources bool
}

func (to *transpileOptions) Args(src string) (godartsass.Args, error) {
	args := godartsass.Args{
		Source:                  src,
		SourceSyntax:            to.syntax,
		OutputStyle:             to.style,
		EnableSourceMap:         to.sourceMaps,
		SourceMapIncludeSources: to.embedSources,
	}

	if args.Source == "" {
		return args, errors.New("sass: source is blank")
	}

	return args, nil
}

type TranspileOption func(*transpileOptions)

func Compressed() TranspileOption {
	return func(to *transpileOptions) {
		to.style = godartsass.OutputStyleCompressed
	}
}

func Expanded() TranspileOption {
	return func(to *transpileOptions) {
		to.style = godartsass.OutputStyleExpanded
	}
}

func SCSS() TranspileOption {
	return func(to *transpileOptions) {
		to.syntax = godartsass.SourceSyntaxSCSS
	}
}

func Sass() TranspileOption {
	return func(to *transpileOptions) {
		to.syntax = godartsass.SourceSyntaxSASS
	}
}

func CSS() TranspileOption {
	return func(to *transpileOptions) {
		to.syntax = godartsass.SourceSyntaxCSS
	}
}

func SourceMaps() TranspileOption {
	return func(to *transpileOptions) {
		to.sourceMaps = true
	}
}

func EmbedSources() TranspileOption {
	return func(to *transpileOptions) {
		to.embedSources = true
	}
}
