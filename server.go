package main

import (
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"net/http"
)

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

func getBackupDirectoryPath(dir string, files []fs.FileInfo, r *http.Request) string {
	for _, file := range files {
		if file.Name() != "saving" && file.Name() != "countries.txt" {
			dir += file.Name() + r.URL.Path + ".txt"
		}
	}
	return dir
}

func countryHandler(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)

	dir := "./date-backups/"
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		fmt.Println(err)
		return
	}

	dir = getBackupDirectoryPath(dir, files, r)

	data, err2 := ioutil.ReadFile(dir)
	if err2 != nil {
		fmt.Println(err)
		return
	}

	fmt.Fprintf(w, "%s", string(data))
}

func countriesHandler(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)

	dir := "./date-backups/countries.txt"

	data, err2 := ioutil.ReadFile(dir)
	if err2 != nil {
		fmt.Fprintf(w, "Error!")
		return
	}

	fmt.Fprintf(w, "%s", string(data))
}

func main() {
	// saveData()

	http.HandleFunc("/", countryHandler)
	http.HandleFunc("/countries", countriesHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
