package main

import (
	"fmt"
	"github.com/bostigger/google-web-scraper/controller"
)

func main() {
	fmt.Println("Google web scrapper cli")
	scrapedData, err := controller.GoogleScraper("bill gates", "com", "en", 1, 10, 10)
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, res := range scrapedData {
		println(res.ResultTitle)
		println(res.ResultDesc)
	}
}
