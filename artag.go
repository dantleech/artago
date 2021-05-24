package main

import (
    "fmt"
    "io"
    "log"
    "net/http"
    "os"
    "path"
)

func main() {
    Start();
}

type config struct {
    address string;
    workspacePath string;
}

func loadConfig() config {
    return config{
        address: ":8080",
        workspacePath: "workspace",
    }
}

type application struct {
    config config;
}

func Start() {
    config := loadConfig()
    log.Println(fmt.Sprintf("Listening for requests on `%s`", config.address))
    application := application {
        config: config,
    }
    err := http.ListenAndServe(config.address, loggingMiddleware(http.HandlerFunc(application.Application)))

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

func (a application)Application(response http.ResponseWriter, request *http.Request) {
    mux := http.NewServeMux()
    mux.Handle("/artifact/upload", http.HandlerFunc(a.artifactUploadHandler))
    mux.Handle("/", http.NotFoundHandler())
    mux.ServeHTTP(response, request)
}

func (a application)artifactUploadHandler(response http.ResponseWriter, request *http.Request) {
    request.ParseMultipartForm(32 << 20)

    for fileName := range request.MultipartForm.File {
        file, _, err := request.FormFile(fileName)
        if err != nil {
            log.Fatal(err)
        }
        defer file.Close()

        if _, err := os.Stat(a.config.workspacePath); os.IsNotExist(err) {
            err := os.MkdirAll(a.config.workspacePath, 0777)
            if err != nil {
                log.Fatal(err)
            }
        }

        destFile, err := os.Create(path.Join(a.config.workspacePath, fileName))

        if err != nil {
            log.Fatal(err)
        }

        io.Copy(destFile, file)
    }
}
