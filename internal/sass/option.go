package sass

import (
	"log/slog"
	"time"
)

// Option is used to configure a *[Runner].
type Option func(*Runner)

// Debug enables debug mode on a *[Runner].
func Debug() Option {
	return func(r *Runner) {
		r.debug = true
	}
}

// SassBinary tells a *[Runner] to use the following binary for transpilation.
// It expects path to be an absolute path to a dart-sass binary.
func SassBinary(path string) Option {
	return func(r *Runner) {
		r.binary = path
	}
}

// Timeout sets the maximum alloted time a *[Runner] may wait for
// transpilation.
func Timeout(dur time.Duration) Option {
	return func(r *Runner) {
		r.timeout = dur
	}
}

// Logger sets the *[slog.Logger] a *[Runner] uses for logging.
func Logger(l *slog.Logger) Option {
	return func(r *Runner) {
		r.l = l
	}
}
