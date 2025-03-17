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
				DataType:    []string{"int"},
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
				DataType:    []string{"text[]"},
				Description: "Product attributes",
				Name:        "attributes",
			},
			{
				DataType:    []string{"string"},
				Description: "Product material",
				Name:        "material",
			},
			{
				DataType:    []string{"string"},
				Description: "Product brand",
				Name:        "productBrand",
			},
			{
				DataType:    []string{"text"},
				Description: "Combined product information for search",
				Name:        "prompt",
			},
			{
				DataType:    []string{"string"},
				Description: "Product category",
				Name:        "category",
			},
			{
				DataType:    []string{"string"},
				Description: "Product price",
				Name:        "price",
			},
			{
				DataType:    []string{"string"},
				Description: "Product sex",
				Name:        "sex",
			},
		},
	}
	return client.Schema().ClassCreator().WithClass(multiModal).Do(ctx)
}

func InsertNewProduct(ctx context.Context, client *weaviate.Client, product Product) {
	attributes := make([]string, len(product.Attributes))
	for i, attr := range product.Attributes {
		attributes[i] = fmt.Sprintf("%s: %s", attr.Name, attr.Value)
	}

	prompt := fmt.Sprintf("%s %s %s %s %s %s %s %s",
		product.Name,
		product.Subtitle,
		product.ProductBrand,
		product.Description,
		product.ExtendedDescription,
		product.Category,
		product.Sex,
		product.Category,
	)

	_, err := client.Data().Creator().WithClassName(ProductsCollectionName).WithProperties(map[string]any{
		"productId":    product.ID,
		"image":        product.Base64Image,
		"filename":     product.ImageLink,
		"attributes":   attributes,
		"material":     product.Material,
		"productBrand": product.ProductBrand,
		"prompt":       prompt,
		"category":     product.Category,
		"price":        product.Price,
		"sex":          product.Sex,
	}).Do(ctx)

	if err != nil {
		fmt.Printf("Błąd podczas wstawiania produktu: %v\n", err)
	} else {
		fmt.Print("Produkt został dodany do Weaviate\n", product.ID)
	}
}

func InsertBatchProducts(ctx context.Context, client *weaviate.Client, products []Product) {
	objects := make([]*models.Object, len(products))

	for i, product := range products {
		fmt.Printf("Dodaję produkt %d do partii\n", product.ID)
		attributes := make([]string, len(product.Attributes))
		for j, attr := range product.Attributes {
			attributes[j] = fmt.Sprintf("%s: %s", attr.Name, attr.Value)
		}

		prompt := fmt.Sprintf("%s %s %s %s %s %s %s %s",
			product.Name,
			product.Subtitle,
			product.ProductBrand,
			product.Description,
			product.ExtendedDescription,
			product.Category,
			product.Sex,
			product.Category,
		)

		objects[i] = &models.Object{
			Class: ProductsCollectionName,
			Properties: map[string]any{
				"productId":    product.ID,
				"image":        product.Base64Image,
				"filename":     product.ImageLink,
				"attributes":   attributes,
				"material":     product.Material,
				"productBrand": product.ProductBrand,
				"prompt":       prompt,
				"category":     product.Category,
				"price":        product.Price,
				"sex":          product.Sex,
			},
		}
	}

	batch := client.Batch().ObjectsBatcher()
	for _, obj := range objects {
		batch.WithObject(obj)
	}

	_, err := batch.Do(ctx)
	if err != nil {
		fmt.Printf("Błąd podczas wstawiania partii produktów: %v\n", err)
	} else {
		fmt.Printf("Pomyślnie dodano partię %d produktów do Weaviate\n", len(products))
	}
}

func ProductExist(client *weaviate.Client, productId int) (bool, error) {
	where := filters.Where().
		WithPath([]string{"productId"}).
		WithValueInt(int64(productId)).
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
