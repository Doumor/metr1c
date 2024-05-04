package rac

import (
	// "fmt"
	"strings"
)

type ErrNoColonInRACLine struct{}

func (e *ErrNoColonInRACLine) Error() string {
	return "error: no colon in a RAC output line"
}

type ErrNoKeyInRACLine struct{}

func (e *ErrNoKeyInRACLine) Error() string {
	return "error: no key found in RAC output line"
}

func extractKeyValue(line string) (string, string, error) {
	if !strings.Contains(line, ":") {
		return "", "", &ErrNoColonInRACLine{}
	}

	data := strings.Split(line, ":")
	key := strings.TrimSpace(data[0])
	value := strings.TrimSpace(data[1])

	key = strings.ReplaceAll(key, `"`, "")
	value = strings.ReplaceAll(value, `"`, "")

	if key == "" {
		return key, value, &ErrNoKeyInRACLine{}
	}

	return key, value, nil
}

