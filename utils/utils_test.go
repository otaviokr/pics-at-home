package utils

import (
	"encoding/json"
	"io/ioutil"
	"net/http/httptest"
	"strings"
	"testing"
)

// Message wraps up a message in JSON-compliant format.
func TestMessage(t *testing.T) {
	expected := map[string]interface{}{"status": false, "message": "this is a test"}

	actual := Message(false, "this is a test")

	if len(actual) != len(expected) {
		t.Fatalf("Wrong number of elements! actual[%d] - expected[%d]", len(actual), len(expected))
	}

	if actual["status"] != expected["status"] {
		t.Fatalf("Status mismatch! actual[%v] - expected[%v]", actual["status"], expected["status"])
	}

	if actual["message"] != expected["message"] {
		t.Fail()
	}
}

// TestRespond sends a JSON response.
//func TestRespond(w http.ResponseWriter, data map[string]interface{}) {
func TestRespond(t *testing.T) {
	expectedCode := 200
	expectedContentType := []string{"application/json"}
	expectedBody := map[string]interface{}{"status": false, "message": "this is a test"}
	w := httptest.NewRecorder()

	Respond(w, expectedBody)

	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != expectedCode {
		t.Fatalf("Status code mismatch: actual[%d] - expected[%d]", resp.StatusCode, expectedCode)
	}

	if resp.Header["Content-Type"][0] != expectedContentType[0] {
		t.Fatalf("Content-type mismatch: actual[%s] - expected[%s]", resp.Header["Content-type"], expectedContentType)
	}

	bodyAsArray, err := json.Marshal(expectedBody)
	if err != nil {
		t.Fatal(err)
	}

	if strings.Compare(string(bodyAsArray), strings.ReplaceAll(string(body), "\n", "")) != 0 {
		t.Fatalf("Body mismatch:\n\n\tactual\n%s\n\n\texpected\n%s", string(body), string(bodyAsArray))
	}
}
