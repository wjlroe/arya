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

type OutputHandler func(stat Stat)

// Refine the below for build errors a bit more...
var go_build_error = regexp.MustCompile(".*:[0-9]*: .*")


// TODO: Extract each output matcher into some kind of plugin

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
	fmt.Printf("cmd: %s\n", cmd)
	fmt.Printf("cmd_args: %s\n", cmd_args)

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

	// TODO: Refactor to bufio here and send each line to the handlers, then collect results
	//stat := GocheckHandler(command.Stdout)
	//stat := GoTestHandler(command.Stdout)

	feedHandlers(command.Stdout)
}

func feedHandlers(input *os.File) {
	buffered_reader := bufio.NewReader(input)
	handlers := []OutputHandler{GocheckHandler}
	quit := make(chan bool, len(handlers))
	stat_objs := make([]Stat, len(handlers))
	this_time := time.LocalTime().Format(time.UnixDate)
	for i,handler := range(handlers) {
		stat_objs[i].quit = quit
		stat_objs[i].time = this_time
		// Start all the handlers up and pass them a stat object
		go handler(stat_objs[i])
	}

	for {
		line, _, err := buffered_reader.ReadLine()
		if err == os.EOF {
			fmt.Println("Going to send the eof signal to all the handlers")
			for _,stat := range(stat_objs) {
				stat.eof <- true
			}
			fmt.Println("Going to wait for all the handlers to quit now")
			for {
				// wait for all the handlers to quit
				<-quit
			}
			fmt.Println("All handlers have quit")
			fmt.Println(stat_objs)
			break
		}
		if err != nil {
			log.Fatal("Some error reading the input: ", err.String())
		}

		line_string := string(line)

		for _,stat := range(stat_objs) {
			stat.lines <- line_string
		}
	}
}
