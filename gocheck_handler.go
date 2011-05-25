package main

import (
	"fmt"
	"regexp"
	"strconv"
)

var go_check_ok = regexp.MustCompile("OK: [0-9]* passed")
var go_check_fail = regexp.MustCompile("OOPS: [0-9]* passed, ([0-9]*) FAILED")

func GocheckHandler(stat Stat) {
	for {
		select {
		case line := <-stat.lines:
processLine(&stat, line)
case <-stat.eof:
stat.quit <- true
}
}
}

// func GocheckHandlerAgain(input io.Reader) {
// 	stat := &Stat{num_errors: 0, matched_lines: []string{}}
// 	matched_lines := []string{}
// 	buffered_reader := bufio.NewReader(input)

// 	for {
// 		line, _, err := buffered_reader.ReadLine()
// 		if err == os.EOF {
// 			break
// 		}
// 		if err != nil {
// 			log.Fatal("Some error reading the input: ", err.String())
// 		}

// 		line_string := string(line)
// 		fmt.Printf("line_string: %s\n", line_string)
// 	}
// }

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

	fmt.Printf("Matched_lines: %s\n", stat.matched_lines)
}