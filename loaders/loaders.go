package loaders

import "strings"

type Loader interface {
	BeforeBuild(fileName string, config *ConfigJSON)
	LoadAndTransformFile(fileName string, config *ConfigJSON) (
		result []byte,
		imports []string,
		err error,
	)
	TransformFile(fileName string, src []byte, config *ConfigJSON) (
		result []byte,
		imports []string,
		err error,
	)
}

type ConfigJSON struct {
	Entry        string
	TemplateHTML string
	BundleDir    string
	WatchFiles   bool
	Env          map[string]string
	DevServer    struct {
		Enable bool
		Port   int
	}
	PermanentCache struct {
		Enable  bool
		DirName string
	}
}

func CreateVarNameFromPath(path string) string {
	newName := strings.Replace(path, "/", "_", -1)
	newName = strings.Replace(newName, ".", "_", -1)
	newName = strings.Replace(newName, "-", "_", -1)
	return newName
}
