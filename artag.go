package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"

	"github.com/dantleech/artag/action"
	"github.com/dantleech/artag/artifact"
	"github.com/dantleech/artag/config"
)

func main() {
	Start()
}

type application struct {
	config config.Config
}

func Start() {
	config := config.LoadConfig("")
	log.Println(fmt.Sprintf("Listening for requests on `%s`", config.Address))
	application := application{
		config: config,
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
	mux.Handle("/", http.FileServer(http.Dir(a.config.PublicDir)))
	mux.ServeHTTP(response, request)
}

func (a application) artifactUploadHandler(response http.ResponseWriter, request *http.Request) {
	request.ParseMultipartForm(32 << 20)

	for fileName := range request.MultipartForm.File {
		file, header, err := request.FormFile(fileName)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		if _, err := os.Stat(a.config.WorkspacePath); os.IsNotExist(err) {
			err := os.MkdirAll(a.config.WorkspacePath, 0777)
			if err != nil {
				log.Fatalf("Could not create workspace at `%s`: %s", a.config.WorkspacePath, err)
			}
		}

		destFilePath := path.Join(a.config.WorkspacePath, fileName)
		destFile, err := os.Create(destFilePath)

		if err != nil {
			log.Fatal(err)
		}

		_, e := io.Copy(destFile, file)

		if e != nil {
			log.Fatalf("Could not copy file: %s", err)
		}

		p := artifact.Processor{
			Rules: a.config.Rules,
			Actions: map[string]artifact.ActionHandler{
				"copy": action.CopyAction,
			},
		}

		artifact := artifact.NewArtifactFromFile(destFile, header)
		log.Printf("Processing file `%s` (%s)", destFilePath, artifact.Name)
		destFile.Close()
		p.Process(artifact)
		os.Remove(destFilePath)
		file.Close()
		log.Printf("Processed: %s", artifact.Name)
	}
}
