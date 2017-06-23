package main

import (

	"fmt"

	"io/ioutil"
	"net/http"

	"log"
	"os"
	"path/filepath"
)




func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Welcome! This is Bridgit, a liaison between Weblicht and the GEF\n")
}

func JobStart(w http.ResponseWriter, r *http.Request) {
	buf, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal("request", err)
	}

	fileName := filepath.Join("data", PseudoUUID())
	f, err := os.Create(fileName)

	defer f.Close()
	_, err = f.Write(buf)
	if err != nil {
		log.Fatal("Error while writing in a file", err)
	}

	//w.Header().Set("Content-Type", "application")
	w.Write(buf)

}