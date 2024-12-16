package sass

import (
	"os"
	"strings"
)

// containsFile iterates files in a directory to check for two things:
// 1) If the file name matches the given name.
// 2) If the file is a partial (i.e. starts with an underscore).
func containsFile(files []os.DirEntry, file string) (bool, bool) {
	for _, e := range files {
		if e.IsDir() {
			return false, false
		}

		// If the file in the directory starts with the given file name
		// (no underscore), return true.
		if strings.HasPrefix(e.Name(), file) {
			return true, false
		}

		// If the file in the directory starts with the given file name
		// (with underscore), return true
		file = "_" + file
		if strings.HasPrefix(e.Name(), file) {
			return true, true
		}
	}

	return false, false
}
