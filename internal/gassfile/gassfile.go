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

type Gassfile struct {
	sources []*Source
}

// Load opens the file at the given path and attempts to load it as a Gassfile.
func Load(path string) (*Gassfile, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("%w: opening file", err)
	}

	defer file.Close()

	sources, err := scan(file)
	if err != nil {
		return nil, fmt.Errorf("%w: collecting sources", err)
	}

	return &Gassfile{
		sources: sources,
	}, nil
}

// Sources returns an iterator over the Gassfile's sources.
func (g *Gassfile) Sources() iter.Seq[*Source] {
	return func(yield func(*Source) bool) {
		for _, s := range g.sources {
			if !yield(s) {
				return
			}
		}
	}
}

// scan creates a scanner for the input r and iterates over each line of the
// file. If a line starts with '#', it's accepted as a comment, otherwise the
// line is split by ' ' and the resulting substrings are used as IO for a
// [Source].
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
