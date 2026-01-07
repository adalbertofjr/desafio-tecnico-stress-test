package stresstest

import (
	"fmt"
	"time"
)

type Report struct {
	url             string
	TotalRequests   int
	RequestsExecs   int
	SuccessRequests int
	FailureRequests map[int]int
	TimeDuration    time.Duration
}

func PrintReport(data Report) {
	fmt.Println("=========================================")
	fmt.Println("Relatório de Stress Test")
	fmt.Println("=========================================")
	fmt.Printf("URL Testada: %s\n", data.url)
	fmt.Printf("Total de Requisições: %d\n", data.TotalRequests)
	fmt.Printf("Requisições Executadas: %d\n", data.RequestsExecs)
	fmt.Printf("Duração Total: %.2fs\n", data.TimeDuration.Seconds())
	fmt.Println("=========================================")
	fmt.Println("Detalhes das Respostas:")
	fmt.Printf("  Código (200): %d Sucesso\n", data.SuccessRequests)
	for code, count := range data.FailureRequests {
		fmt.Printf("  Código (%d): %d Falha\n", code, count)
	}
	fmt.Println("=========================================")
}
