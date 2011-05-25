package main

import (
	"log"
	"fmt"
	"os"
	"regexp"
	"exec"
	"time"
	"bufio"
)

type Stat struct{
	time string
	num_errors int
	matched_lines []string

	lines chan string // channel for lines into the handler
	eof chan bool // signal to handler to clean up and quit
	quit chan bool // signal that handler has quit
}

func (s *Stat) String() string {
	return fmt.Sprintf("Date: %s, Num errors: %d", s.time, s.num_errors)
}

type OutputHandler func(stat *Stat)

// Refine the below for build errors a bit more...
var go_build_error = regexp.MustCompile(".*:[0-9]*: .*")

func main() {
	// Run the rest of the args
	// Read and parse the output
	// (Save every build error matched line and test matched line into a log file for debugging)
	// Spit stats into a stats file per project
	// Print output
	// Print stats

	if len(os.Args) < 2 {
		log.Fatal("Usage: arya some_make_command otherargs")
	}
	cmd := os.Args[1]
	cmd_args := os.Args[2:]

	cmd_name, err := exec.LookPath(cmd)
	if err != nil {
		log.Fatal("cmd: %s could not be found in your PATH", cmd)
	}
	curr_env := os.Environ()
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal("Could not get working directory: %s", err.String())
	}
	command, err := exec.Run(cmd_name, cmd_args, curr_env, cwd, exec.PassThrough, exec.Pipe, exec.MergeWithStdout)
	if err != nil {
		log.Fatal("Error running cmd: %s", err.String())
	}
	_, err = command.Wait(0)
	if err != nil {
		log.Fatal("Cmd exited with error: %s", err.String())
	}

	feedHandlers(command.Stdout)
}

func feedHandlers(input *os.File) {
	buffered_reader := bufio.NewReader(input)
	handlers := []OutputHandler{GocheckHandler}
	num_handlers := len(handlers)
	quit := make(chan bool, num_handlers)
	stat_objs := make([]Stat, num_handlers)
	this_time := time.LocalTime().Format(time.UnixDate)
	for i,handler := range(handlers) {
		stat_objs[i].quit = quit
		stat_objs[i].time = this_time
		stat_objs[i].lines = make(chan string)
		stat_objs[i].eof = make(chan bool)
		// Start all the handlers up and pass them a stat object
		go handler(&stat_objs[i])
	}

	for {
		line, _, err := buffered_reader.ReadLine()
		if err == os.EOF {
			for _,stat := range(stat_objs) {
				stat.eof <- true
			}
			for y := 0; y < num_handlers; y++ {
				// wait for all the handlers to quit
				<-quit
			}
			fmt.Println(stat_objs)
			break
		}
		if err != nil {
			log.Fatal("Some error reading the input: ", err.String())
		}

		line_string := string(line)

		for i,_ := range(stat_objs) {
			stat_objs[i].lines <- line_string
		}
		fmt.Printf("> %s\n", line_string)
	}
}
