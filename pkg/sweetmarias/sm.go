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

	var imageErr error
	doc.Find(`script[type="text/x-magento-init"]`).
		Each(func(_ int, s *goquery.Selection) {
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
			imageErr = json.Unmarshal([]byte(t), &container)

			if imageErr != nil {
				return
			}

			c.Images = make([]string, len(container.A.B.Data))
			for i, d := range container.A.B.Data {
				c.Images[i] = d.Full
			}
		})
	if imageErr != nil {
		return Coffee{}, fmt.Errorf("parsing images json: %w", imageErr)
	}

	var skuErr error
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
		skuErr = json.Unmarshal([]byte(t), &container)
		if skuErr != nil {
			return
		}

		for _, h := range container.A.B.Handles {
			const prefix = "catalog_product_view_sku_"
			if strings.HasPrefix(h, prefix) {
				c.SKU = strings.TrimPrefix(h, prefix)
				break
			}
		}
	})
	if skuErr != nil {
		return Coffee{}, fmt.Errorf("parsing sku json: %w", imageErr)
	}

	return c, nil
}
