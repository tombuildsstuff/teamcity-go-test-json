package runner

import (
	"github.com/tombuildsstuff/teamcity-go-test-json/logger"
)

type ExecuteInput struct {
	// The scope at which tests should be run e.g. `./azurerm/internal/services/resource/tests/...`
	Scope string

	// The prefix which should be used for matching tests e.g. `TestAcc`
	Prefix string

	// The number of tests which should be run in parallel
	Parallelism int

	// The number of times that each test should be run
	Count int

	// Logger is the Test Logger where output should be sent
	Logger logger.TeamCityTestLogger

	// Should parsed output be logged to stdout for debugging purposes?
	Debug bool
}
