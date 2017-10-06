package main

import (
	"log"
	"testing"
	"github.com/EUDAT-GEF/GEF/gefserver/def"
	"github.com/EUDAT-GEF/GEF/gefserver/db"
	"os"
	"github.com/EUDAT-GEF/GEF/gefserver/pier"
)

func TestClient(t *testing.T) {
	config, err := def.ReadConfigFile(configFilePath)
	CheckErr(t, err)

	// overwrite this because when testing we're in a different working directory
	config.Pier.InternalServicesFolder = internalServicesFolder

	db, file, err := db.InitDbForTesting()
	CheckErr(t, err)
	defer db.Close()
	defer os.Remove(file)
	user, _ := AddUserWithToken(t, db, name1, email1)

	pier, err := pier.NewPier(&db, config.Pier, config.TmpDir, config.Timeouts)
	CheckErr(t, err)

	connID, err := pier.AddDockerConnection(0, config.Docker)
	CheckErr(t, err)

	before, err := db.ListServices()
	CheckErr(t, err)

	service, err := pier.BuildService(connID, user.ID, "./clone_test")
	CheckErr(t, err)
	log.Print("test service built: ", service.ID, " ", service.ImageID)
	log.Printf("test service built: %#v", service)

	after, err := db.ListServices()
	CheckErr(t, err)

	errstr := "Cannot find new service in list"
	for _, x := range after {
		if x.ID == service.ID {
			errstr = ""
			break
		}
	}

	if errstr != "" {
		t.Error("before is: ", len(before), before)
		t.Error("service is: ", service)
		t.Error("after is: ", len(after), after)
		t.Error("")
		t.Error(errstr)
		t.Fail()
		return
	}

	job, err := pier.RunService(user.ID, service.ID, testPID, config.Limits, config.Timeouts)
	CheckErr(t, err)
	for job.State.Code == -1 {
		job, err = db.GetJob(job.ID)
		CheckErr(t, err)
	}

	log.Print("test job: ", job.ID)
	// log.Printf("test job: %#v", job)

	jobList, err := db.ListJobs()
	Expect(t, len(jobList) != 0)

	found := false
	for _, j := range jobList {
		if j.ID == job.ID {
			found = true
		}
	}
	Expect(t, found)

	j, err := db.GetJob(job.ID)
	CheckErr(t, err)
	ExpectEquals(t, j.ID, job.ID)
}