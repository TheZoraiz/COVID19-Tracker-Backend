package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

func itExists(name string) bool {
	if _, err := os.Stat(name); os.IsNotExist(err) {
		return false
	} else {
		return true
	}
}

func FetchApiString(url string) (string, error) {
	response, err := http.Get(url)
	if err != nil {
		return "", errors.New("Failed to fetch")
	}
	defer response.Body.Close()
	body, err2 := ioutil.ReadAll(response.Body)
	if err2 != nil {
		return "", errors.New("Failed to read")
	}
	return string(body), nil
}

func saveData() {
	dir := "./date-backups/"

	if !itExists(dir) {
		fmt.Println("Making initial directory...")
		os.Mkdir(dir, 0755)
	} else {
		fmt.Println("Initial directory already exists...")
	}

	dateBackup := dir + time.Now().Format("2006-January-02")

	if !itExists(dateBackup) {
		fmt.Println("This date's backup does't exists\n" +
			"Making this date's backup directory...")

		err := os.Mkdir(dateBackup, 0755)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		fmt.Println("This date's directory exists...\n")
	}

	countries := getCountries()

	dateBackup += "/"

	fmt.Println("Saving individual country data...\n")
	for i := 0; i < len(countries); i++ {
		countryFilePath := dateBackup + countries[i].Slug
		if itExists(countryFilePath + ".txt") {
			fmt.Printf("(%d/%d) %s's data already exists...\n", i+1, len(countries), countries[i].Country)
			continue
		}

		fmt.Printf("(%d/%d) Saving %s's data...\n", i+1, len(countries), countries[i].Country)
		responseString, err := FetchApiString("https://api.covid19api.com/country/" + countries[i].Slug)
		if err != nil {
			fmt.Println("Encountered error...\nRetrying...")
			i--
			continue
		}

		file, err2 := os.Create(dateBackup + countries[i].Slug + ".txt")
		if err2 != nil {
			fmt.Println("Encountered error...\nRetrying...")
			i--
			continue
		}
		defer file.Close()

		_, err3 := file.WriteString(responseString)
		if err3 != nil {
			fmt.Println("Encountered error...\nRetrying...")
			i--
			continue
		}

	}
	fmt.Println("All COVID19 data saved successfully...")
}
