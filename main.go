package main

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

var googleDomains = map[string]string{
	"com": "https://google.com/search?q=",
	"in":  "https://google.co.in/search?q=",
}

var userAgents = []string{}

type SerachResults struct {
	ResultRank  int
	ResultUrl   string
	ResultTitle string
	ResultDesc  string
}

func randomUserAgent() string {
	rand.Seed(time.Now().Unix())
	randNum := rand.Int() % len(userAgents)
	return userAgents[randNum]
}

func buildGoogleUrls(searchTerm string, countryCode string, languageCode string, pages int, resultCount int) ([]string, error) {
	toScrape := []string{}
	//refining the search query before before processing
	searchTerm = strings.Trim(searchTerm, " ")
	searchTerm = strings.Replace(searchTerm, " ", "+", -1)
	if googleBaseUrl, found := googleDomains[countryCode]; found {
		for i := 0; i < pages; i++ {
			start := i * resultCount
			urlQuery := fmt.Sprintf("%s%s&num=%d&hl=%s&start=%d&filter=0", googleBaseUrl, searchTerm, start, languageCode, pages)
			toScrape = append(toScrape, urlQuery)
		}
	} else {
		error := fmt.Errorf("country %s is currently not supported", countryCode)
		return nil, error
	}

	return toScrape, nil
}

func makeRequestToGoogle(page string, proxy interface{}) {

}

func googleResultParsing(response string, counter int) ([]SerachResults, error) {
	return []SerachResults,nil
}

func searchOnGoogle(searchTerm string, countryCode string, languageCode string, pages int, resultCount int, proxy interface{}, backOffTime int) ([]SerachResults, error) {
	results := []SerachResults{}
	resultCounter := 0
	googlePages, err := buildGoogleUrls(searchTerm, countryCode, languageCode, pages, resultCount)
	if err != nil {
		return nil, err
	}
	for _, page := range googlePages {
		res, err := makeRequestToGoogle(page, proxy)
		if err != nil {
			return nil, err
		}
		data, err := googleResultParsing(res, resultCounter)
		if err != nil {
			return nil, err
		}
		resultCounter += len(data)
		for _, result := range data {
			results = append(results, result)
		}
		//wait for making next request to google API
		time.Sleep(time.Duration(backOffTime) * time.Second)

	}

	return results, nil
}

func main() {
	res, err := searchOnGoogle("danish sharma", "com", "en", 1, 30, nil, 5)
	if err != nil {
		panic("unable to fetch results.")
	}
	for _, r := range res {
		fmt.Println(r)
	}
}
