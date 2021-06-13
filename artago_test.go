package main

import (
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dantleech/artago/config"
	"github.com/dantleech/artago/test"
	"github.com/imdario/mergo"
	"github.com/stretchr/testify/assert"
)

func createWorkspace() test.Workspace {
	workspace := test.Workspace{
		WorkspacePath: "workspace",
	}
	workspace.Reset()
	return workspace
}

func initialize(w test.Workspace, c config.Config) (*http.Request, *httptest.ResponseRecorder, application) {
	defaultConfig := config.Config{
		WorkspacePath: w.Path("/"),
	}
	err := mergo.Merge(&defaultConfig, c)

	if err != nil {
		panic(err)
	}

	application := application{
		config: defaultConfig,
	}
	response := httptest.NewRecorder()
	request := createRequest()

	return request, response, application
}

func TestUploadFileAndApplyRule(t *testing.T) {

	workspace := createWorkspace()
	request, response, application := initialize(workspace, config.Config{
		Rules: []config.Rule{
			{
				Predicate: "true",
				Actions: []config.Action{
					{
						Type: "copy",
						Params: map[string]interface{}{
							"destination": workspace.Path("processed.json"),
						},
					},
				},
			},
		},
	})

	handler := http.HandlerFunc(application.Application)
	handler.ServeHTTP(response, request)

	assert.Equal(t, 200, response.Code)
	assert.NoFileExists(t, "workspace/artifact.json", "Temporary file was deleted")
	assert.FileExists(t, "workspace/processed.json", "File was processed")
	assert.Containsf(t, response.Body.String(), "BuildId", "Response contains BuildId")
}

func TestBuildIdHeader(t *testing.T) {
	workspace := createWorkspace()
	request, response, application := initialize(workspace, config.Config{})

	request.Header.Set("BuildId", "ABCDE")

	handler := http.HandlerFunc(application.Application)
	handler.ServeHTTP(response, request)

	assert.JSONEq(t, `{"BuildId": "ABCDE"}`, response.Body.String())
}

func createRequest() *http.Request {
	pr, pw := io.Pipe()

	body, err := ioutil.ReadFile("testdata/json_artifact.json")

	if err != nil {
		panic(err)
	}

	mWriter := multipart.NewWriter(pw)

	// Why is the go function required here??
	go func() {
		defer mWriter.Close()

		part, err := mWriter.CreateFormFile("artifact.json", "testdata/json_artifact.json")
		part.Write(body)

		if err != nil {
			panic(err)
		}
	}()

	request := httptest.NewRequest("POST", "/artifact/upload", pr)
	request.Header.Add("Content-Type", mWriter.FormDataContentType())

	return request
}
