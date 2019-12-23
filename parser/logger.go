package parser

import (
	"log"
	"strings"

	"github.com/tombuildsstuff/teamcity-go-test-json/logger"
)

func (line GoTestJsonLogLine) Log(logger logger.TeamCityTestLogger) error {
	// we don't care about packages for now
	if line.Test == "" {
		return line.logPackage(logger)
	}

	return line.logTest(logger)
}

func (tl GoTestJsonLogLine) logPackage(logger logger.TeamCityTestLogger) error {
	// TODO: log the start/end of packages
	return nil
}

func (tl GoTestJsonLogLine) logTest(logger logger.TeamCityTestLogger) error {
	if strings.EqualFold(tl.Action, "run") {
		logger.TestStart(tl.Test)
		return nil
	}

	if strings.EqualFold(tl.Action, "output") {
		logger.TestStdOut(tl.Test, tl.Output)
		return nil
	}

	if strings.EqualFold(tl.Action, "pass") {
		logger.TestFinished(tl.Test, tl.Duration())
		return nil
	}

	if strings.EqualFold(tl.Action, "fail") {
		logger.TestFailed(tl.Test, tl.Duration(), "Test Failed", "")
		return nil
	}

	if strings.EqualFold(tl.Action, "skip") {
		logger.TestIgnored(tl.Test, "")
		return nil
	}

	if strings.EqualFold(tl.Action, "cont") || strings.EqualFold(tl.Action, "pause") {
		// these happen because we're running tests in parallel
		// but there's nothing notable for TC here - so we can ignore them
		return nil
	}

	log.Printf("[DEBUG] Unexpected Action: %q", tl.Action)
	return nil
}
