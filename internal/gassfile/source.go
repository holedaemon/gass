package gassfile

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"
)

// Source represents a collection of paths used to transpile Sass to CSS.
type Source struct {
	input  string
	output string
}

// NewSource creates a new Source from the given inputs and output.
// Output should be an absolute path to a directory in which transpiled
// CSS files are output.
// Inputs should be one or more absolute paths to a Sass source file. Globs are
// accepted.
func NewSource(input, output string) (*Source, error) {
	if input == "" {
		return nil, errors.New("gassfile: source: input is blank")
	}

	if output == "" {
		return nil, errors.New("gassfile: source: output is blank")
	}

	if !filepath.IsAbs(input) {
		return nil, fmt.Errorf("gassfile: source: %s is not an absolute path", input)
	}

	if !filepath.IsAbs(output) {
		return nil, fmt.Errorf("gassfile: source: %s is not an absolute path", output)
	}

	ext := filepath.Ext(input)
	if !stringInSlice(ext, validExts) {
		return nil, fmt.Errorf("gassfile: source: %s ends with an unrecognized extension", input)
	}

	return &Source{
		input:  input,
		output: output,
	}, nil
}

// Input returns the input path of the Source.
func (s *Source) Input() string {
	return s.input
}

func (s *Source) outputName() string {
	outputExt := filepath.Ext(s.output)

	if strings.EqualFold(outputExt, ".css") {
		return s.output
	}

	inputBase := filepath.Base(s.input)
	inputExt := filepath.Ext(s.input)
	name := fmt.Sprintf("%s.css", strings.TrimSuffix(inputBase, inputExt))
	return filepath.Join(s.output, name)
}

// Output returns the output of the Source.
// If the output ends with a CSS file, that is returned, otherwise the name of
// the input file is used.
func (s *Source) Output() string {
	return s.outputName()
}

// OutputMap returns the map file of the Source's output.
func (s *Source) OutputMap() string {
	return s.outputName() + ".map"
}
