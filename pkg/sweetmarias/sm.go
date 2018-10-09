package sweetmarias // import "github.com/frioux/leatherman/pkg/sweetmarias"

import (
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/pkg/errors"
)

// Coffee has all the details of a Sweet Marias coffee.
type Coffee struct {
	Title    string
	Overview string
	Score    float32
	URL      string

	FarmNotes, CuppingNotes string

	AdditionalAttributes map[string]string
}

// LoadCoffee loads a Coffee from the passed url.
func LoadCoffee(url string) (Coffee, error) {
	res, err := http.Get(url)
	if err != nil {
		return Coffee{}, errors.Wrap(err, "http.Get")
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return Coffee{}, errors.Wrap(err, "goquery.NewDocumentFromReader")
	}

	c := Coffee{URL: url}

	c.Overview = doc.Find("div.overview p").Text()
	if c.Overview == "" {
		c.Overview = doc.Find("div.overview div.value").Text()
	}

	c.Title = doc.Find("h1.page-title span").Text()

	c.FarmNotes = doc.Find("div.origin-notes span").Text()
	if c.FarmNotes == "" {
		c.FarmNotes = doc.Find("div.origin-notes p").Text()
	}

	c.CuppingNotes = doc.Find("div.cupping-notes span").Text()
	if c.CuppingNotes == "" {
		c.CuppingNotes = doc.Find("div.cupping-notes p").Text()
	}

	score, err := strconv.ParseFloat(doc.Find("h5.score-value").Text(), 32)
	if err != nil {
		return Coffee{}, errors.Wrap(err, "strconv.ParseFloat")
	}
	c.Score = float32(score)

	c.AdditionalAttributes = map[string]string{}

	doc.Find("table.additional-attributes-table td").
		Each(func(_ int, s *goquery.Selection) {
			header, _ := s.Attr("data-th")
			c.AdditionalAttributes[header] = strings.Trim(s.Text(), " \n\t")
		})

	return c, nil
}
