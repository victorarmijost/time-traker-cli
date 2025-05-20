package utils

import (
	"fmt"
	"os"
	"path/filepath"
)

var execPath string

func init() {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	execPath = filepath.Dir(ex)
}

func GeAppPath(folder string) string {
	return fmt.Sprintf("%s/%s", execPath, folder)
}
