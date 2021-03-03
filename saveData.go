package main

import (
	"errors"
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	saveData()
}

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

func removeOldData(files []fs.FileInfo) {
	dir := "./date-backups/"
	presentData := time.Now().Format("2006-January-02")

	for _, file := range files {
		if file.Name() != presentData && file.Name() != "countries.txt" {
			fmt.Println("Removing data for " + file.Name() + "...")
			os.RemoveAll(dir + file.Name())
		}
	}
}

func saveData() {
	dir := "./date-backups/"

	if !itExists(dir) {
		fmt.Println("Making initial directory...")
		os.Mkdir(dir, 0777)
	} else {
		fmt.Println("Initial directory already exists...")
	}

	currentDate := time.Now().Format("2006-January-02")
	dateBackup := dir + currentDate

	if !itExists(dateBackup) {
		fmt.Println("This date's backup does't exists\n" +
			"Making this date's backup directory...")

		dateBackup = dir + "saving"
		if !itExists(dateBackup) {
			err := os.Mkdir(dateBackup, 0777)
			if err != nil {
				log.Fatal(err)
			}
		}

	} else {
		fmt.Println("This date's directory exists...\n")
	}

	countries := getCountries()

	dateBackupDirectory := dateBackup + "/"

	fmt.Println("Saving individual country data...\n")
	for i := 0; i < len(countries); i++ {

		if countries[i].Slug == "united-states" {
			continue
		}

		countryFilePath := dateBackupDirectory + countries[i].Slug
		if itExists(countryFilePath + ".txt") {
			fmt.Printf("(%d/%d) %s's data already exists...\n", i+1, len(countries), countries[i].Country)
			continue
		}

		fmt.Printf("(%d/%d) Saving %s's data...\n", i+1, len(countries), countries[i].Country)
		responseString, err := FetchApiString("https://api.covid19api.com/country/" + countries[i].Slug)
		if err != nil {
			fmt.Println("Encountered error fetching...\nSleeping for 5 seconds...")
			time.Sleep(5 * time.Second)
			i--
			continue
		}

		if len(responseString) < 1000 {
			fmt.Println("Encountered error fetching...\nSleeping for 5 seconds...")
			time.Sleep(5 * time.Second)
			i--
			continue
		}

		file, err2 := os.Create(dateBackupDirectory + countries[i].Slug + ".txt")
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
	fmt.Println("All COVID19 data saved successfully\n")

	fmt.Println("Changing direcotory's name to current date...")
	os.Rename(dateBackup, dir+currentDate)

	fmt.Println("Checking for old data...")
	files, _ := ioutil.ReadDir("./date-backups/")
	if len(files) > 1 {
		removeOldData(files)
		fmt.Println("All old data removed\n")
	} else {
		fmt.Println("No old data present\n")
	}

}
