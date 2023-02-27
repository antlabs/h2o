package file

import (
	"io"
	"os"
)

func ReadFile(fileName string) (all []byte, err error) {
	if fileName == "-" {
		all, err = io.ReadAll(os.Stdin)
	} else {
		all, err = os.ReadFile(fileName)
	}
	return
}
