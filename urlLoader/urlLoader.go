package urlLoader

import (
	"io"
	"os"
	"path/filepath"

	"github.com/lvl5hm/go-bundler/loaders"
)

var Loader urlLoader

type urlLoader struct{}

func (u urlLoader) BeforeBuild(fileName string, config *loaders.ConfigJSON) {
	// TODO: check if file exists in destination folder and do not copy
	ext := filepath.Ext(fileName)
	objectName := loaders.CreateVarNameFromPath(fileName)

	dstFileName := config.BundleDir + "/" + objectName + ext
	copyFile(dstFileName, fileName)
}

func (u urlLoader) LoadAndTransformFile(
	fileName string,
	config *loaders.ConfigJSON,
) ([]byte, []string, error) {
	return u.TransformFile(fileName, nil, config)
}

func (u urlLoader) TransformFile(
	fileName string,
	src []byte,
	config *loaders.ConfigJSON,
) ([]byte, []string, error) {
	ext := filepath.Ext(fileName)
	objectName := loaders.CreateVarNameFromPath(fileName)

	res := "moduleFns." + objectName + "=" +
		"function(){return {exports:'" + objectName + ext + "'}};"

	return []byte(res), nil, nil
}

func LoadFile(fileName string, config *loaders.ConfigJSON) ([]byte, []string, error) {
	ext := filepath.Ext(fileName)
	objectName := loaders.CreateVarNameFromPath(fileName)

	dstFileName := config.BundleDir + "/" + objectName + ext
	err := copyFile(dstFileName, fileName)
	if err != nil {
		return nil, nil, err
	}

	res := "moduleFns." + objectName + "=" +
		"function(){return {exports:'" + objectName + ext + "'}};"

	return []byte(res), nil, nil
}

func copyFile(dst, src string) error {
	from, err := os.Open(src)
	if err != nil {
		return err
	}
	defer from.Close()

	to, err := os.OpenFile(dst, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer to.Close()
	_, err = io.Copy(to, from)
	return err
}
