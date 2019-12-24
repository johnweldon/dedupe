package main

import (
	"encoding/hex"
	"fmt"
	"hash/fnv"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
)

var files = map[string][]string{}

func main() {
	if err := filepath.Walk(".", walkFn); err != nil {
		fmt.Fprintf(os.Stderr, "Error walking: %v\n", err)
	}
	for key, val := range files {
		if len(val) > 1 {
			fmt.Printf("Duplicates: %q %+v\n", key, val)
		}
		fd, err := os.Create(path.Join("out", key+".jpg"))
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating file: %v\n", err)
			continue
		}
		defer fd.Close()
		orig, err := os.Open(val[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
			continue
		}
		defer orig.Close()

		if _, err = io.Copy(fd, orig); err != nil {
			fmt.Fprintf(os.Stderr, "Error copying file: %v\n", err)
			continue
		}
	}
}

func walkFn(path string, info os.FileInfo, err error) error {
	if !info.IsDir() && isImage(path) {
		if hash, err := hashFn(path); err == nil {
			fmt.Printf("Path: %q: %s\n", path, hash)
			files[hash] = append(files[hash], path)
		}
	}
	return err
}

func isImage(path string) bool {
	for _, suffix := range []string{".jpg", ".jpeg", ".JPG"} {
		if strings.HasSuffix(path, suffix) {
			return true
		}
	}
	fmt.Fprintf(os.Stderr, "Skipping: %q\n", path)
	return false
}

func hashFn(path string) (string, error) {
	fd, err := os.Open(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening: %v\n", err)
		return "", err
	}
	defer fd.Close()

	hasher := fnv.New128a()
	if _, err := io.Copy(hasher, fd); err != nil {
		fmt.Fprintf(os.Stderr, "Error hashing: %v\n", err)
		return "", err
	}

	return hex.EncodeToString(hasher.Sum(nil)), nil
}
