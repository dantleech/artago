package artag

import (
    "fmt"
    "log"
    "net/http"
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
    response.WriteHeader(http.StatusNotFound)
    response.Write([]byte(fmt.Sprintf("No handler found for URL: %s", request.URL)))
}


