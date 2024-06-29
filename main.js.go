package main

import (
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"sync"
	"time"
)

var (
	targetURL   string
	requests    int
	concurrency int
)

func init() {
	flag.StringVar(&targetURL, "url", "", "URL do serviço a ser testado.")
	flag.IntVar(&requests, "requests", 1, "Número total de requests.")
	flag.IntVar(&concurrency, "concurrency", 1, "Número de chamadas simultâneas.")
	flag.Parse()
}

func main() {
	if targetURL == "" {
		fmt.Println("URL não fornecida. Use a flag --url para fornecer a URL do serviço a ser testado.")
		os.Exit(1)
	}

	_, err := url.ParseRequestURI(targetURL)
	if err != nil {
		fmt.Println("URL fornecida é inválida. Por favor, forneça uma URL válida.")
		os.Exit(1)
	}

	if requests <= 0 {
		fmt.Println("Número de requests deve ser maior que zero. Use a flag --requests para definir o número de requests.")
		os.Exit(1)
	}

	if concurrency <= 0 {
		concurrency = 1
	}

	start := time.Now()

	responses := performRequests()

	printReport(start, responses)
}

func performRequests() map[int]int {
	var wg sync.WaitGroup
	sem := make(chan bool, concurrency)

	responses := make(map[int]int)

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if len(via) >= 10 {
				return http.ErrUseLastResponse
			}
			return nil
		},
	}

	for i := 0; i < requests; i++ {
		wg.Add(1)
		go func() {
			sem <- true
			resp, err := client.Get(targetURL)
			if err != nil {
				fmt.Println(err)
				wg.Done()
				<-sem
				return
			}
			fmt.Printf("Response status: %d\n", resp.StatusCode)
			responses[resp.StatusCode]++
			resp.Body.Close()
			wg.Done()
			<-sem
		}()
	}

	wg.Wait()

	return responses
}

func printReport(start time.Time, responses map[int]int) {
	elapsed := time.Since(start).Round(time.Second)

	fmt.Println("-------------------------------------------------")
	fmt.Printf("Tempo total gasto na execução: %s\n", elapsed)
	fmt.Printf("Quantidade total de requests realizados: %d\n", requests)
	fmt.Printf("Quantidade de requests com status HTTP 200: %d\n", responses[200])
	fmt.Println("Distribuição de outros códigos de status HTTP:")
	for status, count := range responses {
		if status != 200 {
			fmt.Printf("Status %d: %d\n", status, count)
		}
	}
	fmt.Println("-------------------------------------------------")
}
