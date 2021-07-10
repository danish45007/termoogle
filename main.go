package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
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

func getScrapeClient(proxy interface{}) *http.Client {
	// if the proxy type is string
	//return custom http client with proxy enabled
	switch v := proxy.(type) {
	case string:
		proxyUrl, _ := url.Parse(v)
		return &http.Client{Transport: &http.Transport{
			Proxy: http.ProxyURL(proxyUrl),
		}}
	default:
		return &http.Client{}
	}
}

func makeRequestToGoogle(searchUrl string, proxy interface{}) (*http.Response, error) {
	baseClient := getScrapeClient(proxy)
	req, _ := http.NewRequest("GET", searchUrl, nil)
	req.Header.Set("User-Agent", randomUserAgent())

	res, err := baseClient.Do(req)
	if res.StatusCode != 200 {
		err := fmt.Errorf("we recevied a non-200 status code suggesting a ban")
		return nil, err
	}
	if err != nil {
		return nil, err
	}

	return res, nil
}

func googleLinkValidation(url string) bool {
	if url != "" && url != "#" && !strings.HasPrefix(url, "/") {
		return true
	} else {
		return false
	}

}

func googleResultParsing(response *http.Response, rank int) ([]SerachResults, error) {
	doc, err := goquery.NewDocumentFromResponse(response)
	if err != nil {
		return nil, err
	}
	results := []SerachResults{}
	targetTag := doc.Find("div.g")
	//increment rank
	rank++
	for i := range targetTag.Nodes {
		item := targetTag.Eq(i)
		linkTag := item.Find("a")
		link, _ := linkTag.Attr("href")
		titleTag := item.Find("h3.r")
		descTag := item.Find("span.st")
		desc := descTag.Text()
		title := titleTag.Text()
		validLink := strings.Trim(link, " ")
		// data validation
		if googleLinkValidation(validLink) {
			result := SerachResults{
				rank,
				validLink,
				title,
				desc,
			}
			results = append(results, result)
		}
	}
	return results, nil
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

		results = append(results, data...)

		//wait for making next request to google API
		time.Sleep(time.Duration(backOffTime) * time.Second)

	}

	return results, nil
}

func main() {
	res, err := searchOnGoogle("danish sharma", "com", "en", 1, 30, nil, 5)
	if err != nil {
		panic("Unable to fetch results.")
	}
	for _, r := range res {
		fmt.Println(r)
	}
}
