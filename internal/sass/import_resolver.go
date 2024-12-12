package sass

import (
	"bytes"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/bep/godartsass/v2"
	"github.com/holedaemon/gass/internal/gassfile"
)

// importResolver is responsible for resolving imports present within an input
// file. It implements the [godartsass.importResolver] interface.
type importResolver struct {
	baseDir string
}

// CanonicalizeURL creates a canonical version of the given URL, assuming it
// can resolve it. Otherwise, it will return an empty string.
func (r *importResolver) CanonicalizeURL(url string) (string, error) {
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
		ext := filepath.Ext(file)
		if ext == "" {
			// Assume the input is using Scss if there is no explicit extension.
			file += ".scss"
		}

		wholePath := filepath.Join(wholeDir, file)
		return "file://" + wholePath, nil
	}

	return "", nil
}

// Load loads the content of a canonical URL and returns it as a
// [godartsass.Import].
func (r *importResolver) Load(url string) (godartsass.Import, error) {
	clean := strings.TrimPrefix(url, "file://")
	syntax := gassfile.DetermineSourceSyntax(filepath.Ext(clean))
	imp := godartsass.Import{
		SourceSyntax: syntax,
	}

	file, err := os.Open(clean)
	if err != nil {
		return godartsass.Import{}, err
	}

	defer file.Close()

	var buf bytes.Buffer
	if _, err := file.WriteTo(&buf); err != nil {
		return godartsass.Import{}, err
	}

	imp.Content = buf.String()
	return imp, nil
}
