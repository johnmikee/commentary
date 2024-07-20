package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/johnmikee/commentary/internal/commentary"
)

func main() {
	dir := flag.String("dir", ".", "directory to scan for .go files")
	write := flag.Bool("write", false, "write changes to files")
	flag.Parse()

	err := commentary.ProcessDirectory(*dir, *write)
	if err != nil {
		fmt.Println("Error processing directory:", err)
		os.Exit(1)
	}
}
