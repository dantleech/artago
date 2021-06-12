package artifact

import (
	"log"
	"mime/multipart"
	"os"
)

type Artifact struct {
	Path string
	Name string
	Size int64
}

func (a Artifact) OpenFile() *os.File {
	file, err := os.Open(a.Path)
	if err != nil {
		log.Fatalf("Could not open artifact file `%s`", a.Path)
	}

	return file
}

func NewArtifactFromFile(file *os.File, header *multipart.FileHeader) Artifact {
	return Artifact{
		Path: file.Name(),
		Name: header.Filename,
		Size: header.Size,
	}
}
