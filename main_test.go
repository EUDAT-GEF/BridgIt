package main

import (
	"testing"

	"bytes"
	"fmt"
	"log"
	"os"
	"reflect"
	"runtime"
	"strings"

	"github.com/EUDAT-GEF/Bridgit/api"
	"github.com/EUDAT-GEF/Bridgit/utils"

	"encoding/json"
	"net/http"
	"net/http/httptest"

	"github.com/EUDAT-GEF/Bridgit/def"
)

func TestClient(t *testing.T) {
	config, err := utils.ReadConfigFile("./def/config.json")
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

	req, err = http.NewRequest("POST", "/jobs", nil)
	if err != nil {
		t.Fatal(err)
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(app.JobStart)

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

	log.Println("Stopping HTTP server")
	err = app.Server.Shutdown(nil)
	CheckErr(t, err)

	//
	//// overwrite this because when testing we're in a different working directory
	//config.Pier.InternalServicesFolder = internalServicesFolder
	//
	//db, file, err := db.InitDbForTesting()
	//CheckErr(t, err)
	//defer db.Close()
	//defer os.Remove(file)
	//user, _ := AddUserWithToken(t, db, name1, email1)
	//
	//pier, err := pier.NewPier(&db, config.Pier, config.TmpDir, config.Timeouts)
	//CheckErr(t, err)
	//
	//connID, err := pier.AddDockerConnection(0, config.Docker)
	//CheckErr(t, err)
	//
	//before, err := db.ListServices()
	//CheckErr(t, err)
	//
	//service, err := pier.BuildService(connID, user.ID, "./clone_test")
	//CheckErr(t, err)
	//log.Print("test service built: ", service.ID, " ", service.ImageID)
	//log.Printf("test service built: %#v", service)
	//
	//after, err := db.ListServices()
	//CheckErr(t, err)
	//
	//errstr := "Cannot find new service in list"
	//for _, x := range after {
	//	if x.ID == service.ID {
	//		errstr = ""
	//		break
	//	}
	//}
	//
	//if errstr != "" {
	//	t.Error("before is: ", len(before), before)
	//	t.Error("service is: ", service)
	//	t.Error("after is: ", len(after), after)
	//	t.Error("")
	//	t.Error(errstr)
	//	t.Fail()
	//	return
	//}
	//
	//job, err := pier.RunService(user.ID, service.ID, testPID, config.Limits, config.Timeouts)
	//CheckErr(t, err)
	//for job.State.Code == -1 {
	//	job, err = db.GetJob(job.ID)
	//	CheckErr(t, err)
	//}
	//
	//log.Print("test job: ", job.ID)
	//// log.Printf("test job: %#v", job)
	//
	//jobList, err := db.ListJobs()
	//Expect(t, len(jobList) != 0)
	//
	//found := false
	//for _, j := range jobList {
	//	if j.ID == job.ID {
	//		found = true
	//	}
	//}
	//Expect(t, found)
	//
	//j, err := db.GetJob(job.ID)
	//CheckErr(t, err)
	//ExpectEquals(t, j.ID, job.ID)
}

func TestMain(m *testing.M) {
	code := m.Run()
	os.Exit(code)
}

func CheckErr(t *testing.T, err error) {
	if err != nil {
		t.Log(err, caller())
		t.FailNow()
	}
}

func Expect(t *testing.T, condition bool) {
	if !condition {
		t.Log("Expectation failed", caller())
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
