package leaks

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"syscall"
)

func lsof(out *os.File) ([]byte, error) {
	bin := os.Getenv("LEAKS_PATH_TO_LSOF")
	if bin == "" {
		bin = "lsof"
	}
	cmd := exec.Command(bin, "-nP", "-p", fmt.Sprintf("%d", syscall.Getpid()))

	if out != nil {
		out.Seek(0, 0)
		out.Truncate(0)
		cmd.Stdout = out
	}

	if err := cmd.Run(); err != nil {
		return nil, err
	}

	if out != nil {
		out.Seek(0, 0)
		return ioutil.ReadAll(out)

	} else {
		return nil, nil
	}
}

func CanDetectFileLeaks() bool {
	_, err := lsof(nil)
	return err == nil
}

func Files(fn func()) bool {
	return !NoFiles(fn)
}

func NoFiles(fn func()) bool {
	out, err := ioutil.TempFile("", "leaks")
	if err != nil {
		panic("leaks.NoFiles() failed to run lsof: " + err.Error())
	}
	defer out.Close()

	before, err := lsof(out)
	if err != nil {
		panic("leaks.NoFiles() failed to run lsof: " + err.Error())
	}

	fn()

	after, err := lsof(out)
	if err != nil {
		panic("leaks.NoFiles() railed to re-run lsof: " + err.Error())
	}

	if len(before) != len(after) {
		if os.Getenv("LEAKS_DEBUG") != "" {
			diff(before, after)
		}
		return false
	}
	for i := range before {
		if before[i] != after[i] {
			if os.Getenv("LEAKS_DEBUG") != "" {
				diff(before, after)
			}
			return false
		}
	}
	return true
}

func diff(before, after []byte) {
	loss := make(map[string]bool)
	gain := make(map[string]bool)

	for _, l := range bytes.Split(before, []byte{'\n'}) {
		loss[string(l)] = true
	}
	for _, l := range bytes.Split(after, []byte{'\n'}) {
		s := string(l)
		if _, ok := loss[s]; ok {
			delete(loss, s)
		} else {
			gain[s] = true
		}
	}

	for l := range loss {
		fmt.Printf("lost [%s]\n", l)
	}
	for l := range gain {
		fmt.Printf("gain [%s]\n", l)
	}
}
