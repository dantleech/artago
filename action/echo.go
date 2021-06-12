package action

import (
	"log"

	art "github.com/dantleech/artago/artifact"
	"github.com/dantleech/artago/config"
)

type echoParams struct {
	String string `yaml:"string"`
}

func EchoAction(artifact art.Artifact, action config.Action) {
	params := echoParams{}
	art.UnmarshallParams(action.Params, &params)
	log.Println(params.String)
}
