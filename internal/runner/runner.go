// Package runner executes the load test.
package runner

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"
)

// Config holds the load test parameters.
type Config struct {
	URL         string
	Requests    int
	Concurrency int
	Timeout     time.Duration
}

// Result is the aggregated outcome of a run.
type Result struct {
	Duration    time.Duration
	Total       int
	StatusCodes map[int]int // HTTP status -> count
	Errors      int         // requests that failed without an HTTP response
}

// Validate checks the config.
func (c Config) Validate() error {
	if c.URL == "" {
		return errors.New("--url é obrigatório")
	}
	if c.Requests <= 0 {
		return errors.New("--requests deve ser maior que 0")
	}
	if c.Concurrency <= 0 {
		return errors.New("--concurrency deve ser maior que 0")
	}
	return nil
}

// Run performs exactly cfg.Requests requests using at most cfg.Concurrency workers.
func Run(ctx context.Context, cfg Config, client *http.Client) (Result, error) {
	if err := cfg.Validate(); err != nil {
		return Result{}, err
	}
	if client == nil {
		timeout := cfg.Timeout
		if timeout == 0 {
			timeout = 30 * time.Second
		}
		client = &http.Client{Timeout: timeout}
	}

	workers := cfg.Concurrency
	if workers > cfg.Requests {
		workers = cfg.Requests
	}

	jobs := make(chan struct{})
	// status 0 means a transport error (no HTTP response).
	statuses := make(chan int, workers)

	var wg sync.WaitGroup
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for range jobs {
				statuses <- doRequest(ctx, client, cfg.URL)
			}
		}()
	}

	start := time.Now()
	go func() {
		defer close(jobs)
		for i := 0; i < cfg.Requests; i++ {
			jobs <- struct{}{}
		}
	}()
	go func() {
		wg.Wait()
		close(statuses)
	}()

	res := Result{StatusCodes: make(map[int]int)}
	for s := range statuses {
		res.Total++
		if s == 0 {
			res.Errors++
			continue
		}
		res.StatusCodes[s]++
	}
	res.Duration = time.Since(start)
	return res, nil
}

func doRequest(ctx context.Context, client *http.Client, url string) int {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return 0
	}
	resp, err := client.Do(req)
	if err != nil {
		return 0
	}
	defer resp.Body.Close()
	// Drain is intentionally skipped for speed; connection reuse is not required here.
	return resp.StatusCode
}

// String is handy for debugging.
func (r Result) String() string {
	return fmt.Sprintf("total=%d errors=%d codes=%v duration=%s", r.Total, r.Errors, r.StatusCodes, r.Duration)
}
