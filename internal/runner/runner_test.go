package runner

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
)

func TestRunExactRequestsAndStatuses(t *testing.T) {
	var n int64
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		if atomic.AddInt64(&n, 1)%4 == 0 {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	res, err := Run(context.Background(), Config{URL: srv.URL, Requests: 100, Concurrency: 7}, nil)
	if err != nil {
		t.Fatal(err)
	}
	if res.Total != 100 || atomic.LoadInt64(&n) != 100 {
		t.Fatalf("total=%d server=%d, want 100", res.Total, n)
	}
	if res.StatusCodes[200] != 75 || res.StatusCodes[500] != 25 {
		t.Fatalf("unexpected distribution: %v", res.StatusCodes)
	}
}

func TestRunConcurrencyGreaterThanRequests(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {}))
	defer srv.Close()

	res, err := Run(context.Background(), Config{URL: srv.URL, Requests: 3, Concurrency: 50}, nil)
	if err != nil || res.Total != 3 {
		t.Fatalf("res=%v err=%v", res, err)
	}
}

func TestRunConnectionErrors(t *testing.T) {
	res, _ := Run(context.Background(), Config{URL: "http://127.0.0.1:1", Requests: 5, Concurrency: 2}, nil)
	if res.Total != 5 || res.Errors != 5 {
		t.Fatalf("got %v", res)
	}
}

func TestValidate(t *testing.T) {
	if _, err := Run(context.Background(), Config{}, nil); err == nil {
		t.Fatal("expected validation error")
	}
}
