package internal

import (
	"errors"
	"fmt"
	"os"
)

func createDirectory(path string) {
	err := os.MkdirAll(path, 0755)

	if err != nil {
		if Verbose && os.IsExist(err) {
			fmt.Println(err)
		} else {
			FatalError(err)
		}
	}
}

func createFile(path string) {
	var modes int

	if Force {
		modes = os.O_CREATE | os.O_TRUNC
	} else {
		modes = os.O_CREATE | os.O_EXCL
	}

	file, err := os.OpenFile(path, modes, 0644)

	if err != nil {
		if Verbose && errors.Is(err, os.ErrExist) {
			fmt.Println(err)
		} else {
			FatalError(err)
		}
	}

	defer file.Close()
}
