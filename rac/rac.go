/*
rac implements types and methods for running queries with the `rac` tool,
as well as parsing its output.
*/
package rac

import (
	"fmt"
	"os/exec"
	"strings"
)

// ErrNoColonInRACLine is returned when colon delimiter is found in a `rac` output line
type ErrNoColonInRACLine struct{}

func (e *ErrNoColonInRACLine) Error() string {
	return "error: no colon in a RAC output line"
}

// ErrNoKeyInRACLine is returned when a `rac` output line somehow has nothing or only whitespaces before the first colon
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

// RACQuery queries the `rac` tool and parses the output
type RACQuery struct {
	ExecPath   string
	Command    string
	SubCommand string
	Cluster    string
	User       string
	Password   string
	Output     string
	Records    []map[string]string
}

// Run a query against the `rac` tool
func (q *RACQuery) Run() error {
	output, err := exec.Command(q.ExecPath, q.Command, q.SubCommand, q.Cluster, q.User, q.Password).Output()
	if err != nil {
		return fmt.Errorf("error running a '%s %s' command: %w", q.Command, q.SubCommand, err)
	}
	q.Output = string(output)

	return nil
}

// Parse converts `rac` output lines into a slice of map[string]string records
func (q *RACQuery) Parse() error {
	blocks := strings.Split(q.Output, "\n\n")
	fmt.Println(blocks)

	for bindx, block := range blocks {
		record := map[string]string{}

		for lidx, line := range strings.Split(block, "\n") {
			key, value, err := extractKeyValue(line)
			if err != nil {
				return fmt.Errorf("error parsing rac output (block %d, line %d): %w", bindx, lidx, err)
			}
			record[key] = value
		}

		q.Records = append(q.Records, record)
	}

	return nil
}

// CountRecords returns the number of records in the `rac` query
func (q *RACQuery) CountRecords() int {
	return len(q.Records)
}
