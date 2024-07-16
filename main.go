package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os/signal"
	"runtime"
	"strconv"
	"syscall"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	mux := http.NewServeMux()
	mux.HandleFunc("GET /isPrime/{number}", isPrimeHandler())
	mux.HandleFunc("GET /maxProcs", maxProcsHandler())
	mux.HandleFunc("POST /maxProcs/{number}", maxProcsHandler())

	server := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	go func() {
		slog.Info("server is starting", slog.String("addr", server.Addr))

		err := server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			panic(err)
		}
	}()

	<-ctx.Done()
	slog.Info("server is shutting down", slog.String("reason", ctx.Err().Error()))

	err := server.Shutdown(ctx)
	if err != nil {
		panic(err)
	}
}

// maxProcsHandler sets or returns the current value of GOMAXPROCS and the total number of CPUs as seen by the Go
// runtime. If the URL contains a number path value the handler will change the current GOMAXPROCS to the given number.
func maxProcsHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		numberStr := r.PathValue("number")

		if numberStr == "" {
			fmt.Fprintf(w, "Go is using %d CPUs and there are %d CPUs available\n",
				runtime.GOMAXPROCS(-1), // n <= 1 does not change the setting and returns the current number
				runtime.NumCPU(),
			)
			return
		}

		number, err := strconv.Atoi(numberStr)
		if err != nil {
			http.Error(w, fmt.Sprintf("%q is not a number", numberStr), http.StatusBadRequest)
			return
		}

		previous := runtime.GOMAXPROCS(number)
		fmt.Fprintf(w, "Go is now using %d CPUs, previously it was using %d CPUs\n", number, previous)
	}
}

// isPrimeHandler reads a number from the URL path and returns a text explaining if the number is a prime number or not.
func isPrimeHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		numberStr := r.PathValue("number")
		number, err := strconv.ParseUint(numberStr, 10, 64)
		if err != nil {
			http.Error(w, fmt.Sprintf("%q is not a number: %s", numberStr, err), http.StatusBadRequest)
			return
		}

		if isPrime(number) {
			fmt.Fprintf(w, "%d is a prime number\n", number)
		} else {
			fmt.Fprintf(w, "%d is not prime number\n", number)
		}
	}
}

// isPrime returns true if n is a prime number. The implementation is intentionally inefficient to demonstrate the
// effects of CPU limits in container environments.
func isPrime(n uint64) bool {
	if n <= 1 {
		return false
	}

	prime := true
	for i := range n {
		if i <= 1 {
			continue
		}
		if n%i == 0 {
			prime = false
		}
	}

	return prime
}
