package save

import (
	"fmt"
	"go/format"
	"os"
)

func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}

	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func TmplFile(fileName string, getTmpl func() []byte) {

	if b, _ := exists(fileName); !b {

		buf := getTmpl()
		fmtType, err := format.Source(buf)
		if err != nil {
			fmt.Printf("%s fail:%s\n", fileName, err)
			os.Stdout.Write(buf)
			return
		}

		//os.Stdout.Write(fmtType)
		os.WriteFile(fileName, fmtType, 0644)
	} else {
		fmt.Printf("%s 已经存在，忽略\n", fileName)
	}
}
