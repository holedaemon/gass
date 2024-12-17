package gassfile

import (
	"os"
	"strings"

	"github.com/bep/godartsass/v2"
)

const gassPathPrefix = "$GASS_PATH_"

// DetermineSourceSyntax returns the [godartsass.SourceSyntax] of a file based
// on its extension.
func DetermineSourceSyntax(ext string) godartsass.SourceSyntax {
	ext = strings.ToLower(ext)
	switch ext {
	case ".sass":
		return godartsass.SourceSyntaxSASS
	case ".css":
		return godartsass.SourceSyntaxCSS
	case ".scss":
		fallthrough
	default:
		return godartsass.SourceSyntaxSCSS
	}
}

func resolveSource(path string) (string, error) {
	if strings.HasPrefix(path, gassPathPrefix) {
		sep := strings.IndexRune(path, os.PathSeparator)
	}
}
