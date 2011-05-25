package main

import (
	"fmt"
	"regexp"
	"strconv"
)

var go_check_ok = regexp.MustCompile("OK: [0-9]* passed")
var go_check_fail = regexp.MustCompile("OOPS: [0-9]* passed, ([0-9]*) FAILED")

func GocheckHandler(stat *Stat) {
L:
	for {
		select {
		case line := <-stat.lines:
			processLine(stat, line)
		case <-stat.eof:
			stat.quit <- true
			break L
		}
	}
}

func processLine(stat *Stat, line string) {
	fail_matches := go_check_fail.FindStringSubmatch(line)
	if fail_matches != nil {
		fail_num, err := strconv.Atoi(fail_matches[1])
		if err != nil {
			fmt.Printf("Could not convert %s into integer\n", fail_matches[1])
		} else {
			stat.num_errors += fail_num // assumes one match per line
		}
		stat.matched_lines = append(stat.matched_lines, line)
	}

	//fmt.Printf("Matched_lines: %s\n", stat.matched_lines)
}