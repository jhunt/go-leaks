go-leaks - A Leak Detecter for Go Unit Tests
============================================

Go seduces with its garbage-collected memory model -- never worry
about resource allocation and de-allocation again!!

Except for file descriptors.  This is bad, in production code:

    package leaky

    import (
      "net/http"
    )

    func Ping() bool {
      res, err := http.Get("https://jameshunt.us/")
      if err != nil {
        return false
      }

      return res.StatusCode == 200
    }

... because it leaks a file descriptor.

Now, with `go-leaks`, we can test for this type of bad behavior:

    package leaky_test

    import (
      "testing"

      "github.com/jhunt/go-leaks"

      "my/code/leaky"
    )

    func TestPing(t *testing.T) {
      if !leaky.Ping() {
        t.Errorf("failed to ping!")
      }
      if leaks.Files(func() { leaky.Ping() }) {
        t.Errorf("Ping() leaked one or more file descriptors!")
      }
    }

This test fails.

    →  go test ./examples/leaky/
    --- FAIL: TestPing (0.65s)
        leaky_test.go:16: Ping() leaked one or more file descriptors!
    FAIL
    FAIL	github.com/jhunt/go-leaks/examples/leaky	0.668s
    FAIL

Happy Testing!
