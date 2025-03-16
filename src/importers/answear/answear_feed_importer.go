package answear

import (
	"context"
	"ecommerce-agent/src/db"
	"fmt"
	"path/filepath"
	"runtime"
)

func ImportAnswearFeed() {
	ctx := context.Background()
	client := db.GetDb()
	_, filename, _, ok := runtime.Caller(0)

	if !ok {
		fmt.Println("Nie można określić ścieżki bieżącego pliku")
		return
	}

	baseDir := filepath.Dir(filename)

	xmlPath := filepath.Join(baseDir, "assets", "products.xml")

	fmt.Printf("Ścieżka do pliku XML: %s\n", xmlPath)

	if err := CreateProductsClass(ctx, client); err != nil {
		fmt.Printf("Błąd podczas tworzenia klasy: %v\n", err)
	}

	ReadXml(xmlPath, func(product Product) {
		exists, err := ProductExist(client, product.ID)

		if err != nil {
		}

		if !exists {
			InsertNewProduct(ctx, client, product)
			return
		}

		fmt.Printf("Produkt %d już istnieje w Weaviate\n", product.ID)
	})
}
