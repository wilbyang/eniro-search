// Package google provides a function to do Eniro searches using the Eniro Web
// Search API. See https://developers.google.com/web-search/docs/
//
// This package is an example to accompany https://blog.golang.org/context.
// It is not intended for use by others.
//
// Eniro has since disabled its search API,
// and so this package is no longer useful.
package main

import (
	"context"
	"encoding/json"
	"net/http"
)

const baseURL = "https://api.eniro.com/api/v2.0/SE/search"

type Result struct {
	Companies []struct {
		EniroID string `json:"eniroId"`
		Name    string `json:"name"`
		Phones  []struct {
			Number string `json:"number"`
			Text   string `json:"text"`
		} `json:"phones"`
		ProfilePageLink string `json:"profilePageLink"`
		Addresses       []struct {
			PostalCode   string `json:"postalCode"`
			StreetName   string `json:"streetName"`
			StreetNumber string `json:"streetNumber"`
			PostalArea   string `json:"postalArea"`
		} `json:"addresses"`
	} `json:"companies"`
	SearchLevel string `json:"searchLevel"`
	Hits        struct {
		Companies int `json:"companies"`
	} `json:"hits"`
}

// Search sends query to Eniro search and returns the result.
func Search(ctx context.Context, query string) (Result, error) {
	// Prepare the Eniro Search API request.
	//?profile=interview&key=114773894699415832&country=se&version=1.1.3&search_word=pizza

	req, err := http.NewRequest("GET", baseURL, nil)
	if err != nil {
		return Result{}, err
	}
	q := req.URL.Query()
	req.Header.Set("Authorization", "Basic eWFuZy53aWxieUBnbWFpbC5jb206RUlzcjR3VUtZZ0FiQ203YVF3UWN4N2VfNHRSNFBxOFBPM2JlcnBlS0V0QQ==")
	q.Set("query", query)

	req.URL.RawQuery = q.Encode()

	// Issue the HTTP request and handle the response. The httpDo function
	// cancels the request if ctx.Done is closed.
	var result Result
	err = httpDo(ctx, req, func(resp *http.Response, err error) error {
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		// unserialize the response to go struct
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return err
		}

		return nil
	})
	// httpDo waits for the closure we provided to return, so it's safe to
	// read result here.
	return result, err
}

// httpDo issues the HTTP request and calls f with the response. If ctx.Done is
// closed while the request or f is running, httpDo cancels the request, waits
// for f to exit, and returns ctx.Err. Otherwise, httpDo returns f's error.
func httpDo(ctx context.Context, req *http.Request, f func(*http.Response, error) error) error {
	// Run the HTTP request in a goroutine and pass the response to f.
	c := make(chan error, 1)
	req = req.WithContext(ctx)
	go func() { c <- f(http.DefaultClient.Do(req)) }()
	select {
	case <-ctx.Done():
		<-c // Wait for f to return.
		return ctx.Err()
	case err := <-c:
		return err
	}
}
