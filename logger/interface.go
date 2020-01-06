package logger

import "github.com/tombuildsstuff/teamcity-go-test-json/models"

type TestResultLogger interface {
	Log(result models.TestResult) string
}
