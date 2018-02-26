package htmlLoader

import (
	"path/filepath"

	"github.com/lvl5hm/go-bundler/loaders"
)

func LoadFile(fileName, bundleDir string) ([]byte, []string, error) {
	ext := filepath.Ext(fileName)
	objectName := loaders.CreateVarNameFromPath(fileName)

	dstFileName := bundleDir + "/" + objectName + ext
	err := copyFile(dstFileName, fileName)
	if err != nil {
		return nil, nil, err
	}

	res := "moduleFns." + objectName + "=" +
		"function(){return {exports:'" + objectName + ext + "'}};"

	return []byte(res), nil, nil
}
