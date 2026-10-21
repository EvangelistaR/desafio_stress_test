package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/rafaelevangelista/stress-test-cli/internal/report"
	"github.com/rafaelevangelista/stress-test-cli/internal/runner"
)

func main() {
	var cfg runner.Config
	flag.StringVar(&cfg.URL, "url", "", "URL do serviço a ser testado")
	flag.IntVar(&cfg.Requests, "requests", 0, "número total de requisições")
	flag.IntVar(&cfg.Concurrency, "concurrency", 1, "número de chamadas simultâneas")
	flag.DurationVar(&cfg.Timeout, "timeout", 30*time.Second, "timeout por requisição")
	flag.Parse()

	if err := cfg.Validate(); err != nil {
		fmt.Fprintln(os.Stderr, "erro:", err)
		flag.Usage()
		os.Exit(2)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	fmt.Printf("Iniciando: %s | requests=%d | concurrency=%d\n", cfg.URL, cfg.Requests, cfg.Concurrency)
	res, err := runner.Run(ctx, cfg, nil)
	if err != nil {
		fmt.Fprintln(os.Stderr, "erro:", err)
		os.Exit(1)
	}
	report.Print(os.Stdout, res)
}
