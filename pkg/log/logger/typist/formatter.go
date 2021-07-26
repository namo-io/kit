package typist

type Formatter interface {
	Format(*Record) string
}

const (
	NewLine = "\r\n"
)
