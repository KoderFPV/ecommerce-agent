package answear

import (
	"context"
	"ecommerce-agent/src/db"
	"fmt"
	"path/filepath"
	"runtime"
	"sync"
)

const (
	batchSize = 20
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

	// Kanały do komunikacji między goroutines
	productChan := make(chan Product, batchSize)
	processedChan := make(chan Product, batchSize)
	doneChan := make(chan bool)

	// Grupa oczekiwania dla wszystkich goroutines
	var wg sync.WaitGroup

	// Uruchomienie workerów do przetwarzania obrazów
	for i := 0; i < batchSize; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for product := range productChan {
				base64Image, err := FetchImageAsBase64(product.ImageLink)
				if err != nil {
					fmt.Printf("Błąd podczas pobierania obrazu dla produktu %d: %v\n", product.ID, err)
					continue
				}
				product.Base64Image = base64Image
				processedChan <- product
			}
		}()
	}

	// Uruchomienie workera do zapisywania do bazy
	wg.Add(1)
	go func() {
		defer wg.Done()
		products := make([]Product, 0, batchSize)
		for product := range processedChan {
			products = append(products, product)
			if len(products) == batchSize {
				// Zapisz partię produktów do bazy
				for _, p := range products {
					InsertNewProduct(ctx, client, p)
				}
				products = products[:0] // Wyczyść slice
			}
		}
		// Zapisz pozostałe produkty
		for _, p := range products {
			InsertNewProduct(ctx, client, p)
		}
		doneChan <- true
	}()

	// Licznik przetworzonych produktów
	count := 0

	// Funkcja do zamykania kanałów po zakończeniu
	go func() {
		wg.Wait()
		close(processedChan)
	}()

	// Odczyt XML i wysyłanie produktów do przetwarzania
	ReadXml(xmlPath, func(product Product) {
		productChan <- product
		count++
		fmt.Printf("Wysłano produkt do przetwarzania: %d\n", count)
	})

	close(productChan)
	<-doneChan

	fmt.Printf("Zakończono import %d produktów\n", count)
}
