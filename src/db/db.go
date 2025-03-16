package db

import (
	"github.com/weaviate/weaviate-go-client/v4/weaviate"
	"sync"
)

var (
	dbConnection *weaviate.Client
	once         sync.Once
)

func GetDb() *weaviate.Client {
	once.Do(func() {
		config := weaviate.Config{
			Host:   "localhost:8080",
			Scheme: "http",
		}
		client := weaviate.New(config)
		dbConnection = client
	})
	return dbConnection
}
