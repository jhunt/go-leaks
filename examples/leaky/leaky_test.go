package leaky_test

import (
	"testing"

	"github.com/jhunt/go-leaks"

	"github.com/jhunt/go-leaks/examples/leaky"
)

func TestPing(t *testing.T) {
	if !leaky.Ping() {
		t.Errorf("failed to ping!")
	}
	if leaks.Files(func() { leaky.Ping() }) {
		t.Errorf("Ping() leaked one or more file descriptors!")
	}
}
