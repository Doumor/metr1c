package rac

import (
	"errors"
	"reflect"
	"testing"
)

func TestExtractKeyValueNormal(t *testing.T) {
	input := "connection     : 3f97c035-b8e6-4b25-a72c-887b51a72b67"
	expectedKey := "connection"
	expectedValue := "3f97c035-b8e6-4b25-a72c-887b51a72b67"

	actualKey, actualValue, err := extractKeyValue(input)
	if err != nil {
		t.Error(err)
	}

	if expectedKey != actualKey || expectedValue != actualValue {
		t.Fatalf("(actual) %#v: %#v != %#v: %#v (expected)\n", actualKey, actualValue, expectedKey, expectedValue)
	}
}

func TestExtractKeyValueQuotes(t *testing.T) {
	input := `application     : "WebServerExtension"`
	expectedKey := "application"
	expectedValue := "WebServerExtension"

	actualKey, actualValue, err := extractKeyValue(input)
	if err != nil {
		t.Error(err)
	}

	if expectedKey != actualKey || expectedValue != actualValue {
		t.Fatalf("(actual) %#v: %#v != %#v: %#v (expected)\n", actualKey, actualValue, expectedKey, expectedValue)
	}
}

func TestExtractKeyValueNoColon(t *testing.T) {
	input := `application      "WebServerExtension"`

	_, _, err := extractKeyValue(input)
	if !errors.Is(err, &ErrNoColonInRACLine{}) {
		t.Errorf("expected error not raised: %s", &ErrNoColonInRACLine{})
	}
}

func TestExtractKeyValueNoKey(t *testing.T) {
	input := `: "WebServerExtension"`

	_, _, err := extractKeyValue(input)
	if !errors.Is(err, &ErrNoKeyInRACLine{}) {
		t.Errorf("expected error not raised: %s", &ErrNoKeyInRACLine{})
	}
}

