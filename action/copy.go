package action

import (
	"fmt"
	"io"
	"os"
	"path"

	"github.com/dantleech/artag/config"
	"github.com/dantleech/artag/processor"
)

type copyParams struct {
	Destination string `yaml:"destination"`
}

func CopyAction(artifact processor.Artifact, action config.Action) {
	params := copyParams{}
	processor.UnmarshallParams(action.Params, &params)
	ensureDirectoryExists(params.Destination)

	df, err := os.Create(params.Destination)

	if err != nil {
		panic(fmt.Sprintf("Could not create file `%v`: %v", params.Destination, err))
	}

	defer df.Close()

	artifactFile := artifact.OpenFile()
	io.Copy(df, artifactFile)
	artifactFile.Close()
}

func ensureDirectoryExists(filePath string) {
	dn := path.Dir(filePath)
	if _, err := os.Stat(dn); os.IsExist(err) {
		return
	}
	os.MkdirAll(dn, 0777)
}
