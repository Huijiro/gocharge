package gocharge

type HandlerFunc[W any, R any] func(w Response[W], r Request[R]) error
