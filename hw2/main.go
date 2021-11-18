// This program is was created to delete duplicated files from directory.

package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/xfiendx4life/gb_gobest/hw2/files_deleter"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	delete   bool
	dir      string
	ErrHelp  = errors.New("flag: help requested")
	logLevel = zap.LevelFlag("loglevel", zap.InfoLevel, "set logging level")
	filelog  string
	logger   *zap.SugaredLogger
)

func init() {
	flag.BoolVar(&delete, "delete", false, "set true if you want to delete duplicate files")
	flag.StringVar(&dir, "dir", ".", "choose direcory to work with")
	flag.StringVar(&filelog, "filelog", "", "choose file for logs, leave empty to use stderr")

}

func InitLogger(level *zapcore.Level, filelog string) {
	var output io.Writer
	var encoder zapcore.Encoder
	// choosing file or stderr
	if filelog != "" {
		output, _ = os.Create(filelog)                                     // we are going to use file as log output
		encoder = zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()) // using json for file
	} else {
		output = os.Stderr
		encoder = zapcore.NewConsoleEncoder(zap.NewProductionEncoderConfig()) // using simple console

	}

	writeSyncer := zapcore.AddSync(output)
	core := zapcore.NewCore(encoder, writeSyncer, level)
	logger = zap.New(core).Sugar()

}

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stdout, "Usage of %s:\n", os.Args[0])
		fmt.Fprintln(os.Stdout, "This program was created to delete duplicated files from directory.")
		flag.PrintDefaults()
	}
	flag.Parse()
	InitLogger(logLevel, filelog)
	logger.Debugf("Checking debug")
	logger.Errorf("Checking error")
	logger.Info("checking info")
	n, _ := files_deleter.Delete(dir, delete)
	fmt.Println(n)

}
