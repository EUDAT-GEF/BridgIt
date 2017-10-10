package main

import (
	"testing"

	"reflect"
	"bytes"
	"runtime"
	"strings"
	"fmt"
	"log"
	"os"
)

func TestClient(t *testing.T) {
	config, err := ReadConfigFile(configFilePath)
	CheckErr(t, err)

	router := NewRouter()
	srv := startHttpServer(config.PortNumber, router)





	log.Println("Stopping HTTP server")
	if err := srv.Shutdown(nil); err != nil {
		panic(err)
	}

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