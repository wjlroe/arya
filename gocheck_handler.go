package main

import (
	"bufio"
	"io"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
)

var go_check_ok = regexp.MustCompile("OK: [0-9]* passed")
var go_check_fail = regexp.MustCompile("OOPS: [0-9]* passed, ([0-9]*) FAILED")

func GocheckHandler(input io.Reader) *Stat {
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
		fail_matches := go_check_fail.FindStringSubmatch(line_string)
		if fail_matches != nil {
			fail_num, err := strconv.Atoi(fail_matches[1])
			if err != nil {
				fmt.Printf("Could not convert %s into integer\n", fail_matches[1])
			} else {
				stat.num_errors += fail_num // assumes one match per line
			}
			matched_lines = append(matched_lines, line_string)
		}
	}
	fmt.Printf("Matched_lines: %s\n", matched_lines)
	stat.matched_lines = matched_lines
	return stat
}