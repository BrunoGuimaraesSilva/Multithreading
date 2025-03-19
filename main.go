package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Address struct {
	Street       string `json:"logradouro"`
	Neighborhood string `json:"bairro"`
	City         string `json:"localidade"`
	State        string `json:"uf"`
	ZipCode      string `json:"cep"`
}

func main() {
	var zipCode string
	fmt.Print("Enter the ZIP code: ")
	fmt.Scan(&zipCode)

	timeout := time.Second
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	resultCh := make(chan string, 1)

	brasilapi := "https://brasilapi.com.br/api/cep/v1/" + zipCode
	viacep := "http://viacep.com.br/ws/" + zipCode + "/json/"

	go fetchZipCode(ctx, brasilapi, resultCh, "BrasilAPI")
	go fetchZipCode(ctx, viacep, resultCh, "ViaCEP")

	select {
	case result := <-resultCh:
		fmt.Println(result)
	case <-ctx.Done():
		fmt.Println("Error: Request timeout")
	}
}

func fetchZipCode(ctx context.Context, baseURL string, resultCh chan string, apiName string) {

	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, baseURL, nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		resultCh <- fmt.Sprintf("Error in %s: %v", apiName, err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		resultCh <- fmt.Sprintf("%s returned an error: %s", apiName, resp.Status)
		return
	}

	var address Address
	if err := json.NewDecoder(resp.Body).Decode(&address); err != nil {
		resultCh <- fmt.Sprintf("Error decoding %s response: %v", apiName, err)
		return
	}

	resultCh <- fmt.Sprintf("%s: %s, %s, %s - %s (%s)", apiName, address.Street, address.Neighborhood, address.City, address.State, address.ZipCode)
}
