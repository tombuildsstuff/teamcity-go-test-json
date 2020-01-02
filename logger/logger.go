package logger

import (
	"fmt"
	"log"
	"strings"
	"time"
)

type testData struct {
	panic    bool
	started  time.Time
	stdErr   []string
	stdOut   []string
	messages []string
}

func (td *testData) append(message string) {
	td.messages = append(td.messages, message)
}

func (td *testData) appendErr(message string) {
	td.stdErr = append(td.stdErr, message)
}

func (td *testData) appendOut(message string) {
	if td.panic {
		td.appendErr(message)
		return
	}

	td.stdOut = append(td.stdOut, message)
}

type TeamCityTestLogger struct {
	// TODO: we're going to have to keep track of which tests have been/are being run here
	// since we can only output one StdOut and StdErr line per test
	tests map[string]*testData
}

func NewTeamCityTestLogger() TeamCityTestLogger {
	return TeamCityTestLogger{
		tests: map[string]*testData{},
	}
}

func (tl TeamCityTestLogger) TestSuiteStarted(name string) {
	fmt.Printf("##teamcity[testSuiteStarted name='%s']\n", name)
}

func (tl TeamCityTestLogger) TestStart(name string) {
	tl.tests[name] = &testData{
		stdErr:  []string{},
		stdOut:  []string{},
		panic:   false,
		started: time.Now(),
	}
}

func (tl TeamCityTestLogger) TestIgnored(name string, message string) {
	if message == "" {
		tl.tests[name].append(fmt.Sprintf("##teamcity[testIgnored name='%s']\n", name))
	} else {
		message = sanitizeInput(message)
		tl.tests[name].append(fmt.Sprintf("##teamcity[testIgnored name='%s' message='%s']\n", name, message))
	}

	tl.TestFinished(name, -1)
}

func (tl TeamCityTestLogger) TestFinished(name string, duration int64) {
	test, ok := tl.tests[name]
	if ok {
		fmt.Printf("##teamcity[testStarted name='%s' captureStandardOutput='false']\n", name)

		// anything else
		log.Printf(strings.Join(test.messages, "\n"))

		// we can only have one StdErr/StdOut per test
		out := strings.Join(test.stdOut, "\n")
		out = strings.TrimSuffix(out, "\n")
		out = sanitizeInput(out)
		fmt.Printf("##teamcity[testStdOut name='%s' out='%s']\n", name, out)

		err := strings.Join(test.stdErr, "\n")
		err = strings.TrimSuffix(err, "\n")
		err = sanitizeInput(err)
		fmt.Printf("##teamcity[testStdErr name='%s' out='%s']\n", name, err)
	}

	// output the stderr/stdout
	fmt.Printf("##teamcity[testFinished name='%s' duration='%d']\n", name, duration)

	delete(tl.tests, name)
}

func (tl TeamCityTestLogger) TestStdErr(name string, err string) {
	tl.tests[name].appendErr(err)
}

func (tl TeamCityTestLogger) TestStdOut(name string, out string) {
	if strings.HasPrefix(out, "panic:") {
		tl.tests[name].panic = true
	}

	tl.tests[name].appendOut(out)
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

	tl.tests[name].append(fmt.Sprintf("##teamcity[testFailed %s]\n", strings.Join(vals, " ")))
	tl.TestFinished(name, duration)
}

func (tl TeamCityTestLogger) TestSuiteFinished(name string) {
	fmt.Printf("##teamcity[testSuiteFinished name='%s']\n", name)
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
