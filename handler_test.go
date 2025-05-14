package gocharge_test

import (
	"testing"

	gocharge "github.com/huijiro/go-charge"
)

func TestHandler(t *testing.T) {
	server.Get("/test", func(w gocharge.Response[any], r gocharge.Request[any]) error {
		return nil
	},
	)
}
