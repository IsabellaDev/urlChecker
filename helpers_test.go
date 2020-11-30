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

func Test404GetStatusCode(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
		w.Write([]byte("Not Found"))
	}))
	defer ts.Close()
	link := ts.URL
	result, err := getStatusCode(link)
	if err != nil {
		fmt.Println(err)
	}

	expected := 404
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected: %v, but got: %v", expected, result)
	}
}

func Test_removeDuplicate(t *testing.T) {
	type args struct {
		urls []string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		// TODO: Add test cases.
		{"name", args{[]string{"http://www.google.ca", "http://www.google.ca"}}, []string{"http://www.google.ca"}},
		{"name", args{[]string{"http://www.google.ca", "http://zyang.ca", "http://www.google.ca"}}, []string{"http://www.google.ca", "http://zyang.ca"}},
		{"name", args{[]string{}}, []string{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := removeDuplicate(tt.args.urls); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("removeDuplicate() = %v, want %v", got, tt.want)
			}
		})
	}
}
