/*
rac implements types and methods for running queries with the `rac` tool,
as well as parsing its output.
*/
package rac

import (
	"bytes"
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
	Option     string
	Cluster    string
	User       string
	Password   string
	Output     string
	Records    []map[string]string
}

// Run a query against the `rac` tool
func (q *RACQuery) Run() error {
	var stdout, stderr bytes.Buffer
	cmd := exec.Command(q.ExecPath, q.Command, q.SubCommand, q.Option, q.Cluster, q.User, q.Password)

	// Cannot pass empty arg
	if q.Option == "" {
		cmd = exec.Command(q.ExecPath, q.Command, q.SubCommand, q.Cluster, q.User, q.Password)
	}

	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("error running a '%s %s %s' command (%w): %s", q.Command, q.SubCommand, q.Option, err, stderr.String())
	}
	q.Output = string(stdout.String())

	return nil
}

// Parse converts `rac` output lines into a slice of map[string]string records
func (q *RACQuery) Parse() error {
	outputCleaned := strings.TrimSpace(q.Output)
	blocks := strings.Split(outputCleaned, "\n\n")

	// If rac's output is (effectively) empty, just return.
	// Since Split always returns a slice of length >= 1, check the first item's length in characters
	if len(blocks[0]) == 0 {
		return nil
	}
	// Debug
	//fmt.Println(blocks)

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
