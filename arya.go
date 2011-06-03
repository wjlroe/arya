package main

import (
	"log"
	"fmt"
	"os"
	"io"
	"regexp"
	"exec"
	"time"
	"bytes"
	"bufio"
	"strings"
	"io/ioutil"
	"path/filepath"
)

type Stat struct{
	time string
	num_errors int
	matched_lines []string
	matched bool
	project_name string

	lines chan string // channel for lines into the handler
	eof chan bool // signal to handler to clean up and quit
	quit chan bool // signal that handler has quit
}

func (s Stat) String() string {
	return fmt.Sprintf("Date: %s, Num errors: %d", s.time, s.num_errors)
}

func (stat *Stat) save() {
	home_dir := os.Getenv("HOME")
	stat_csv := filepath.Join(home_dir, ".arya", stat.project_name + ".csv")
	fmt.Printf("CSV: %s\n", stat_csv)

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
	cmd_args := strings.Join(os.Args[2:], " ")

	cmd_name, err := exec.LookPath(cmd)
	if err != nil {
		log.Fatal("cmd: %s could not be found in your PATH", cmd)
	}

	var buffer []byte
	output_buffer := bytes.NewBuffer(buffer)
	cmd_obj := exec.Command(cmd_name, cmd_args)
	cmd_obj.Stdout = output_buffer
	cmd_obj.Stderr = output_buffer
	//err = cmd_obj.Start()
	// TODO: refactor to run and stream output to handlers...
	// Problem is it sends EOF right away
	err = cmd_obj.Run() // blocks until finished
	// if err != nil {
	// 	log.Fatal("Error running cmd: ", err.String())
	// }

	stats := feedHandlers(output_buffer)
	processStats(stats)
	//cmd_obj.Wait()
}

func feedHandlers(input io.Reader) []Stat {
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
			fmt.Printf("%v\n", stat_objs)
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
	return stat_objs
}

func processStats(stats []Stat) {
	var matched_stat *Stat
	for i,_ := range(stats) {
		if stats[i].matched {
			if matched_stat != nil {
				log.Fatal("More than one output handler matched - this is bad.")
			}
			matched_stat = &stats[i]
		}
	}
	if matched_stat == nil {
		log.Fatal("No output handler matched! No stats recorded.")
	}
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal("Couldn't get the current working directory! ", err.String())
	}
	project_root := filepath.Join(cwd, ".arya")
	fi, err := os.Stat(project_root)
	if err != nil {
		log.Fatal("Could not find .arya file. ", err.String())
	}
	if fi.IsRegular() {
		arya_file, err := os.Open(project_root)
		if err != nil {
			log.Fatal("Couldn't open .arya file! ", err.String())
		}
		defer arya_file.Close()
		var content []byte
		content, err = ioutil.ReadAll(arya_file)
		if err != nil {
			log.Fatal("Error reading .arya file ", err.String())
		}
		content_string := strings.TrimSpace(string(content))
		if content_string == "" {
			matched_stat.project_name = filepath.Base(cwd)
		} else {
			matched_stat.project_name = content_string
		}
		fmt.Printf("Project is: %s\n", matched_stat.project_name)
	} else {
		log.Fatal(".arya not a real file")
	}
	matched_stat.save()
}
