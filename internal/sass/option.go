package sass

import (
	"log/slog"
	"time"
)

// Option is used to configure a *[Transpiler].
type Option func(*Transpiler)

// Debug enables debug mode on a *[Transpiler].
func Debug() Option {
	return func(t *Transpiler) {
		t.debug = true
	}
}

// SassBinary tells a *[Transpiler] to use the following binary for transpilation.
// It expects path to be an absolute path to a dart-sass binary.
func SassBinary(path string) Option {
	return func(t *Transpiler) {
		t.binary = path
	}
}

// Timeout sets the maximum alloted time a *[Transpiler] may wait for
// transpilation.
func Timeout(dur time.Duration) Option {
	return func(t *Transpiler) {
		t.timeout = dur
	}
}

// Logger sets the *[slog.Logger] a *[Transpiler] uses for logging.
func Logger(l *slog.Logger) Option {
	return func(t *Transpiler) {
		t.l = l
	}
}
