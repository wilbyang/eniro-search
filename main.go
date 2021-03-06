// The server program issues Eniro search requests
//
// The /search endpoint accepts these query params:
//   q=the Eniro search query
//   timeout=a timeout for the request, in time.Duration format
//
// For example, http://localhost:8080/search?q=pizza&timeout=1s serves the
// first few Eniro search results for "pizza" or a "deadline exceeded" error
// if the timeout expires.
package main

import (
	"context"
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"time"
)

var (
	listenAddr = flag.String("listenAddr", ":8080", "Server address")
)

func main() {
	flag.Parse()
	http.HandleFunc("/search", handleSearch)
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)
	log.Fatal(http.ListenAndServe(*listenAddr, nil))
}

// handleSearch handles URLs like /search?q=pizza&timeout=1s by forwarding the
// query to google.Search. If the query param includes timeout, the search is
// canceled after that duration elapses.
func handleSearch(w http.ResponseWriter, req *http.Request) {
	// ctx is the Context for this handler. Calling cancel closes the
	// ctx.Done channel, which is the cancellation signal for requests
	// started by this handler.
	var (
		ctx    context.Context
		cancel context.CancelFunc
	)
	timeout, err := time.ParseDuration(req.FormValue("timeout"))
	if err == nil {
		// The request has a timeout, so create a context that is
		// canceled automatically when the timeout expires.
		ctx, cancel = context.WithTimeout(context.Background(), timeout)
	} else {
		ctx, cancel = context.WithCancel(context.Background())
	}
	defer cancel() // Cancel ctx as soon as handleSearch returns.

	// Check the search query.
	query := req.FormValue("q")
	if query == "" {
		http.Error(w, "no query", http.StatusBadRequest)
		return
	}

	// Run the Eniro search and print the results.
	start := time.Now()
	result, err := Search(ctx, query)
	elapsed := time.Since(start)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var ret = struct {
		Result           Result
		Timeout, Elapsed time.Duration
	}{
		Result:  result,
		Timeout: timeout,
		Elapsed: elapsed,
	}
	if err := json.NewEncoder(w).Encode(ret); err != nil {
		log.Print(err)
		return
	}
}
