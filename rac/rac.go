package rac

import (
	"fmt"
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

	data := strings.SplitN(line, ":", 2)
	key := strings.TrimSpace(data[0])
	value := strings.TrimSpace(data[1])

	key = strings.ReplaceAll(key, `"`, "")
	value = strings.ReplaceAll(value, `"`, "")

	if key == "" {
		return key, value, &ErrNoKeyInRACLine{}
	}

	return key, value, nil
}

func Parse(output string) ([]map[string]string, error) {
	var records []map[string]string
	blocks := strings.Split(output, "\n\n")
	fmt.Println(blocks)

	for _, block := range blocks {
		record := map[string]string{}

		for idx, line := range strings.Split(block, "\n") {
			key, value, err := extractKeyValue(line)
			if err != nil {
				return nil, fmt.Errorf("error parsing rac output (line %d): %w", idx, err)
			}
			record[key] = value
		}

		records = append(records, record)
	}

	return records, nil
}
