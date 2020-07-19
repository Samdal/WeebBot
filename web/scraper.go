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

//mURL prefix map for the site url
//there were more options
//but they have fased out, filler
//is the only one left
var mURL = map[string]string{
	"fill": "https://www.animefillerlist.com/shows/",
}

//builds the mURL (basicly only changes " " to "-")
func buildURL(searchTerm string, shorturl string) string {
	//filler
	searchTerm = strings.Replace(searchTerm, " ", "-", -1)
	return fmt.Sprintf("%s%s", mURL[shorturl], searchTerm)
}

//Sends http GET request
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

//parser result and finds the information wanted
func resultParser(response *http.Response, shortURL string) ([]Structure, error) {
	doc, err := goquery.NewDocumentFromResponse(response)
	if err != nil {
		return nil, err
	}
	results := []Structure{}
	//if statement can change what to scrape depending on website chosen (fill in this case)
	if shortURL == "fill" {
		var item = doc.Find("div.CopyProtected")

		//the information we want
		classes := [4]string{"div.manga_canon", "div.mixed_canon/filler", "div.filler", "div.anime_canon"}

		for _, search := range classes {

			//searches for information
			selection := item.Find(search)

			//saves it as text
			episodes := selection.Find("span.Episodes")
			episodestxt := episodes.Text()

			if episodestxt == "" {
				episodestxt = "0"
			}
			//saves result in the structure struct
			//the %//DEVIDER//% is used for splitting the
			//information later
			result := Structure{
				"%//DEVIDER//%",
				episodestxt,
			}
			//append the results together
			results = append(results, result)
		}

	}
	return results, err
}

//Scrape scrapes on the web for defined websites
func Scrape(searchTerm string, shortURL string) ([]Structure, error) {

	mURL := buildURL(searchTerm, shortURL)

	response, err := request(mURL)
	if err != nil {
		return nil, err
	}

	Scrapes, err := resultParser(response, shortURL)
	if err != nil {
		return nil, err
	}

	return Scrapes, nil
}
