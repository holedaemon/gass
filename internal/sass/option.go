package sass

import (
	"log/slog"
	"time"
)

type Option func(*Runner)

func Debug() Option {
	return func(r *Runner) {
		r.debug = true
	}
}

func DartBinary(path string) Option {
	return func(r *Runner) {
		r.binary = path
	}
}

func Timeout(dur time.Duration) Option {
	return func(r *Runner) {
		r.timeout = dur
	}
}

func Logger(l *slog.Logger) Option {
	return func(r *Runner) {
		r.l = l
	}
}
