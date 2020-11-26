package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestExtractURL(t *testing.T) {
	result := extractURL("https://www.google.ca, https://www.facebook.com")

	expected := []string{"https://www.google.ca", "https://www.facebook.com"}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected: %v, but got: %v", expected, result)
	}
}

func TestGetStatusCode(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte("Ok"))
	}))
	defer ts.Close()
	link := ts.URL
	result, err := getStatusCode(link)
	if err != nil {
		fmt.Println(err)
	}

	expected := 200
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected: %v, but got: %v", expected, result)
	}
}
func Test400GetStatusCode(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(400)
		w.Write([]byte("Not Found"))
	}))
	defer ts.Close()
	link := ts.URL
	result, err := getStatusCode(link)
	if err != nil {
		fmt.Println(err)
	}

	expected := 400
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected: %v, but got: %v", expected, result)
	}
}
