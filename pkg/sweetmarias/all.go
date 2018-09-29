package sweetmarias // import "github.com/frioux/leatherman/pkg/sweetmarias"

import (
	"fmt"
	"math/rand"
	"net/http"

	"github.com/PuerkitoBio/goquery"
	"github.com/pkg/errors"
)

const allURL = "https://www.sweetmarias.com/green-coffee.html?product_list_limit=all&sm_status=1"

var errStatusCode = errors.New("status code error")

// AllCoffees returns a list of URLs for each coffee in Sweet Maria's inventory.
func AllCoffees() ([]string, error) {
	res, err := http.Get(allURL)
	if err != nil {
		return nil, errors.Wrap(err, "http.Get")
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return nil, errors.Wrap(errStatusCode,
			fmt.Sprintf("%d %s", res.StatusCode, res.Status))
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, errors.Wrap(err, "goquery.NewDocumentFromReader")
	}

	coffees := []string{}

	doc.Find("table#table-products-list tbody tr.product a.product-item-link").Each(func(_ int, s *goquery.Selection) {
		link, ok := s.Attr("href")
		if !ok {
			return
		}
		coffees = append(coffees, link)
	})

	rand.Shuffle(len(coffees), func(i, j int) {
		coffees[i], coffees[j] = coffees[j], coffees[i]
	})

	return coffees, nil
}
