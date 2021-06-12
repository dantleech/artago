package action

import (
	"fmt"
	"io"
	"os"
	"path"

	art "github.com/dantleech/artag/artifact"
	"github.com/dantleech/artag/config"
)

type copyParams struct {
	Destination string `yaml:"destination"`
}

func CopyAction(artifact art.Artifact, action config.Action) {
	params := copyParams{}
	art.UnmarshallParams(action.Params, &params)
	dest := art.ResolveArtifactParameter(artifact, params.Destination)
	ensureDirectoryExists(dest)

	df, err := os.Create(dest)

	if err != nil {
		panic(fmt.Sprintf("Could not create file `%v`: %v", dest, err))
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
