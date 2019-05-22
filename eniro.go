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

const baseURL = "https://api.eniro.com/partnerapi/cs/search/basic"

type Result struct {
	Title        string    `json:"title"`
	Query        string    `json:"query"`
	TotalHits    uint32    `json:"totalHits"`
	TotalCount   uint32    `json:"totalCount"`
	StartIndex   uint32    `json:"startIndex"`
	ItemsPerPage uint16    `json:"itemsPerPage"`
	Adverts      []Advert2 `json:"adverts"`
}

// use map to make it possible to do json serialization projection
// check here
// https://stackoverflow.com/questions/17306358/removing-fields-from-struct-or-hiding-them-in-json-response
type Advert2 map[string]interface{}

type Advert struct {
	EniroId        string       `json:"eniroId"`
	CompanyInfo    CompanyInfo  `json:"companyInfo"`
	Address        Address      `json:"address"`
	Location       Location     `json:"location"`
	PhoneNumbers   PhoneNumbers `json:"phoneNumbers"`
	CompanyReviews string       `json:"companyReviews"`
	Homepage       string       `json:"homepage"`
	Facebook       string       `json:"facebook"`
	InfoPageLink   string       `json:"infoPageLink"`
}

type CompanyInfo struct {
	CompanyName string `json:"companyName"`
	OrgNumber   string `json:"orgNumber"`
	CompanyText string `json:"companyText"`
}

type Address struct {
	StreetName string `json:"streetName"`
	PostCode   string `json:"postCode"`
	PostArea   string `json:"postArea"`
	PostBox    string `json:"postBox"`
}

type Coordinate struct {
	Use       string  `json:"use"`
	Longitude float32 `json:"longitude"`
	Latitude  float64 `json:"latitude"`
}
type Coordinates []Coordinate

type Location struct {
	Coordinates Coordinates `json:"coordinates"`
}
type PhoneNumber struct {
	Kind        string `json:"type"`
	PhoneNumber string `json:"phoneNumber"`
	Label       string `json:"label"`
}
type PhoneNumbers []PhoneNumber

// Search sends query to Eniro search and returns the result.
func Search(ctx context.Context, query string) (Result, error) {
	// Prepare the Eniro Search API request.
	//?profile=interview&key=114773894699415832&country=se&version=1.1.3&search_word=pizza

	req, err := http.NewRequest("GET", baseURL, nil)
	if err != nil {
		return Result{}, err
	}
	q := req.URL.Query()
	q.Set("profile", "interview")
	q.Set("key", "114773894699415832")
	q.Set("country", "se")
	q.Set("version", "1.1.3")
	q.Set("search_word", query)

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
