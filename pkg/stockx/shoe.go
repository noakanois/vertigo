package stockx

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type ProductDetails struct {
	ID          int
	Name        string
	Subtitle    string
	LastSale    string
	ProductName string
	MainPicture string
	Attributes  map[string]string
	Description string
}

func GetShoeInformation(url string) (ProductDetails, error) {
	res, err := http.Get(url)
	if err != nil {
		return ProductDetails{}, fmt.Errorf("error fetching page: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return ProductDetails{}, fmt.Errorf("status code error: %d %s", res.StatusCode, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return ProductDetails{}, err
	}

	product := ProductDetails{
		Attributes: make(map[string]string),
	}

	primaryTitleSelection := doc.Find("h1[data-component='primary-product-title']")
	product.Name = strings.TrimSpace(primaryTitleSelection.Contents().Not("span").Text())
	product.Subtitle = strings.TrimSpace(primaryTitleSelection.Find("span[data-component='secondary-product-title']").Text())

	product.LastSale = strings.TrimSpace(doc.Find(".css-1q8ctst").First().Text())

	doc.Find(".css-17w8l66 div").Each(func(index int, item *goquery.Selection) {
		key := strings.TrimSpace(item.Find("span.chakra-text").Text())
		value := strings.TrimSpace(item.Find("p.chakra-text").Text())
		product.Attributes[key] = value
	})

	product.Description = strings.TrimSpace(doc.Find(".css-1k2nzv4").Text())

	doc.Find("img[data-image-type='360']").Each(func(i int, s *goquery.Selection) {
		srcSet, exists := s.Attr("srcset")
		if exists {
			urls := strings.Split(srcSet, ",")
			for _, url := range urls {
				trimmedUrl := strings.Split(strings.TrimSpace(url), " ")[0]
				parts := strings.Split(trimmedUrl, "/")
				for j, part := range parts {
					if part == "360" && j+2 < len(parts) {
						product.ProductName = parts[j+1]
						break
					}
				}
			}
		}
	})

	product.MainPicture = "https://images.stockx.com/images/" + product.ProductName+ "-Product.jpg"
	return product, err
}
