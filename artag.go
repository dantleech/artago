package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
)

func Start() {
    address := ":8080"
    log.Println(fmt.Sprintf("Listening for requests on `%s`", address))
    err := http.ListenAndServe(address, loggingMiddleware(http.HandlerFunc(Application)))

    if err != nil {
        log.Fatal(err)
    }
}

func loggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
        log.Println(fmt.Sprintf("[%s] %s", request.Method, request.URL));
        next.ServeHTTP(response, request);
    });
}

func Application(response http.ResponseWriter, request *http.Request) {
    mux := http.NewServeMux()
    mux.Handle("/artifact/upload", http.HandlerFunc(artifactUploadHandler))
    mux.Handle("/", http.NotFoundHandler())
    mux.ServeHTTP(response, request)
}

func artifactUploadHandler(response http.ResponseWriter, request *http.Request) {
    request.ParseMultipartForm(32 << 20)
    outDir := "workspace"


    for fileName := range request.MultipartForm.File {
        file, _, err := request.FormFile(fileName)
        if err != nil {
            log.Fatal(err)
        }
        defer file.Close()

        if _, err := os.Stat(outDir); os.IsNotExist(err) {
            err := os.MkdirAll(outDir, 0777)
            if err != nil {
                log.Fatal(err)
            }
        }

        destFile, err := os.Create(path.Join(outDir, fileName))

        if err != nil {
            log.Fatal(err)
        }

        io.Copy(destFile, file)
    }
}


