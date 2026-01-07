package main

import (
	"adalbertofjr/desafio-tecnico-stress-test/stresstest"
	"flag"
	"fmt"
)

func main() {
	url := flag.String("url", "", "The URL to stress test")
	requests := flag.Int("requests", 0, "Number of requests to perform")
	concurrency := flag.Int("concurrency", 0, "Number of multiple requests to make at a time")

	flag.Parse()

	// Validação
	if *url == "" {
		fmt.Println("Erro: URL é obrigatória")
		flag.Usage()
		return
	}

	if *requests <= 0 || *concurrency <= 0 {
		fmt.Println("Erro: requests e concurrency devem ser maiores que zero")
		flag.Usage()
		return
	}

	stresstest.Execute(*url, *requests, *concurrency)
}
