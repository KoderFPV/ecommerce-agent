package main

import (
	"fmt"

	"ecommerce-agent/src/importers/answear"
)

func main() {
	fmt.Println("Rozpoczynam import produktów Answear...")
	answear.ImportAnswearFeed()
	fmt.Println("Import zakończony.")
}
