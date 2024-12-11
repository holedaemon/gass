package sass

import (
	"errors"

	"github.com/bep/godartsass/v2"
)

// transpileOptions are used to configure a *[godartsass.Transpiler].
type transpileOptions struct {
	syntax       godartsass.SourceSyntax
	style        godartsass.OutputStyle
	sourceMaps   bool
	embedSources bool
}

// Args converts specified options into a [godartsass.Args].
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

// TranspileOption are used to configure a *[godartsass.Transpiler].
type TranspileOption func(*transpileOptions)

// Compressed configures a *[godartsass.Transpiler] to output minified CSS.
func Compressed() TranspileOption {
	return func(to *transpileOptions) {
		to.style = godartsass.OutputStyleCompressed
	}
}

// Expanded configures a *[godartsass.Transpiler] to output unminified CSS.
func Expanded() TranspileOption {
	return func(to *transpileOptions) {
		to.style = godartsass.OutputStyleExpanded
	}
}

// SCSS configures a *[godartsass.Transpiler] to expect input files to use
// Scss syntax.
func SCSS() TranspileOption {
	return func(to *transpileOptions) {
		to.syntax = godartsass.SourceSyntaxSCSS
	}
}

// Sass configures a *[godartsass.Transpiler] to expect input files to use
// Sass syntax.
func Sass() TranspileOption {
	return func(to *transpileOptions) {
		to.syntax = godartsass.SourceSyntaxSASS
	}
}

// CSS configures a *[godartsass.Transpiler] to expect input files to use
// CSS syntax.
func CSS() TranspileOption {
	return func(to *transpileOptions) {
		to.syntax = godartsass.SourceSyntaxCSS
	}
}

// SourceMaps configures a *[godartsass.Transpiler] to output source maps
// when transpiling.
func SourceMaps() TranspileOption {
	return func(to *transpileOptions) {
		to.sourceMaps = true
	}
}

// EmbedSources configures a *[godartsass.Transpiler] to embed sources into
// output source maps.
func EmbedSources() TranspileOption {
	return func(to *transpileOptions) {
		to.embedSources = true
	}
}
