package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
)

func main() {
	Start()
}

type application struct {
	config    Config
	processor Processor
}

func Start() {
	config := LoadConfig("")
	log.Println(fmt.Sprintf("Listening for requests on `%s`", config.Address))
	application := application{
		config: config,
		processor: Processor{
			Rules: config.Rules,
		},
	}
	err := http.ListenAndServe(config.Address, loggingMiddleware(http.HandlerFunc(application.Application)))

	if err != nil {
		log.Fatal(err)
	}
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		log.Println(fmt.Sprintf("[%s] %s", request.Method, request.URL))
		next.ServeHTTP(response, request)
	})
}

func (a application) Application(response http.ResponseWriter, request *http.Request) {
	mux := http.NewServeMux()
	mux.Handle("/artifact/upload", http.HandlerFunc(a.artifactUploadHandler))
	mux.Handle("/", http.NotFoundHandler())
	mux.ServeHTTP(response, request)
}

func (a application) artifactUploadHandler(response http.ResponseWriter, request *http.Request) {
	request.ParseMultipartForm(32 << 20)

	for fileName := range request.MultipartForm.File {
		file, _, err := request.FormFile(fileName)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		if _, err := os.Stat(a.config.WorkspacePath); os.IsNotExist(err) {
			err := os.MkdirAll(a.config.WorkspacePath, 0777)
			if err != nil {
				log.Fatal(err)
			}
		}

		destFile, err := os.Create(path.Join(a.config.WorkspacePath, fileName))

		if err != nil {
			log.Fatal(err)
		}

		a.processor.process(NewArtifact(destFile))
	}
}
