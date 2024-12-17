package gassfile

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/bep/godartsass/v2"
)

// Source represents a collection of paths used to transpile Sass to CSS.
type Source struct {
	input  string
	output string
}

func NewSource(input, output string) (*Source, error) {
	if input == "" {
		return nil, errors.New("gassfile: source: input is blank")
	}

	if output == "" {
		return nil, errors.New("gassfile: source: output is blank")
	}

	var err error
	input, err = resolveSource(input)
	if err != nil {
		return nil, fmt.Errorf("resolving input: %w", err)
	}

	output, err = resolveSource(output)
	if err != nil {
		return nil, fmt.Errorf("resolving output: %w", err)
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

// InputDir returns the directory of the Source's input.
func (s *Source) InputDir() string {
	return filepath.Dir(s.input)
}

// Syntax returns the syntax used by the input based on its file extension.
func (s *Source) Syntax() godartsass.SourceSyntax {
	return DetermineSourceSyntax(filepath.Ext(s.input))
}

// Relative returns a slice containing absolute paths to directories relative
// to the input of the Source.
func (s *Source) Relative() ([]string, error) {
	dir := filepath.Dir(s.input)
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	rd := make([]string, 0)
	for _, e := range entries {
		if e.IsDir() {
			rd = append(rd, filepath.Join(dir, e.Name()))
		}
	}

	return rd, nil
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
