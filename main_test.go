package main

import (
	"testing"

	"bytes"
	"fmt"
	"log"
	"reflect"
	"runtime"
	"strings"

	"github.com/EUDAT-GEF/BridgIt/api"
	"github.com/EUDAT-GEF/BridgIt/utils"

	"encoding/json"
	"net/http"
	"net/http/httptest"

	"github.com/EUDAT-GEF/BridgIt/def"
)

var accessToken = "yJYcu3KjXqawyaeMIKPBuJc1ArCkAGFJIDQwgf89wP5JBOEl"
var inputFile = "http://weblicht.sfs.uni-tuebingen.de/clrs/storage/1507625318314.txt"
var testService = "fake"
var testConfig = "./tests/test_config.json"

func TestClient(t *testing.T) {
	config, err := utils.ReadConfigFile(testConfig)
	CheckErr(t, err)

	app := api.NewApp(config)
	go app.Start()

	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(app.Index)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	expected := app.Info
	var infoReply def.Info
	err = json.NewDecoder(rr.Body).Decode(&infoReply)

	CheckErr(t, err)
	ExpectEquals(t, infoReply, expected)

	req, err = http.NewRequest("POST", "/jobs?service="+testService+"&token="+accessToken+"&input="+inputFile, nil)
	if err != nil {
		t.Fatal(err)
	}

	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(app.JobStart)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expectedLink := app.Config.GEFAddress + "/api/volumes/OutputVolume/results.txt?content&access_token=" + accessToken

	CheckErr(t, err)
	ExpectEquals(t, rr.Body.String(), expectedLink)

	log.Println("Stopping HTTP server")
	err = app.Server.Shutdown(nil)
	CheckErr(t, err)
}

func CheckErr(t *testing.T, err error) {
	if err != nil {
		t.Log(err, caller())
		t.FailNow()
	}
}

func ExpectEquals(t *testing.T, left, right interface{}) {
	if !reflect.DeepEqual(left, right) {
		t.Logf("Not Equals:\n%#v\n%#v\n@%s", left, right, caller())
		t.FailNow()
	}
}

func ExpectNotEquals(t *testing.T, left, right interface{}) {
	if reflect.DeepEqual(left, right) {
		t.Logf("Equals (should not be):\n%#v\n%#v\n@%s", left, right, caller())
		t.FailNow()
	}
}

func ExpectNotNil(t *testing.T, value interface{}) {
	if value == nil {
		t.Log("Unexpected NIL value", caller())
		t.FailNow()
	}
}

func caller() string {
	var b bytes.Buffer
	for i := 2; i < 5; i++ {
		_, file, line, ok := runtime.Caller(i)
		if ok &&
			!strings.HasSuffix(file, "/src/testing/testing.go") &&
			!strings.Contains(file, "/src/runtime/") {
			b.WriteString(fmt.Sprint("\n", file, ":", line))
		}
	}
	return b.String()
}
