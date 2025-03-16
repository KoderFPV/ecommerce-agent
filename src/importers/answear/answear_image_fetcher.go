package answear

import (
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
)

func FetchImageAsBase64(imageURL string) (string, error) {
	response, err := http.Get(imageURL)

	if err != nil {
		return "", fmt.Errorf("błąd podczas pobierania obrazu: %v", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return "", fmt.Errorf("nieprawidłowy kod odpowiedzi: %d", response.StatusCode)
	}

	imageData, err := io.ReadAll(response.Body)

	if err != nil {
		return "", fmt.Errorf("błąd podczas odczytywania danych obrazu: %v", err)
	}

	base64String := base64.StdEncoding.EncodeToString(imageData)

	return base64String, nil
}
