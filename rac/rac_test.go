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

func TestExtractKeyValueColonsInValue(t *testing.T) {
	input := `connected-at   : 2024-04-06T22:00:03`
	expectedKey := "connected-at"
	expectedValue := "2024-04-06T22:00:03"

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

func TestRACQueryParseSingleBlock(t *testing.T) {
	input := `connection     : 3f97c035-b8e6-4b25-a72c-887b51a72b67
	conn-id        : 1168
	application    : "WebServerExtension"
	connected-at   : 2024-04-06T22:00:03
	blocked-by-ls  : 0`

	expected := []map[string]string{
		{
			"connection":    "3f97c035-b8e6-4b25-a72c-887b51a72b67",
			"conn-id":       "1168",
			"application":   "WebServerExtension",
			"connected-at":  "2024-04-06T22:00:03",
			"blocked-by-ls": "0",
		},
	}

	query := RACQuery{
		Output: input,
	}
	err := query.Parse()
	if err != nil {
		t.Error(err)
	}

	actual := query.Records
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("(actual) %#v != %#v (expected)\n", actual, expected)
	}
}

func TestRACQueryParseMultipleBlocksOK(t *testing.T) {
	input := `connection     : 3f97c035-b8e6-4b25-a72c-887b51a72b67
	conn-id        : 1168
	application    : "WebServerExtension"
	connected-at   : 2024-04-06T22:00:03
	blocked-by-ls  : 0

	connection     : 5f3777f8-75a2-4b6e-a9da-7d24cac4bb21
	conn-id        : 678
	application    : "WebServerDoodad"
	connected-at   : 2024-05-04T21:26:01
	blocked-by-ls  : 1`

	expected := []map[string]string{
		{
			"connection":    "3f97c035-b8e6-4b25-a72c-887b51a72b67",
			"conn-id":       "1168",
			"application":   "WebServerExtension",
			"connected-at":  "2024-04-06T22:00:03",
			"blocked-by-ls": "0",
		},
		{
			"connection":    "5f3777f8-75a2-4b6e-a9da-7d24cac4bb21",
			"conn-id":       "678",
			"application":   "WebServerDoodad",
			"connected-at":  "2024-05-04T21:26:01",
			"blocked-by-ls": "1",
		},
	}

	query := RACQuery{
		Output: input,
	}
	err := query.Parse()
	if err != nil {
		t.Error(err)
	}

	actual := query.Records
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("(actual) %#v != %#v (expected)\n", actual, expected)
	}
}

func TestRACQueryParseMultipleBlocksCount(t *testing.T) {
	input := `connection     : 3f97c035-b8e6-4b25-a72c-887b51a72b67
	conn-id        : 1168
	application    : "WebServerExtension"
	connected-at   : 2024-04-06T22:00:03
	blocked-by-ls  : 0

	connection     : 5f3777f8-75a2-4b6e-a9da-7d24cac4bb21
	conn-id        : 678
	application    : "WebServerDoodad"
	connected-at   : 2024-05-04T21:26:01
	blocked-by-ls  : 1`

	expected := 2

	query := RACQuery{
		Output: input,
	}
	err := query.Parse()
	if err != nil {
		t.Error(err)
	}

	actual := query.CountRecords()
	if actual != expected {
		t.Fatalf("(actual) %#v != %#v (expected)\n", actual, expected)
	}
}
