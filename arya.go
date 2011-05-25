package main

import (
	"log"
	"fmt"
	"os"
	"regexp"
	"exec"
	"time"
)

type Stat struct{
	time string
	num_errors int
	matched_lines []string
}

func (s *Stat) String() string {
	return fmt.Sprintf("Date: %s, Num errors: %d", s.time, s.num_errors)
}

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
	stat := GocheckHandler(command.Stdout)
	//stat := GoTestHandler(command.Stdout)
	stat.time = time.LocalTime().Format(time.UnixDate)
	fmt.Printf("Stats: %s\n", stat)
}

func feedHandlers(stdout os.File, lines chan string, eof chan bool) {
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

		lines <- line_string // send the line to the multiplexer

	}
}

type OutputHandler func(lines chan string, eof chan bool)

func demultiplexer(handlers OutputHandler) (lines chan string) {
	lines := make(chan string)
	go func() {
		for {
			line <- lines

		}
	}()
	return lines
}