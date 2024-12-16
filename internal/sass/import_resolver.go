package sass

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"

	"github.com/bep/godartsass/v2"
	"github.com/holedaemon/gass/internal/gassfile"
)

// importResolver is responsible for resolving imports present within an input
// file. It implements the [godartsass.ImportResolver] interface.
type importResolver struct {
	baseDir string
}

// CanonicalizeURL creates a canonical version of the given URL, assuming it
// can resolve it. Otherwise, it will return an empty string.
func (r *importResolver) CanonicalizeURL(url string) (string, error) {
	// Return the directory and file of the URL.
	// If dir is empty, then the file is in the base directory.
	dir, file := filepath.Split(url)
	if dir == "" {
		dir = r.baseDir
	} else {
		dir = filepath.Join(r.baseDir, dir)
	}

	// Get a slice containing the files in the directory.
	files, err := os.ReadDir(dir)
	if err != nil {
		return "", err
	}

	// Iterate each file in files and check if it matches the file in URL.
	// We also check to see if the file is a partial, e.g. starts with an
	// underscore.
	hasFile, isPartial := containsFile(files, file)
	if hasFile {
		name := file

		// If the file is a partial, make sure we reflect that in the
		// canonical URL.
		if isPartial {
			name = "_" + name
		}

		// If the file's extension is blank, assume it is using SCSS syntax.
		ext := filepath.Ext(file)
		if ext == "" {
			name = name + ".scss"
		}

		// Create the entire path and return it with the appropriate protocol.
		wholePath := filepath.Join(dir, name)
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
