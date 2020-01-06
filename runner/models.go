package runner

type ExecuteInput struct {
	// The scope at which tests should be run e.g. `./azurerm/internal/services/resource/tests/...`
	Scope string

	// The prefix which should be used for matching tests e.g. `TestAcc`
	Prefix string

	// The number of tests which should be run in parallel
	Parallelism int

	// The number of times that each test should be run
	Count int

	// Should parsed output be logged to stdout for debugging purposes?
	Debug bool

	// The timeout in hours which should be used for the tests
	TimeoutInHours int
}
