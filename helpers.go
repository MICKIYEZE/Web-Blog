package main

import (
    "strings"
    "testing"
)

func AssertEqual(t *testing.T, expected, actual interface{}) {
    if expected != actual {
        t.Errorf("Expected %v, got %v", expected, actual)
    }
}

func AssertNoError(t *testing.T, err error) {
    if err != nil {
        t.Fatalf("Unexpected error: %v", err)
    }
}

func AssertError(t *testing.T, err error, expectedMsg string) {
    if err == nil {
        t.Fatal("Expected an error but got nil")
    }

    if !strings.Contains(err.Error(), expectedMsg) {
        t.Errorf("Expected error message to contain %q, got %q", expectedMsg, err.Error())
    }
}