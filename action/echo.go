package action

import (
	"log"

	"github.com/dantleech/artag/config"
	"github.com/dantleech/artag/processor"
)

type echoParams struct {
	String string `yaml:"string"`
}

func EchoAction(artifact processor.Artifact, action config.Action) {
	params := echoParams{}
	processor.UnmarshallParams(action.Params, &params)
	log.Println(params.String)
}
