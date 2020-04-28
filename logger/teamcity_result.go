package logger

import (
	"fmt"

	uuid "github.com/satori/go.uuid"
	"github.com/tombuildsstuff/teamcity-go-test-json/models"
)

type TeamCityTestResultLogger struct {
	logger TeamCityTestResultFormatter
}

func (rl TeamCityTestResultLogger) Log(result models.TestResult) string {
	flowId := uuid.NewV4().String()

	// output the package name for logical grouping
	builder := rl.logger.TestSuiteStarted(result.Package, flowId)

	// always has to be a start
	builder += rl.logger.TestStart(result.TestName, flowId)

	switch result.Result {
	case models.Successful:
		builder += rl.logger.TestStdOut(result.TestName, result.StdOut, flowId)
		break

	case models.Failed:
		builder += rl.logger.TestStdOut(result.TestName, result.StdOut, flowId)
		builder += rl.logger.TestStdErr(result.TestName, result.StdErr, flowId)
		builder += rl.logger.TestFailed(result.TestName, "Test Failed", "", flowId)
		break

	case models.Ignored:
		builder += rl.logger.TestStdOut(result.TestName, result.StdOut, flowId)
		builder += rl.logger.TestIgnored(result.TestName, "Test ignored", flowId)
		break

	default:
		panic(fmt.Errorf("Unexpected Test State %q - this is a bug in the test runner", result.Result))
	}

	// e.g. 76.12s -> 76120ms
	testDuration := int64(result.Duration * 1000)

	// always has to be a finish
	builder += rl.logger.TestFinished(result.TestName, testDuration, flowId)

	// output the package name for logical grouping
	builder += rl.logger.TestSuiteFinished(result.Package, flowId)

	return builder
}
