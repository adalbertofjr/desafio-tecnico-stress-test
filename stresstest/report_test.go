package stresstest

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"
	"time"
)

func TestPrintReport(t *testing.T) {
	tests := []struct {
		name     string
		report   Report
		expected []string
	}{
		{
			name: "Relatório com apenas sucessos",
			report: Report{
				url:             "http://example.com",
				TotalRequests:   100,
				RequestsExecs:   100,
				SuccessRequests: 100,
				FailureRequests: map[int]int{},
				TimeDuration:    5 * time.Second,
			},
			expected: []string{
				"Relatório de Stress Test",
				"URL Testada: http://example.com",
				"Total de Requisições: 100",
				"Requisições Executadas: 100",
				"Duração Total: 5.00s",
				"Código (200): 100 Sucesso",
			},
		},
		{
			name: "Relatório com sucessos e falhas",
			report: Report{
				url:             "http://test.com",
				TotalRequests:   200,
				RequestsExecs:   200,
				SuccessRequests: 150,
				FailureRequests: map[int]int{404: 30, 500: 20},
				TimeDuration:    10500 * time.Millisecond,
			},
			expected: []string{
				"Relatório de Stress Test",
				"URL Testada: http://test.com",
				"Total de Requisições: 200",
				"Requisições Executadas: 200",
				"Duração Total: 10.50s",
				"Código (200): 150 Sucesso",
			},
		},
		{
			name: "Relatório vazio",
			report: Report{
				url:             "http://empty.com",
				TotalRequests:   0,
				RequestsExecs:   0,
				SuccessRequests: 0,
				FailureRequests: map[int]int{},
				TimeDuration:    0,
			},
			expected: []string{
				"Relatório de Stress Test",
				"URL Testada: http://empty.com",
				"Total de Requisições: 0",
				"Requisições Executadas: 0",
				"Duração Total: 0.00s",
				"Código (200): 0 Sucesso",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Capturar stdout
			old := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			PrintReport(tt.report)

			// Restaurar stdout
			w.Close()
			os.Stdout = old

			var buf bytes.Buffer
			io.Copy(&buf, r)
			output := buf.String()

			// Verificar se todas as strings esperadas estão na saída
			for _, expected := range tt.expected {
				if !strings.Contains(output, expected) {
					t.Errorf("Saída não contém '%s'\nSaída completa:\n%s", expected, output)
				}
			}

			// Verificar códigos de falha específicos
			for code := range tt.report.FailureRequests {
				expectedCode := fmt.Sprintf("Código (%d):", code)
				if !strings.Contains(output, expectedCode) {
					t.Errorf("Saída não contém código de falha %d", code)
				}
			}
		})
	}
}

func TestReportStruct(t *testing.T) {
	// Testar criação e campos da struct Report
	report := Report{
		url:             "http://localhost:8080",
		TotalRequests:   50,
		RequestsExecs:   50,
		SuccessRequests: 45,
		FailureRequests: map[int]int{404: 5},
		TimeDuration:    2 * time.Second,
	}

	if report.url != "http://localhost:8080" {
		t.Errorf("URL incorreta: %s", report.url)
	}

	if report.TotalRequests != 50 {
		t.Errorf("TotalRequests incorreto: %d", report.TotalRequests)
	}

	if report.RequestsExecs != 50 {
		t.Errorf("RequestsExecs incorreto: %d", report.RequestsExecs)
	}

	if report.SuccessRequests != 45 {
		t.Errorf("SuccessRequests incorreto: %d", report.SuccessRequests)
	}

	if report.FailureRequests[404] != 5 {
		t.Errorf("FailureRequests[404] incorreto: %d", report.FailureRequests[404])
	}

	if report.TimeDuration != 2*time.Second {
		t.Errorf("TimeDuration incorreto: %v", report.TimeDuration)
	}
}

func TestPrintReportWithMultipleFailureCodes(t *testing.T) {
	// Capturar stdout
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	report := Report{
		url:             "http://multi.com",
		TotalRequests:   100,
		RequestsExecs:   100,
		SuccessRequests: 70,
		FailureRequests: map[int]int{
			404: 10,
			500: 15,
			503: 5,
		},
		TimeDuration: 3 * time.Second,
	}

	PrintReport(report)

	// Restaurar stdout
	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	io.Copy(&buf, r)
	output := buf.String()

	// Verificar que todos os códigos de erro aparecem
	failureCodes := []int{404, 500, 503}
	for _, code := range failureCodes {
		expectedCode := fmt.Sprintf("Código (%d):", code)
		if !strings.Contains(output, expectedCode) {
			t.Errorf("Saída não contém código de falha %d", code)
		}
	}

	// Verificar estrutura básica
	if !strings.Contains(output, "Relatório de Stress Test") {
		t.Error("Saída não contém título do relatório")
	}

	if !strings.Contains(output, "Código (200): 70 Sucesso") {
		t.Error("Saída não contém contagem de sucessos")
	}
}
