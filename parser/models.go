package parser

import "fmt"

type goTestJsonLogLine struct {
	Time string `json:"Time"`

	// Action is the type of log message this is: e.g. output, skip, pass, fail, cont, run, pause
	Action string `json:"Action"`

	// Package is the path to the Go Package
	Package string `json:"Package"`

	// Output is the stdout
	Output string `json:"Output"`

	// Test is the name of the test being run
	Test string `json:"Test,omitempty"`

	// Elapsed is the time taken for this test in seconds (e.g. 12.34s)
	Elapsed float64 `json:"Elapsed,omitempty"`
}

func (tl goTestJsonLogLine) CacheKey() string {
	return fmt.Sprintf("%s-%s", tl.Package, tl.Test)
}

func (tl goTestJsonLogLine) Duration() int64 {
	// e.g. 76.12
	return int64(tl.Elapsed * 1000)
}
