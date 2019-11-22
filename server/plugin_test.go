package main

import (
	"testing"
)

func TestUrbanDictionary(t *testing.T) {
	_, err := UrbanDictionary("/define bonfire")
	if err != nil {
		t.Fatal(err.Error())
	}
}
