package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

const errorMessage = `the tool expects exactly two arguments: "<path prefix> <dir name>", both must be non-empty and very strict paths` + "\n"

const WeirdFileMode = os.ModeAppend |
	os.ModeExclusive |
	os.ModeTemporary |
	os.ModeSymlink |
	os.ModeDevice |
	os.ModeNamedPipe |
	os.ModeSocket |
	os.ModeSetuid |
	os.ModeSetgid |
	os.ModeCharDevice |
	os.ModeSticky |
	os.ModeIrregular

func buildFileList(out *[]string, prefix string) {
	fi, err := os.Lstat(prefix)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed reading directory %q: %s\n", prefix, err)
		os.Exit(1)
	}
	if !fi.IsDir() {
		fmt.Fprintf(os.Stderr, "file is not a directory %q: %s\n", prefix, fi.Mode())
		os.Exit(1)
	}
	if fi.Mode()&WeirdFileMode != 0 {
		fmt.Fprintf(os.Stderr, "file has one of the blacklisted flags on it %q: %s\n", prefix, fi.Mode())
		os.Exit(1)
	}

	dirContents, err := ioutil.ReadDir(prefix)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed reading directory %q: %s\n", prefix, err)
		os.Exit(1)
	}
	for _, f := range dirContents {
		m := f.Mode()
		if m.IsDir() {
			buildFileList(out, filepath.Join(prefix, f.Name()))
		} else {
			if m&WeirdFileMode != 0 {
				fmt.Fprintf(os.Stderr, "file has one of the blacklisted flags on it %q: %s\n", filepath.Join(prefix, f.Name()), m)
				os.Exit(1)
			}
			if m.IsRegular() {
				*out = append(*out, filepath.Join(prefix, f.Name()))
			}
		}
	}
	*out = append(*out, prefix)
}

func pathCheck(a, b string) error {
	if filepath.Clean(a) != a || strings.TrimSpace(a) != a {
		return fmt.Errorf("first argument must be a clean path")
	}
	if a == "" {
		return fmt.Errorf("first argument must be non-empty")
	}
	if a == "." {
		return fmt.Errorf("first argument must not be a dot")
	}
	if a == "/" {
		return fmt.Errorf("first argument must not be a root")
	}
	if !filepath.IsAbs(a) {
		return fmt.Errorf("first argument must be an absolute path")
	}

	if filepath.Clean(b) != b || strings.TrimSpace(b) != b {
		return fmt.Errorf("second argument must be a clean path")
	}
	if strings.Contains(b, "..") {
		return fmt.Errorf("second argument must not contain ..")
	}
	if b == "" {
		return fmt.Errorf("second argument must be non-empty")
	}
	if b == "." {
		return fmt.Errorf("second argument must not be a dot")
	}
	if b == "/" {
		return fmt.Errorf("second argument must not be a root")
	}
	if filepath.IsAbs(b) {
		return fmt.Errorf("second argument must not be an absolute path")
	}
	return nil
}

func main() {
	if len(os.Args) != 3 {
		fmt.Fprintf(os.Stderr, errorMessage)
		os.Exit(1)
	}

	prefix := os.Args[1]
	dir := os.Args[2]
	if err := pathCheck(prefix, dir); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}

	var allFiles []string
	fullPath := filepath.Join(prefix, dir)
	buildFileList(&allFiles, fullPath)
	for _, f := range allFiles {
		if err := os.Remove(f); err != nil {
			fmt.Fprintf(os.Stderr, "failed removing file %q: %s\n", f, err)
			os.Exit(1)
		}
	}
	os.Exit(0)
}
