package main

import (
	"github.com/kjk/lzmadec"
)

func uncompress7Z(file string) error {
	a, err := lzmadec.NewArchive(file)
	if err != nil {
		return err
	}

	for _, e := range a.Entries {
		err = a.ExtractToFile(e.Path, e.Path)
		if err != nil {
			return err
		}
	}

	return nil
}
