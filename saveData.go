package main

import (
	"errors"
	"fmt"
	"io/fs"
	"io/ioutil"
	"net/http"
	"os"
	"sync"
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
	dir := "./server/date-backups/"
	presentData := time.Now().Format("2006-January-02")

	for _, file := range files {
		if file.Name() != presentData && file.Name() != "countries.txt" {
			fmt.Println("Removing data for " + file.Name() + "...")
			os.RemoveAll(dir + file.Name())
		}
	}
}

func saveData() {
	dir := "./server/date-backups/"

	if !itExists(dir) {
		fmt.Println("Making initial directory...")
		if err := os.Mkdir(dir, 0777); err != nil {
			panic(err)
		}
	} else {
		fmt.Println("Initial directory already exists...\n")
	}

	currentDate := time.Now().Format("2006-January-02")
	dateBackup := dir + currentDate

	if !itExists(dateBackup) {
		fmt.Println("This date's backup does't exists\n" +
			"Making this date's backup directory...")

		dateBackup = dir + "saving"
		if !itExists(dateBackup) {
			if err := os.Mkdir(dateBackup, 0777); err != nil {
				panic(err)
			}
		}

	} else {
		fmt.Println("This date's directory exists...")
	}

	countries := getCountries()

	dateBackupDirectory := dateBackup + "/"

	fmt.Println("\nSaving individual country data...")

	var wg sync.WaitGroup
	for i := 0; i < len(countries); {

		interval := 5
		if interval > len(countries)-i {
			interval = len(countries) - i
		}

		wg.Add(interval)
		for j := 0; j < interval; j++ {
			go saveCountry(dateBackupDirectory, countries[i][0], countries[i][1], &wg)
			i++
		}
		wg.Wait()
		fmt.Printf("Progress: (%d/%d)\n\n", i, len(countries))
		time.Sleep(time.Millisecond * 300)
	}

	fmt.Println("All COVID19 data saved successfully\n")

	fmt.Println("Changing direcotory's name to current date...")
	os.Rename(dateBackup, dir+currentDate)

	fmt.Println("Checking for old data...")
	files, _ := ioutil.ReadDir("./server/date-backups/")
	if len(files) > 2 {
		removeOldData(files)
		fmt.Println("All old data removed\n")
	} else {
		fmt.Println("No old data present\n")
	}

}

func saveCountry(dateBackupDirectory string, country string, name string, wg *sync.WaitGroup) {

	if country == "united-states" {
		wg.Done()
		return
	}

	countryFilePath := dateBackupDirectory + country
	if itExists(countryFilePath + ".txt") {
		fmt.Printf("%s's data already exists...\n", name)
		wg.Done()
		return
	}

	fmt.Printf("Saving %s's data...\n", name)
	responseString, err := FetchApiString("https://api.covid19api.com/country/" + country)
	if err != nil {
		fmt.Println("Encountered errors... Retrying in 1 second...\r")
		time.Sleep(1 * time.Second)
		saveCountry(dateBackupDirectory, country, name, wg)
		return
	}

	if len(responseString) < 2000 {
		fmt.Println("Encountered errors... Retrying in 1 second...\r")
		time.Sleep(1 * time.Second)
		saveCountry(dateBackupDirectory, country, name, wg)
		return
	}

	file, err2 := os.Create(dateBackupDirectory + country + ".txt")
	if err2 != nil {
		fmt.Println("Encountered errors... Retrying in 1 second...\r")
		saveCountry(dateBackupDirectory, country, name, wg)
		return
	}
	defer file.Close()

	if _, err3 := file.WriteString(responseString); err3 != nil {
		fmt.Println("Encountered errors... Retrying in 1 second...\r")
		saveCountry(dateBackupDirectory, country, name, wg)
		return
	}

	wg.Done()
}
