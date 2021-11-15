// This program is was created to delete duplicated files from directory.

package main

import (
	"errors"
	"filesDeleter/files_deleter"
	"flag"
	"fmt"
	"os"
)

var (
	delete  bool
	dir     string
	ErrHelp = errors.New("flag: help requested")
	
)



func init() {
	flag.BoolVar(&delete, "delete", false, "set true if you want to delete duplicate files")
	flag.StringVar(&dir, "dir", ".", "choose direcory to work with")
}

func main() {
	flag.Usage   = func() {
		fmt.Fprintf(os.Stdout, "Usage of %s:\n", os.Args[0])
		fmt.Fprintln(os.Stdout, "This program was created to delete duplicated files from directory.")
		flag.PrintDefaults()
	}
	flag.Parse()
	n, _ := files_deleter.Delete(dir, delete)
	fmt.Println(n)
}
