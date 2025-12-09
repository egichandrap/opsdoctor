package utils

import (
	"encoding/json"
	"fmt"
	"os"
)

var noColor bool
var verbose bool
var outFormat string

// Mapping warna untuk PrintColor()
var colors = map[string]string{
	"reset":  "\033[0m",
	"red":    "\033[31m",
	"green":  "\033[32m",
	"yellow": "\033[33m",
	"blue":   "\033[34m",
	"cyan":   "\033[36m",
	"white":  "\033[37m",
	"bold":   "\033[1m",
}

func SetNoColor(v bool) {
	noColor = v
}

func SetVerbose(v bool) {
	verbose = v
}

func SetOutFormat(f string) {
	outFormat = f
}

func IsJSONOutput() bool {
	return outFormat == "json"
}

func VPrintf(format string, a ...interface{}) {
	if verbose {
		fmt.Printf(format, a...)
	}
}

func PrintColor(color string, msg string) {
	if noColor || IsJSONOutput() {
		fmt.Println(msg)
		return
	}
	code, ok := colors[color]
	if !ok {
		fmt.Println(msg)
		return
	}
	fmt.Println(code + msg + colors["reset"])
}

func ExportJSON(filename string, data interface{}) error {
	b, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filename, b, 0644)
}
