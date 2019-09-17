package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
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

func main() {
	if len(os.Args) != 3 {
		fmt.Fprintf(os.Stderr, errorMessage)
		os.Exit(1)
	}

	prefix := os.Args[1]
	if filepath.Clean(prefix) != prefix {
		fmt.Fprintf(os.Stderr, "first argument must be a clean path\n")
		os.Exit(1)
	}
	if prefix == "" {
		fmt.Fprintf(os.Stderr, "first argument must be non-empty\n")
		os.Exit(1)
	}
	if !filepath.IsAbs(prefix) {
		fmt.Fprintf(os.Stderr, "first argument must be an absolute path\n")
		os.Exit(1)
	}

	dir := os.Args[2]
	if filepath.Clean(dir) != dir {
		fmt.Fprintf(os.Stderr, "second argument must be a clean path\n")
		os.Exit(1)
	}
	if dir == "" {
		fmt.Fprintf(os.Stderr, "second argument must be non-empty\n")
		os.Exit(1)
	}
	if dir == "." {
		fmt.Fprintf(os.Stderr, "second argument must not be a dot\n")
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
