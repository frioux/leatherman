package sweetmarias // import "github.com/frioux/leatherman/pkg/sweetmarias"

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/frioux/leatherman/internal/lmhttp"
)

// Coffee has all the details of a Sweet Marias coffee.
type Coffee struct {
	Title    string
	Overview string
	Score    float32
	URL      string
	SKU      string

	FarmNotes, CuppingNotes string

	Images               []string
	AdditionalAttributes map[string]string
}

func extractImages(doc *goquery.Document) ([]string, error) {
	var (
		err    error
		images []string
	)
	doc.Find(`script[type="text/x-magento-init"]`).Each(func(_ int, s *goquery.Selection) {
		t := s.Text()
		if !strings.Contains(s.Text(), "mage/gallery/gallery-ext") {
			return
		}
		type imageContainer struct {
			A struct {
				B struct {
					Data []struct {
						Full string `json:"full"`
					}
				} `json:"mage/gallery/gallery-ext"`
			} `json:"[data-gallery-role=gallery-placeholder]"`
		}
		var container imageContainer
		if err = json.Unmarshal([]byte(t), &container); err != nil {
			return
		}

		images = make([]string, len(container.A.B.Data))
		for i, d := range container.A.B.Data {
			images[i] = d.Full
		}
	})
	if err != nil {
		return nil, fmt.Errorf("parsing images json: %w", err)
	}

	return images, nil
}

func extractSKU(doc *goquery.Document) (string, error) {
	var (
		err error
		sku string
	)
	doc.Find(`script[type="text/x-magento-init"]`).Each(func(_ int, s *goquery.Selection) {
		t := s.Text()
		if !strings.Contains(s.Text(), "view_sku") {
			return
		}
		type skuContainer struct {
			A struct {
				B struct {
					Handles []string `json:"handles"`
				} `json:"pageCache"`
			} `json:"body"`
		}
		var container skuContainer
		if err = json.Unmarshal([]byte(t), &container); err != nil {
			return
		}

		for _, h := range container.A.B.Handles {
			const prefix = "catalog_product_view_sku_"
			if strings.HasPrefix(h, prefix) {
				sku = strings.TrimPrefix(h, prefix)
				break
			}
		}
	})
	if err != nil {
		return "", fmt.Errorf("parsing sku json: %w", err)
	}

	return sku, nil

}

// LoadCoffee loads a Coffee from the passed url.
func LoadCoffee(ctx context.Context, url string) (Coffee, error) {
	res, err := lmhttp.Get(ctx, url)
	if err != nil {
		return Coffee{}, fmt.Errorf("http.Get: %w", err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return Coffee{}, fmt.Errorf("goquery.NewDocumentFromReader: %w", err)
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

	if scoreStr := doc.Find("h5.score-value").Text(); scoreStr != "" {
		score, err := strconv.ParseFloat(scoreStr, 32)
		if err != nil {
			return Coffee{}, fmt.Errorf("strconv.ParseFloat: %w", err)
		}
		c.Score = float32(score)
	}

	c.AdditionalAttributes = map[string]string{}

	doc.Find("table.additional-attributes-table td").
		Each(func(_ int, s *goquery.Selection) {
			header, _ := s.Attr("data-th")
			c.AdditionalAttributes[header] = strings.Trim(s.Text(), " \n\t")
		})

	c.Images, err = extractImages(doc)
	if err != nil {
		return Coffee{}, err
	}

	c.SKU, err = extractSKU(doc)
	if err != nil {
		return Coffee{}, err
	}

	return c, nil
}
