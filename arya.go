package main

import (
	"log"
	"fmt"
	"os"
)

func main() {
	// Run the rest of the args
	// Read and parse the output
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


}