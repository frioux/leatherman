package sweetmarias // import "github.com/frioux/leatherman/pkg/sweetmarias"

import (
	"context"
	"errors"
	"fmt"
	"io"
	"math/rand"

	"github.com/PuerkitoBio/goquery"
	"github.com/frioux/leatherman/internal/lmhttp"
)

var errStatusCode = errors.New("status code error")

// AllCoffees returns a list of URLs for each coffee in Sweet Maria's inventory.
func AllCoffees() ([]string, error) {
	const allURL = "https://www.sweetmarias.com/green-coffee.html?product_list_limit=all&sm_status=1"

	res, err := lmhttp.Get(context.TODO(), allURL)
	if err != nil {
		return nil, fmt.Errorf("http.Get: %w", err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return nil, fmt.Errorf("%d %s: %w", res.StatusCode, res.Status, errStatusCode)
	}

	return allCoffees(res.Body)
}

func allCoffees(r io.Reader) ([]string, error) {
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return nil, fmt.Errorf("goquery.NewDocumentFromReader: %w", err)
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
