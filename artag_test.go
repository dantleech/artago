package main

import (
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUploadFile(t *testing.T) {
    pr, pw := io.Pipe()

    body, err := ioutil.ReadFile("testdata/json_artifact.json")

    if err != nil {
        t.Error(err)
    }

    mWriter := multipart.NewWriter(pw)

    // Why is the go function required here??
    go func() {
        defer mWriter.Close()

        part, err := mWriter.CreateFormFile("artifact.json", "testdata/json_artifact.json")
        part.Write(body)

        if err != nil {
            t.Error(err)
        }
    }()

    request := httptest.NewRequest("POST", "/artifact/upload", pr)
    response := httptest.NewRecorder()
    request.Header.Add("Content-Type", mWriter.FormDataContentType())
    application := application{
        config: Config {
            workspacePath: "workspace",
        },
    }

    handler:= http.HandlerFunc(application.Application);
    handler.ServeHTTP(response, request)

    if response.Code != 200 {
        t.Errorf("Expected code 200 got: %d", response.Code)
    }
}
