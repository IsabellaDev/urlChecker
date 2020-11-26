package main

import (
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

func TestParseFromTelescope(t *testing.T) {

}
