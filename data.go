package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

func fetchCountries() (string, error) {
	response, err := http.Get("https://api.covid19api.com/summary")
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()
	body, err2 := ioutil.ReadAll(response.Body)
	if err2 != nil {
		panic(err2)
	}
	return string(body), nil
}

type Data struct {
	Countries []struct {
		Country string `json:"Country"`
		Slug    string `json:"Slug"`
	} `json:"Countries"`
}

func getCountries() map[int][]string {
	jsonString, _ := fetchCountries()
	var data Data

	if err := json.Unmarshal([]byte(jsonString), &data); err != nil {
		panic(err)
	}

	countriesSlice := map[int][]string{}

	for index, element := range data.Countries {
		countriesSlice[index] = []string{element.Slug, element.Country}
	}

	return countriesSlice
}
