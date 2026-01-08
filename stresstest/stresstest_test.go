package stresstest

import (
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
)

func TestIncrementCount(t *testing.T) {
	mu := sync.Mutex{}
	counter := 0

	// Testar incremento sequencial
	result := incrementCount(&mu, &counter)
	if result != 1 || counter != 1 {
		t.Errorf("Esperado 1, obtido %d", result)
	}

	result = incrementCount(&mu, &counter)
	if result != 2 || counter != 2 {
		t.Errorf("Esperado 2, obtido %d", result)
	}
}

func TestIncrementCountConcurrency(t *testing.T) {
	mu := sync.Mutex{}
	counter := 0
	wg := sync.WaitGroup{}
	iterations := 1000

	// Testar incremento concorrente
	for range iterations {
		wg.Add(1)
		go func() {
			defer wg.Done()
			incrementCount(&mu, &counter)
		}()
	}

	wg.Wait()

	if counter != iterations {
		t.Errorf("Esperado %d, obtido %d", iterations, counter)
	}
}

func TestIncrementaStatusCount(t *testing.T) {
	tests := []struct {
		name            string
		statusCode      int
		expectedSuccess int
		expectedFailure map[int]int
	}{
		{
			name:            "Status 200 deve contar como sucesso",
			statusCode:      200,
			expectedSuccess: 1,
			expectedFailure: map[int]int{},
		},
		{
			name:            "Status 201 deve contar como sucesso",
			statusCode:      201,
			expectedSuccess: 1,
			expectedFailure: map[int]int{},
		},
		{
			name:            "Status 299 deve contar como sucesso",
			statusCode:      299,
			expectedSuccess: 1,
			expectedFailure: map[int]int{},
		},
		{
			name:            "Status 404 deve contar como falha",
			statusCode:      404,
			expectedSuccess: 0,
			expectedFailure: map[int]int{404: 1},
		},
		{
			name:            "Status 500 deve contar como falha",
			statusCode:      500,
			expectedSuccess: 0,
			expectedFailure: map[int]int{500: 1},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mu := sync.Mutex{}
			data := Report{
				SuccessRequests: 0,
				FailureRequests: make(map[int]int),
			}

			incrementaStatusCount(&mu, tt.statusCode, &data)

			if data.SuccessRequests != tt.expectedSuccess {
				t.Errorf("Sucesso: esperado %d, obtido %d", tt.expectedSuccess, data.SuccessRequests)
			}

			if len(data.FailureRequests) != len(tt.expectedFailure) {
				t.Errorf("Falhas: esperado %d códigos, obtido %d", len(tt.expectedFailure), len(data.FailureRequests))
			}

			for code, count := range tt.expectedFailure {
				if data.FailureRequests[code] != count {
					t.Errorf("Status %d: esperado %d, obtido %d", code, count, data.FailureRequests[code])
				}
			}
		})
	}
}

func TestMakeRequest(t *testing.T) {
	// Criar servidor HTTP de teste
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	mu := sync.Mutex{}
	wg := sync.WaitGroup{}
	data := Report{
		SuccessRequests: 0,
		RequestsExecs:   0,
		FailureRequests: make(map[int]int),
	}

	wg.Add(1)
	makeRequest(&mu, &wg, server.URL, &data)

	if data.RequestsExecs != 1 {
		t.Errorf("Esperado 1 requisição executada, obtido %d", data.RequestsExecs)
	}

	if data.SuccessRequests != 1 {
		t.Errorf("Esperado 1 requisição com sucesso, obtido %d", data.SuccessRequests)
	}
}

func TestMakeRequestWithDifferentStatusCodes(t *testing.T) {
	tests := []struct {
		name            string
		statusCode      int
		expectedSuccess int
	}{
		{"Status 200", 200, 1},
		{"Status 201", 201, 1},
		{"Status 404", 404, 0},
		{"Status 500", 500, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.statusCode)
			}))
			defer server.Close()

			mu := sync.Mutex{}
			wg := sync.WaitGroup{}
			data := Report{
				SuccessRequests: 0,
				RequestsExecs:   0,
				FailureRequests: make(map[int]int),
			}

			wg.Add(1)
			makeRequest(&mu, &wg, server.URL, &data)

			if data.SuccessRequests != tt.expectedSuccess {
				t.Errorf("Esperado %d sucesso, obtido %d", tt.expectedSuccess, data.SuccessRequests)
			}
		})
	}
}

func TestExecute(t *testing.T) {
	// Criar servidor HTTP de teste
	requestCount := 0
	mu := sync.Mutex{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		requestCount++
		mu.Unlock()
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	// Executar teste com poucos requests para ser rápido
	requests := 10
	concurrency := 3

	Execute(server.URL, requests, concurrency)

	if requestCount != requests {
		t.Errorf("Esperado %d requests ao servidor, obtido %d", requests, requestCount)
	}
}

func TestExecuteConcurrency(t *testing.T) {
	// Testar se a concorrência está sendo respeitada
	activeConcurrent := 0
	maxConcurrent := 0
	mu := sync.Mutex{}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		activeConcurrent++
		if activeConcurrent > maxConcurrent {
			maxConcurrent = activeConcurrent
		}
		mu.Unlock()

		// Simular processamento
		// time.Sleep(10 * time.Millisecond)

		mu.Lock()
		activeConcurrent--
		mu.Unlock()

		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	requests := 20
	concurrency := 5

	Execute(server.URL, requests, concurrency)

	// A concorrência máxima deve ser igual ou menor que o limite definido
	if maxConcurrent > concurrency {
		t.Errorf("Concorrência máxima de %d excedeu o limite de %d", maxConcurrent, concurrency)
	}
}
