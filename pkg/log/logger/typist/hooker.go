package typist

import "context"

type Hooker interface {
	Name() string
	Fire(context.Context, Level, *Record) error
}
