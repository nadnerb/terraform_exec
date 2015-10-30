package util

import (
	"path/filepath"
	"fmt"
	"os"
)

// I am guessing/hoping there is a nicer way to do this
func HasFilesWithExtension(directory string, extension string) (bool, error) {
	dirname := directory + string(filepath.Separator)

	d, err := os.Open(dirname)
	if err != nil {
		fmt.Println(err)
		return false, err
	}
	defer d.Close()

	files, err := d.Readdir(-1)
	if err != nil {
		fmt.Println(err)
		return false, err
	}

	for _, file := range files {
		if file.Mode().IsRegular() {
			if filepath.Ext(file.Name()) == extension {
				return true, nil
			}
		}
	}
	return false, nil
}

