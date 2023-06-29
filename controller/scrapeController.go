package controller

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/bostigger/google-web-scraper/helpers"
	"github.com/bostigger/google-web-scraper/model"
	"net/http"
	"strings"
	"time"
)

func buildGoogleUrls(searchTerm string, countryCode string, languageCode string, pages int, count int) ([]string, error) {
	var toScrape []string
	searchTerm = strings.Trim(searchTerm, "")
	searchTerm = strings.Replace(searchTerm, " ", "+", -1)
	if baseUrl, exists := helpers.DOMAINS[countryCode]; exists {
		for i := 0; i < pages; i++ {
			scrapeUrl := fmt.Sprintf("%s%s&num=%d&hl=%s&start=1&filter=0", baseUrl, searchTerm, count, languageCode)
			toScrape = append(toScrape, scrapeUrl)
		}
	} else {
		err := fmt.Errorf("country code %s not supported currently", countryCode)
		return nil, err
	}
	return toScrape, nil

}

func scrapeClientRequest(url string) (*http.Response, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", helpers.PickRandomUserAgent())
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	fmt.Println(res.Status)
	if res.StatusCode != 200 {
		err := fmt.Errorf("scrapping gone wrong with not 200 response")
		return nil, err
	}
	return res, nil
}

func googleResultParsing(response *http.Response, rank int) ([]model.SearchResult, error) {
	doc, err := goquery.NewDocumentFromResponse(response)
	if err != nil {
		return nil, err
	}
	var results []model.SearchResult
	sel := doc.Find("div.g")
	rank++
	for i := range sel.Nodes {
		item := sel.Eq(i)
		linkTag := item.Find("a")
		link, _ := linkTag.Attr("href")
		titleTag := item.Find("h3.r")
		descTag := item.Find("span.st")
		desc := descTag.Text()
		title := titleTag.Text()
		link = strings.Trim(link, " ")

		if link != "" && link != "#" && !strings.HasPrefix(link, "/") {
			result := model.SearchResult{
				ResultRank:  rank,
				ResultTitle: title,
				ResultDesc:  desc,
				ResultUrl:   link,
			}
			results = append(results, result)
			rank++
		}
	}
	return results, nil
}
func GoogleScraper(searchKeyword string, countryCode string, languageCode string, pages int, count int, waitTime int) ([]model.SearchResult, error) {
	var results []model.SearchResult
	resultCounter := 0
	googlePages, err := buildGoogleUrls(searchKeyword, countryCode, languageCode, pages, count)
	if err != nil {
		return nil, err
	}
	for _, page := range googlePages {
		res, err := scrapeClientRequest(page)
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
		time.Sleep(time.Duration(waitTime) * time.Second)
	}
	return results, nil
}
