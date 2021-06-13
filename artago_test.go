package main

import (
	"io"
	"io/ioutil"
	"log"
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

func initialize(w test.Workspace, c config.Config) (*httptest.ResponseRecorder, application) {
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

	return response, application
}

func createRequest(files map[string]string) *http.Request {
	pr, pw := io.Pipe()
	request := httptest.NewRequest("POST", "/artifact/upload", pr)

	mWriter := multipart.NewWriter(pw)
	go func() {
		defer mWriter.Close()
		for fieldName, path := range files {
			body, err := ioutil.ReadFile(path)

			if err != nil {
				panic(err)
			}

			part, err := mWriter.CreateFormFile(fieldName, path)
			part.Write(body)

			if err != nil {
				panic(err)
			}
		}
	}()
	request.Header.Add("Content-Type", mWriter.FormDataContentType())

	return request
}

func TestUploadFileAndApplyRule(t *testing.T) {

	request := createRequest(map[string]string{
		"file1": "testdata/json_artifact.json",
	})
	workspace := createWorkspace()
	response, application := initialize(workspace, config.Config{
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
	request := createRequest(map[string]string{
		"file1": "testdata/json_artifact.json",
	})
	response, application := initialize(workspace, config.Config{})

	request.Header.Set("BuildId", "ABCDE")

	handler := http.HandlerFunc(application.Application)
	handler.ServeHTTP(response, request)

	assert.JSONEq(t, `{"BuildId": "ABCDE", "Results": {}}`, response.Body.String())
}

func TestUploadMultipleFiles(t *testing.T) {
	workspace := createWorkspace()
	request := createRequest(map[string]string{
		"file1": "testdata/json_artifact.json",
		"file2": "testdata/json_artifact.json",
		"file3": "testdata/json_artifact.json",
	})
	response, application := initialize(workspace, config.Config{
		Rules: []config.Rule{
			{
				Predicate: "true",
				Actions: []config.Action{
					{
						Type: "publishLink",
						Params: map[string]interface{}{
							"name":     "foo",
							"template": "http://hello",
						},
					},
				},
			},
		},
	})
	request.Header.Set("BuildId", "ABCDE")

	handler := http.HandlerFunc(application.Application)
	handler.ServeHTTP(response, request)

	log.Printf("%v", response.Body.String())
	assert.JSONEq(t, `{"BuildId": "ABCDE", "Results": {"links":{"foo":"http://hello"}}}`, response.Body.String())
}
