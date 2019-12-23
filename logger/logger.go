package logger

import (
	"fmt"
	"strings"
)

type TeamCityTestLogger struct {
}

func (tl TeamCityTestLogger) TestSuiteStarted(name string) {
	fmt.Printf("##teamcity[testSuiteStarted name='%s']\n", name)
}

func (tl TeamCityTestLogger) TestStart(name string) {
	fmt.Printf("##teamcity[testStarted name='%s' captureStandardOutput='false']\n", name)
}

func (tl TeamCityTestLogger) TestIgnored(name string, message string) {
	if message == "" {
		fmt.Printf("##teamcity[testIgnored name='%s']\n", name)
	} else {
		message = sanitizeInput(message)
		fmt.Printf("##teamcity[testIgnored name='%s' message='%s']\n", name, message)
	}
}

func (tl TeamCityTestLogger) TestFinished(name string, duration int64) {
	fmt.Printf("##teamcity[testFinished name='%s' duration='%d']\n", name, duration)
}

func (tl TeamCityTestLogger) TestStdErr(name string, err string) {
	err = sanitizeInput(strings.TrimSuffix(err, "\n"))
	fmt.Printf("##teamcity[testStdErr name='%s' out='%s']\n", name, err)
}

func (tl TeamCityTestLogger) TestStdOut(name string, out string) {
	out = sanitizeInput(strings.TrimSuffix(out, "\n"))
	fmt.Printf("##teamcity[testStdOut name='%s' out='%s']\n", name, out)
}

func (tl TeamCityTestLogger) TestFailed(name string, duration int64, message, details string) {
	// ##teamcity[testFailed name='MyTest.test1' message='failure message' details='message and stack trace']
	vals := []string{
		fmt.Sprintf("name='%s'", name),
		fmt.Sprintf("duration='%d'", duration),
	}

	message = sanitizeInput(message)
	if message != "" {
		vals = append(vals, fmt.Sprintf("message='%s'", message))
	}

	details = sanitizeInput(details)
	if details != "" {
		vals = append(vals, fmt.Sprintf("details='%s'", details))
	}

	fmt.Printf("##teamcity[testFailed %s]\n", strings.Join(vals, " "))
}

func (tl TeamCityTestLogger) TestSuiteFinished(name string) {
	fmt.Printf("##teamcity[testSuiteFinished name='%s']\n", name)
}

func sanitizeInput(input string) string {
	output := input

	// from https://confluence.jetbrains.com/display/TCD9/Build+Script+Interaction+with+TeamCity#BuildScriptInteractionwithTeamCity-ReportingTests
	output = strings.ReplaceAll(output, "|", "||")
	output = strings.ReplaceAll(output, "'", "|'")
	output = strings.ReplaceAll(output, "\n", "|n")
	output = strings.ReplaceAll(output, "\r", "|r")
	output = strings.ReplaceAll(output, "[", "|[")
	output = strings.ReplaceAll(output, "]", "|]")

	return output
}
