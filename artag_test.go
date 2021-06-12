package main

import (
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dantleech/artag/config"
	"github.com/dantleech/artag/test"
	"github.com/stretchr/testify/assert"
)

func TestUploadFile(t *testing.T) {
	workspace := test.Workspace{
		WorkspacePath: "workspace",
	}
	workspace.Reset()
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
		config: config.Config{
			Address:       "",
			WorkspacePath: workspace.Path("/"),
			Rules: []config.Rule{
				{
					Predicate: "true",
					Actions: []config.Action{
						{
							Type: "copy",
							Params: map[string]interface{}{
								"destination": "workspace/processed.json",
							},
						},
					},
				},
			},
			PublicDir: "",
		},
	}

	handler := http.HandlerFunc(application.Application)
	handler.ServeHTTP(response, request)

	assert.Equal(t, 200, response.Code)

	assert.NoFileExists(t, "workspace/artifact.json", "Temporary file was deleted")
	assert.FileExists(t, "workspace/processed.json", "File was processed")
}
