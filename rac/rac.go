package rac

import (
	"fmt"
	"strings"
)

type ErrMalformedRACLine struct{}

func (e *ErrMalformedRACLine) Error() string {
	return "could not split the RAC output line"
}

func extractKeyValue(line string) (string, string, error) {
	data := strings.Split(line, ":")
	key := strings.TrimSpace(data[0])
	value := strings.TrimSpace(data[1])

	key = strings.ReplaceAll(key, `"`, "")
	value = strings.ReplaceAll(value, `"`, "")

	if key == "" || value == "" {
		return key, value, &ErrMalformedRACLine{}
	}

	return key, value, nil
}

