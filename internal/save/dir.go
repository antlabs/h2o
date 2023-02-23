package save

import (
	"os"
	"path"
)

func Mkdir(dir string, packageName string) (dirWithPackage string) {
	dir = dir + "/" + packageName
	dir = path.Clean(dir)
	os.MkdirAll(dir, 0755)
	return dir
}

func MkdirAndClean(dir string) string {
	dir = path.Clean(dir)
	os.MkdirAll(dir, 0755)
	return dir
}
