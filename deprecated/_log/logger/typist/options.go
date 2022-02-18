package typist

type Options func(*typist)

// WithVerbose is show level(debug, trace)
func WithVerboseDisable() Options {
	return func(t *typist) {
		t.level = InfoLevel
	}
}

// WithFormatter is setter formatter
func WithFormatter(formatter Formatter) Options {
	return func(t *typist) {
		t.formatter = formatter
	}
}

// WithCllerIgnorePackageFile is setter caller ignore package name
func WithCllerIgnorePackageFile(fileName string) Options {
	return func(t *typist) {
		t.callerIgnorePackageFile = fileName
	}
}
