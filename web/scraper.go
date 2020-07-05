package scraper

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

//Structure contains the structure of the scrape
type Structure struct {
	Tittle string
	Desc   string
}

var shortURL string

//URL prefix map
var URL = map[string]string{
	"fill":   "https://www.animefillerlist.com/shows/",
	"define": "https://www.urbandictionary.com/define.php?term=",
	"news":   "https://www.animenewsnetwork.com/",
}

func buildURL(searchTerm string, shorturl string) string {
	shortURL = shorturl
	//filler
	if shortURL == "fill" {
		searchTerm = strings.Replace(searchTerm, " ", "-", -1)
		return fmt.Sprintf("%s%s", URL[shortURL], searchTerm)
	}
	//"define"
	searchTerm = strings.Replace(searchTerm, " ", "%20", -1)
	return fmt.Sprintf("%s%s", URL[shortURL], searchTerm)
}

func request(searchURL string) (*http.Response, error) {
	baseClient := &http.Client{}
	req, _ := http.NewRequest("GET", searchURL, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.100 Safari/537.36")
	res, err := baseClient.Do(req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func resultParser(response *http.Response, times int) ([]Structure, error) {
	doc, err := goquery.NewDocumentFromResponse(response)
	if err != nil {
		return nil, err
	}
	results := []Structure{}
	//if statement can change what to scrape depending on website chosen (fill in this case)
	if shortURL == "fill" {
		var item = doc.Find("div.CopyProtected")

		classes := [4]string{"div.manga_canon", "div.mixed_canon/filler", "div.filler", "div.anime_canon"}

		for _, search := range classes {

			selection := item.Find(search)

			episodes := selection.Find("span.Episodes")
			episodestxt := episodes.Text()

			if episodestxt == "" {
				episodestxt = "0"
			}
			result := Structure{
				"%//DEVIDER//%",
				episodestxt,
			}
			results = append(results, result)
		}

	}
	return results, err
}

//Scrape scrapes on the web for defined websites
func Scrape(searchTerm string, shortURL string, amount int) ([]Structure, error) {
	var url string
	if shortURL != "news" {
		url = buildURL(searchTerm, shortURL)
	} else {
		url = shortURL
	}
	res, err := request(url)
	if err != nil {
		return nil, err
	}
	Scrapes, err := resultParser(res, amount)
	if err != nil {
		return nil, err
	}
	return Scrapes, nil
}
