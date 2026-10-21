// Package report prints the load test summary.
package report

import (
	"fmt"
	"io"
	"sort"

	"github.com/rafaelevangelista/stress-test-cli/internal/runner"
)

// Print writes the final report to w.
func Print(w io.Writer, r runner.Result) {
	fmt.Fprintln(w, "========== Relatório do Teste de Carga ==========")
	fmt.Fprintf(w, "Tempo total gasto:            %s\n", r.Duration)
	fmt.Fprintf(w, "Total de requests realizados: %d\n", r.Total)
	fmt.Fprintf(w, "Requests com status 200:      %d\n", r.StatusCodes[200])

	fmt.Fprintln(w, "Distribuição de outros status:")
	codes := make([]int, 0, len(r.StatusCodes))
	for c := range r.StatusCodes {
		if c != 200 {
			codes = append(codes, c)
		}
	}
	sort.Ints(codes)
	if len(codes) == 0 && r.Errors == 0 {
		fmt.Fprintln(w, "  (nenhum)")
	}
	for _, c := range codes {
		fmt.Fprintf(w, "  %d: %d\n", c, r.StatusCodes[c])
	}
	if r.Errors > 0 {
		fmt.Fprintf(w, "  erros de conexão/timeout: %d\n", r.Errors)
	}
}
