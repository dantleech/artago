package test

import (
	"io/ioutil"
	"os"
	"path"
)

type Workspace struct {
	WorkspacePath string
}

func (w Workspace) Reset() {
	os.RemoveAll(w.WorkspacePath)
	os.MkdirAll(w.WorkspacePath, 0777)
}

func (w Workspace) Path(subPath string) string {
	return path.Join(w.WorkspacePath, subPath)
}

func (w Workspace) Put(subPath string, contents []byte) {
	err := ioutil.WriteFile(w.Path(subPath), contents, 0777)
	if err != nil {
		panic(err)
	}
}

func (w Workspace) Open(subPath string) *os.File {
	file, err := os.Open(w.Path(subPath))
	if err != nil {
		panic(err)
	}
	return file
}

func (w Workspace) Get(subPath string) []byte {
	bytes, err := ioutil.ReadFile(w.Path(subPath))
	if err != nil {
		panic(err)
	}
	return bytes
}
