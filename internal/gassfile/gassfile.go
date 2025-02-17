// Package gassfile provides support for the 'gassfile': a file containing
// a list of SCSS sources to transpile to CSS.
package gassfile

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"iter"
	"os"
	"strings"
)

var (
	// ErrInvalidSource is returned when a [Source] is invalid.
	ErrInvalidSource = errors.New("gassfile: source is invalid")

	validExts = []string{".scss", ".sass", ".css"}
)

// Gassfile is a plaintext file containing a collecton of *[Source], delimited
// by newlines.
type Gassfile struct {
	sources []*Source
}

// Sources returns an iterator over g's list of *[Source].
func (g *Gassfile) Sources() iter.Seq[*Source] {
	return func(yield func(*Source) bool) {
		for _, s := range g.sources {
			if !yield(s) {
				return
			}
		}
	}
}

// NewFromReader creates a new *[Gassfile] from the given [io.Reader].
func NewFromReader(r io.Reader) (*Gassfile, error) {
	sources, err := scan(r)
	if err != nil {
		return nil, fmt.Errorf("collecting sources: %w", err)
	}

	return &Gassfile{
		sources: sources,
	}, nil
}

// New creates a new *[Gassfile] from the given path.
func New(path string) (*Gassfile, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("opening file: %w", err)
	}

	defer file.Close()

	return NewFromReader(file)
}

// scan creates a scanner for the input r and iterates over each line of the
// file. If a line starts with '#', it's accepted as a comment, otherwise the
// line is split by ' ' and the resulting substrings are used as IO for a
// *[Source].
func scan(r io.Reader) ([]*Source, error) {
	scanner := bufio.NewScanner(r)
	sources := make([]*Source, 0)

	for scanner.Scan() {
		text := scanner.Text()
		paths := strings.Split(text, " ")

		if strings.HasPrefix(text, "#") {
			continue
		}

		if len(paths) != 2 {
			return nil, fmt.Errorf("%w: must be 2 paths", ErrInvalidSource)
		}

		input := paths[0]
		output := paths[1]
		src, err := NewSource(input, output)
		if err != nil {
			return nil, err
		}

		sources = append(sources, src)
	}

	return sources, nil
}
