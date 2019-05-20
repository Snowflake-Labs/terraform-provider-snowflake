package instrumentedsql

// Opt is a functional option type for the wrapped driver
type Opt func(*wrappedDriver)

// WithLogger sets the logger of the wrapped driver to the provided logger
func WithLogger(l Logger) Opt {
	return func(w *wrappedDriver) {
		w.Logger = l
	}
}

// WithTracer sets the tracer of the wrapped driver to the provided tracer
func WithTracer(t Tracer) Opt {
	return func(w *wrappedDriver) {
		w.Tracer = t
	}
}
