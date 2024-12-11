package sass

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/bep/godartsass/v2"
)

// ImportResolver is responsible for resolving imports present within an input
// file. It implements the [godartsass.ImportResolver] interface.
type ImportResolver struct {
	baseDir string
}

// CanonicalizeURL creates a canonical version of the given URL, assuming it
// can resolve it. Otherwise, it will return an empty string.
func (r *ImportResolver) CanonicalizeURL(url string) (string, error) {
	dir, file := filepath.Split(url)
	wholeDir := filepath.Join(r.baseDir, dir)

	entries, err := os.ReadDir(wholeDir)
	if err != nil {
		return "", err
	}

	hasFile := slices.ContainsFunc(entries, func(e os.DirEntry) bool {
		if e.IsDir() {
			return false
		}

		name := strings.TrimPrefix(e.Name(), "_")
		return strings.HasPrefix(name, file)
	})

	if hasFile {
		wholePath := filepath.Join(wholeDir, file)
		return "file://" + wholePath, nil
	}

	return "", nil
}

// Load loads the content of a canonical URL and returns it as a
// [godartsass.Import].
func (r *ImportResolver) Load(url string) (godartsass.Import, error) {
	fmt.Println(url)
	return godartsass.Import{}, nil
}
