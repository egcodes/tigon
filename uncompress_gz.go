package main

import (
	"bufio"
	"compress/gzip"
	"io"
	"os"
)

func uncompressGZ(file string) error {
	reader, err := os.Open(file)

	gzr, err := gzip.NewReader(reader)
	if err != nil {
		return err
	}
	defer gzr.Close()

	f, err := os.Create(Config.Path.Raw + "/" + FolderName + "/" + gzr.Name)
	if err != nil {
		return err
	}
	w := bufio.NewWriter(f)
	if _, err := io.Copy(w, gzr); err != nil {
		return err
	}

	w.Flush()
	f.Sync()
	f.Close()

	if err := gzr.Close(); err != nil {
		return err
	}

	return nil
}
