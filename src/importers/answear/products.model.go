package answear

import (
	"context"
	"fmt"
	"github.com/weaviate/weaviate-go-client/v4/weaviate"
	"github.com/weaviate/weaviate-go-client/v4/weaviate/filters"
	"github.com/weaviate/weaviate/entities/models"
)

var ProductsCollectionName = "Products"

func CreateProductsClass(ctx context.Context, client *weaviate.Client) error {
	multiModal := &models.Class{
		Class:       ProductsCollectionName,
		Description: "Sample class holding all the images",
		ModuleConfig: map[string]any{
			"multi2vec-clip": map[string]any{
				"imageFields": []string{"image"},
			},
		},
		VectorIndexType: "hnsw",
		Vectorizer:      "multi2vec-clip",
		Properties: []*models.Property{
			{
				DataType:    []string{"string"},
				Description: "The name of the file",
				Name:        "productId",
			},
			{
				DataType:    []string{"string"},
				Description: "The name of the file",
				Name:        "filename",
			},
			{
				DataType:    []string{"blob"},
				Description: "Base64 encoded image",
				Name:        "image",
			},
			{
				DataType:    []string{"string"},
				Description: "Temp field to check if save to db",
				Name:        "temp",
			},
		},
	}
	return client.Schema().ClassCreator().WithClass(multiModal).Do(ctx)
}

func InsertNewProduct(ctx context.Context, client *weaviate.Client, product Product) {
	base64Image, _err := FetchImageAsBase64(product.ImageLink)

	if _err != nil {
		fmt.Printf("Error fetching image: %v\n", _err)
		return
	}

	object := &models.Object{
		Class: ProductsCollectionName,
		Properties: map[string]any{
			"productId": product.ID,
			"image":     base64Image,
			"filename":  product.ImageLink,
			"temp":      "jestem tu uuuu, co ja tutaj robie",
		},
	}

	res, err := client.Data().Creator().WithClassName(ProductsCollectionName).WithProperties(object).Do(ctx)

	fmt.Println(res)

	if err != nil {
		fmt.Printf("Błąd podczas wstawiania produktu: %v\n", err)
	} else {
		fmt.Print("Produkt został dodany do Weaviate\n", product.ID)
	}

	return
}

func ProductExist(client *weaviate.Client, productId int) (bool, error) {
	where := filters.Where().
		WithPath([]string{"productId"}).
		WithValueText(string(productId)).
		WithOperator(filters.Equal)

	result, err := client.GraphQL().Get().
		WithClassName(ProductsCollectionName).
		WithWhere(where).
		Do(context.Background())

	if err != nil {
		return false, err
	}

	if result.Data == nil {
		return false, fmt.Errorf("brak danych w odpowiedzi")
	}

	getResult, ok := result.Data["Get"]
	if !ok || getResult == nil {
		return false, nil
	}

	classData, ok := getResult.(map[string]any)[ProductsCollectionName]
	if !ok || classData == nil {
		return false, nil
	}

	objects, ok := classData.([]any)
	if !ok {
		return false, fmt.Errorf("nieprawidłowy format danych")
	}

	return len(objects) > 0, nil
}
