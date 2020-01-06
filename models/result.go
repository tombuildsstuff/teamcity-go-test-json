package models

type Result string

var (
	Failed     Result = "Failed"
	Ignored    Result = "Ignored"
	Pending    Result = "Pending"
	Successful Result = "Successful"
)

type TestResult struct {
	Package  string
	TestName string
	Result   Result
	Duration float64
	StdOut   string
	StdErr   string
}
