package internal

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

func loadYamlFile(fileName string, obj interface{}) error {
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		return err
	}
	err = yaml.UnmarshalStrict(data, obj)
	return err
}

func absPath(path string, dir string) (string, error) {
	if !filepath.IsAbs(path) {
		if dir != "" {
			path = filepath.Join(dir, path)
		}
		return filepath.Abs(path)
	}
	return path, nil
}

func errServiceName(err error, name string) error {
	return fmt.Errorf("Error in service %s: %s", name, err.Error())
}
