package main

import (
	"bufio"
	"io"
	"fmt"
	"log"
	"os"
	"regexp"
)

var go_test_fail = regexp.MustCompile("--- FAIL: .*") // need to be counted
var go_test_pass = regexp.MustCompile("^PASS$")

func GoTestHandler(input io.Reader) *Stat {
	stat := &Stat{num_errors: 0, matched_lines: []string{}}
	matched_lines := []string{}
	buffered_reader := bufio.NewReader(input)

	for {
		line, _, err := buffered_reader.ReadLine()
		if err == os.EOF {
			break
		}
		if err != nil {
			log.Fatal("Some error reading the input: ", err.String())
		}

		line_string := string(line)
		fmt.Printf("line_string: %s\n", line_string)
		fail_matches := go_test_fail.FindIndex(line)
		if fail_matches != nil {
			stat.num_errors += 1 // assumes one match per line
			matched_lines = append(matched_lines, line_string)
		}
	}
	fmt.Printf("Matched_lines: %s\n", matched_lines)
	stat.matched_lines = matched_lines
	return stat
}