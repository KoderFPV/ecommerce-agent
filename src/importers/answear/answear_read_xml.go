package answear

import (
	"encoding/xml"
	"fmt"
	"os"
)

type Product struct {
	Material                           string           `xml:"material"`
	ProductBrand                       string           `xml:"productBrand"`
	AttributesJson                     []string         `xml:"attributesJson"`
	ColorDataJson                      string           `xml:"colorDataJson"`
	SizeJson                           string           `xml:"sizeJson"`
	OnOfferFor30Days                   int              `xml:"onOfferFor30Days"`
	PercentageDiscountValueFromMinimal string           `xml:"percentageDiscountValueFromMinimal"`
	URL                                string           `xml:"url"`
	ID                                 int              `xml:"id"`
	Name                               string           `xml:"name"`
	Subtitle                           string           `xml:"subtitle"`
	Slug                               string           `xml:"slug"`
	FrontendUuid                       string           `xml:"frontendUuid"`
	PbbCode                            string           `xml:"pbbCode"`
	Sport                              int              `xml:"sport"`
	OfferType                          string           `xml:"offerType"`
	Mpn                                string           `xml:"mpn"`
	Price                              string           `xml:"price"`
	PriceRegular                       string           `xml:"priceRegular"`
	PriceMinimal                       string           `xml:"priceMinimal"`
	PriceIsDiscounted                  int              `xml:"priceIsDiscounted"`
	DiscountValue                      float64          `xml:"discountValue"`
	Sale                               int              `xml:"sale"`
	Description                        string           `xml:"description"`
	ExtendedDescription                string           `xml:"extendedDescription"`
	ImageLink                          string           `xml:"imageLink"`
	AdditionalImageLink                []string         `xml:"additionalImageLink"`
	Category                           string           `xml:"category"`
	Sex                                string           `xml:"sex"`
	Categories                         []string         `xml:"categories>item"`
	Pseudocategory                     []Pseudocategory `xml:"pseudocategory>item"`
	ColorData                          ColorData        `xml:"colorData"`
	Attributes                         []Attribute      `xml:"attributes>item"`
	Variations                         string           `xml:"variations"`
	AllVariations                      []Variation      `xml:"allVariations>item"`
	Availability                       string           `xml:"availability"`
	ActiveFrom                         string           `xml:"activeFrom"`
	Sorting                            Sorting          `xml:"sorting"`
	Base64Image                        string
}

type Pseudocategory struct {
	ID   int    `xml:"id"`
	Name string `xml:"name"`
}

type ColorData struct {
	ID        int    `xml:"id"`
	Name      string `xml:"name"`
	FrontName string `xml:"frontName"`
	Code      string `xml:"code"`
	Hex       string `xml:"hex"`
	Image     string `xml:"image"`
}

type Attribute struct {
	ID    int    `xml:"id"`
	Name  string `xml:"name"`
	Value string `xml:"value"`
}

type Variation struct {
	ID           int    `xml:"id"`
	Ean          string `xml:"ean"`
	SizeId       int    `xml:"sizeId"`
	SizeName     string `xml:"sizeName"`
	Availability string `xml:"availability"`
}

type Sorting struct {
	Global int `xml:"GLOBAL"`
}

func ReadXml(filePath string, onProductRead func(Product)) {
	file, err := os.Open(filePath)

	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		return
	}

	defer file.Close()

	decoder := xml.NewDecoder(file)

	for {
		tok, err := decoder.Token()

		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			fmt.Printf("Error decoding token: %v\n", err)
			return
		}

		switch se := tok.(type) {

		case xml.StartElement:
			if se.Name.Local == "item" {
				var product Product

				err := decoder.DecodeElement(&product, &se)

				if err != nil {
					fmt.Printf("Error decoding element: %v\n", err)
					return
				}

				onProductRead(product)
			}
		}
	}
}
