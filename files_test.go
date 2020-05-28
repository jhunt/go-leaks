package leaks_test

import (
	"io/ioutil"
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"testing"

	"github.com/jhunt/go-leaks"
)

func TestFileLeaks(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "File Leaks Test Suite")
}

func leaky() {
	_, err := ioutil.TempFile("", "leaks")
	Ω(err).ShouldNot(HaveOccurred())
}

func airtight() {
	f, err := ioutil.TempFile("", "leaks")
	defer f.Close()
	Ω(err).ShouldNot(HaveOccurred())
}

func exchange(f *os.File) func() {
	return func() {
		f.Close()
		leaky()
	}
}

var _ = Describe("leaks", func() {
	Context("with a leaky function", func() {
		if !leaks.CanDetectFileLeaks() {
			Skip("leaks library cannot detect file leaks")
		}

		It("should detect our leaks", func() {
			Ω(leaks.Files(leaky)).Should(BeTrue())
			Ω(leaks.NoFiles(leaky)).Should(BeFalse())
		})
	})

	Context("with an airtight function", func() {
		if !leaks.CanDetectFileLeaks() {
			Skip("leaks library cannot detect file leaks")
		}

		It("should detect no leaks", func() {
			Ω(leaks.Files(airtight)).Should(BeFalse())
			Ω(leaks.NoFiles(airtight)).Should(BeTrue())
		})
	})

	Context("with a change to the contents, but not the counts, of open files", func() {
		if !leaks.CanDetectFileLeaks() {
			Skip("leaks library cannot detect file leaks")
		}

		It("should detect the leak as a loss+gain", func() {
			f, err := ioutil.TempFile("", "leaks")
			Ω(err).ShouldNot(HaveOccurred())
			Ω(leaks.Files(exchange(f))).Should(BeTrue())

			f, err = ioutil.TempFile("", "leaks")
			Ω(err).ShouldNot(HaveOccurred())
			Ω(leaks.NoFiles(exchange(f))).Should(BeFalse())
		})
	})

	Context("without lsof", func() {
		It("cannot detect file descriptor leaks", func() {
			os.Setenv("LEAKS_PATH_TO_LSOF", "/path/to/nowhere")
			Ω(leaks.CanDetectFileLeaks()).Should(BeFalse())
			os.Unsetenv("LEAKS_PATH_TO_LSOF")
		})
	})
})
