package logger

import (
	"fmt"
	"strings"
)

type TeamCityTestResultFormatter struct{}

func (tl TeamCityTestResultFormatter) TestSuiteStarted(name string, flowId string) string {
	return fmt.Sprintf("##teamcity[testSuiteStarted name='%s' flowId='%s']\n", name, flowId)
}

func (tl TeamCityTestResultFormatter) TestStart(name string, flowId string) string {
	return fmt.Sprintf("##teamcity[testStarted name='%s' captureStandardOutput='false' flowId='%s']\n", name, flowId)
}

func (tl TeamCityTestResultFormatter) TestIgnored(name string, message string, flowId string) string {
	vals := []string{
		fmt.Sprintf("name='%s'", name),
		fmt.Sprintf("flowId='%s'", flowId),
	}

	message = sanitizeInput(message)
	if message != "" {
		vals = append(vals, fmt.Sprintf("message='%s'", message))
	}

	return fmt.Sprintf("##teamcity[testIgnored %s]\n", strings.Join(vals, " "))
}

func (tl TeamCityTestResultFormatter) TestFinished(name string, duration int64, flowId string) string {
	return fmt.Sprintf("##teamcity[testFinished name='%s' duration='%d' flowId='%s']\n", name, duration, flowId)
}

func (tl TeamCityTestResultFormatter) TestStdErr(name string, err string, flowId string) string {
	err = strings.TrimSuffix(err, "\n")
	err = sanitizeInput(err)
	if len(err) > 0 {
		return fmt.Sprintf("##teamcity[testStdErr name='%s' out='%s' flowId='%s']\n", name, err, flowId)
	}

	return ""
}

func (tl TeamCityTestResultFormatter) TestStdOut(name string, out string, flowId string) string {
	out = strings.TrimSuffix(out, "\n")
	out = sanitizeInput(out)
	if len(out) > 0 {
		return fmt.Sprintf("##teamcity[testStdOut name='%s' out='%s' flowId='%s']\n", name, out, flowId)
	}

	return ""
}

func (tl TeamCityTestResultFormatter) TestFailed(name string, message, details string, flowId string) string {
	// ##teamcity[testFailed name='MyTest.test1' message='failure message' details='message and stack trace']
	vals := []string{
		fmt.Sprintf("name='%s'", name),
		fmt.Sprintf("flowId='%s'", flowId),
	}

	message = sanitizeInput(message)
	if message != "" {
		vals = append(vals, fmt.Sprintf("message='%s'", message))
	}

	details = sanitizeInput(details)
	if details != "" {
		vals = append(vals, fmt.Sprintf("details='%s'", details))
	}

	return fmt.Sprintf("##teamcity[testFailed %s]\n", strings.Join(vals, " "))
}

func (tl TeamCityTestResultFormatter) TestSuiteFinished(name string, flowId string) string {
	return fmt.Sprintf("##teamcity[testSuiteFinished name='%s' flowId='%s']\n", name, flowId)
}

func sanitizeInput(input string) string {
	output := input

	// from https://confluence.jetbrains.com/display/TCD9/Build+Script+Interaction+with+TeamCity#BuildScriptInteractionwithTeamCity-ReportingTests
	output = strings.Replace(output, "|", "||", -1)
	output = strings.Replace(output, "'", "|'", -1)
	output = strings.Replace(output, "\n", "|n", -1)
	output = strings.Replace(output, "\r", "|r", -1)
	output = strings.Replace(output, "[", "|[", -1)
	output = strings.Replace(output, "]", "|]", -1)
	output = strings.Replace(output, "|n|n", "|n", -1)

	return output
}
