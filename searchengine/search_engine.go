package searchengine

import (
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/danish45007/my-google.com/domains"
	"github.com/danish45007/my-google.com/urlshortner"
	"github.com/olekukonko/tablewriter"
)

var userAgents = []string{
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.100 Safari/537.36",
	"Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.100 Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.100 Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_6) AppleWebKit/604.1.38 (KHTML, like Gecko) Version/11.0 Safari/604.1.38",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:56.0) Gecko/20100101 Firefox/56.0",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13) AppleWebKit/604.1.38 (KHTML, like Gecko) Version/11.0 Safari/604.1.38",
}

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
	googleBaseUrl, found := domains.GetGoogleDomain(countryCode)
	if found && len(googleBaseUrl) > 0 {
		for i := 0; i < pages; i++ {
			start := i * resultCount
			urlQuery := fmt.Sprintf("%s%s&num=%d&hl=%s&start=%d&filter=0", googleBaseUrl, searchTerm, resultCount, languageCode, start)
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
		titleTag := item.Find("h3")
		descTag := item.Find("span")
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
			rank++
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
		for _, result := range data {
			results = append(results, result)
		}
		//wait for making next request to google API
		time.Sleep(time.Duration(backOffTime) * time.Second)

	}

	return results, nil
}

func createTableOnTerminal(searchTerm string, countryCode string, languageCode string, pages int, resultCount int, proxy interface{}, backOffTime int) {
	res, err := searchOnGoogle(searchTerm, countryCode, languageCode, pages, resultCount, proxy, backOffTime)
	if err != nil {
		fmt.Println(err)
	}
	results := [][]string{}

	for _, r := range res {
		result := []string{}
		result = append(result, strconv.Itoa(r.ResultRank), r.ResultTitle, r.ResultDesc, urlshortner.UrlShortner(r.ResultUrl))
		results = append(results, result)
	}
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Rank", "Title", "Description ", "Url"})
	table.SetBorder(true)
	table.SetCaption(true, "Search Results")
	table.SetRowLine(true)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.AppendBulk(results)

	table.Render()
}

func SearchEngine(searchTerm string, countryCode string, languageCode string, pages int, resultCount int, proxy interface{}, backOffTime int) {
	createTableOnTerminal(searchTerm, countryCode, languageCode, pages, resultCount, proxy, backOffTime)
}
