package main

import (
	"flag"
	"log"
	"os"

	"github.com/tombuildsstuff/teamcity-go-test-json/logger"
	"github.com/tombuildsstuff/teamcity-go-test-json/runner"
)

func main() {
	count := flag.Int("count", 1, "The number of times that each test should be run")
	prefix := flag.String("prefix", "", "The Test Prefix, for example `TestAcc`")
	scope := flag.String("scope", "", "The directory scope where tests should be run. This'll be suffixed with `/...` if not specified.")
	parallelism := flag.Int("parallelism", 0, "The number of tests which should be run in parallel where possible")
	timeout := flag.Int("timeout", 0, "The maximum test duration in hours")

	flag.Parse()

	input := runner.ExecuteInput{
		Logger: logger.NewTeamCityTestLogger(),
		Debug:  os.Getenv("DEBUG") != "",
	}

	if scope != nil && *scope != "" {
		input.Scope = *scope
	}

	if prefix != nil && *prefix != "" {
		input.Prefix = *prefix
	}

	if parallelism != nil && *parallelism != 0 {
		input.Parallelism = *parallelism
	}

	if count != nil && *count != 0 {
		input.Count = *count
	}

	if timeout != nil && *timeout != 0 {
		input.TimeoutInHours = *timeout
	}

	if err := input.Execute(); err != nil {
		log.Printf("[ERROR] Error running executor: %+v", err)
		os.Exit(1)
	}
}
