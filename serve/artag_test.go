package artag

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUploadFile(t *testing.T) {
    writer := httptest.NewRecorder()
    request := httptest.NewRequest("POST", "/artifact/upload", nil)
    handler:= http.HandlerFunc(Application);
    handler.ServeHTTP(writer, request)
    if writer.Code != 200 {
        t.Errorf("Expected code 200 got: %d", writer.Code)
    }
}
