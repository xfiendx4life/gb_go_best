// This program is was created to delete duplicated files from directory.

package main

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/xfiendx4life/gb_gobest/hw2/files_deleter"
	loggermaker "github.com/xfiendx4life/gb_gobest/hw2/logger_maker"
	"go.uber.org/zap"
)

var (
	delete   bool
	dir      string
	ErrHelp  = errors.New("flag: help requested")
	logLevel = zap.LevelFlag("loglevel", zap.InfoLevel, "set logging level")
	filelog  string
)

func init() {
	flag.BoolVar(&delete, "delete", false, "set true if you want to delete duplicate files")
	flag.StringVar(&dir, "dir", ".", "choose direcory to work with")
	flag.StringVar(&filelog, "filelog", "", "choose file for logs, leave empty to use stderr")

}

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stdout, "Usage of %s:\n", os.Args[0])
		fmt.Fprintln(os.Stdout, "This program was created to delete duplicated files from directory.")
		flag.PrintDefaults()
	}
	flag.Parse()
	logger := loggermaker.InitLogger(logLevel, filelog)
	n, _ := files_deleter.Delete(dir, delete, logger)
	fmt.Println(n)

}
