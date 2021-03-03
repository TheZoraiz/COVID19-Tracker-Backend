package main

import (
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

func getCurrentDirectory() string {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	return dir
}

func getBackupDirectoryPath(dir string, files []fs.FileInfo, r *http.Request) string {
	for _, file := range files {
		if file.Name() != "saving" && file.Name() != "countries.txt" && file.Name() != "Server" {
			dir += file.Name() + r.URL.Path + ".txt"
		}
	}
	return dir
}

func countryHandler(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)

	dir := getCurrentDirectory() + "/date-backups/"
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		fmt.Println(err)
		return
	}

	dir = getBackupDirectoryPath(dir, files, r)

	data, err2 := ioutil.ReadFile(dir)
	if err2 != nil {
		fmt.Println(err2)
		return
	}

	fmt.Println("Sending " + dir)
	fmt.Fprintf(w, "%s", string(data))
}

func countriesHandler(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)

	dir := getCurrentDirectory() + "/date-backups/countries.txt"

	data, err := ioutil.ReadFile(dir)
	if err != nil {
		fmt.Println(err)
		fmt.Fprintf(w, "Error!")
		return
	}

	fmt.Println("Sending " + dir)
	fmt.Fprintf(w, "%s", string(data))
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.HandleFunc("/", countryHandler)
	http.HandleFunc("/countries", countriesHandler)

	fmt.Println("Server listening on port:" + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
