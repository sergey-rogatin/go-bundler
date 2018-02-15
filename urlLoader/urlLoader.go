package urlLoader

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/lvl5hm/go-bundler/jsLoader"
)

func LoadFile(fileName, bundleDir string) ([]byte, []string, error) {
	ext := filepath.Ext(fileName)
	objectName := jsLoader.CreateVarNameFromPath(fileName)

	dstFileName := bundleDir + "/" + objectName + ext
	copyFile(dstFileName, fileName)

	res := "moduleFns." + objectName + "=" +
		"function(){return {exports:'" + objectName + ext + "'}};"

	return []byte(res), nil, nil
}

func copyFile(dst, src string) {
	from, err := os.Open(src)
	if err != nil {
		fmt.Println(err)
	}
	defer from.Close()

	to, err := os.OpenFile(dst, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println(err)
	}
	defer to.Close()
	io.Copy(to, from)
}
