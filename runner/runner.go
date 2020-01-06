package runner

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/tombuildsstuff/teamcity-go-test-json/logger"
	"github.com/tombuildsstuff/teamcity-go-test-json/models"
	"github.com/tombuildsstuff/teamcity-go-test-json/parser"
)

var loggedAtLeastOneResult = false

type Executor struct {
	logger logger.TestResultLogger
	parser parser.TestResultParser
}

func NewExecutor() Executor {
	executor := Executor{
		logger: logger.TeamCityTestResultLogger{},
	}
	executor.parser = parser.NewResultsParser(executor.logTestResult)
	return executor
}

func (e *Executor) Execute(input ExecuteInput) error {
	args := input.toArgs()
	if input.Debug {
		log.Printf("[DEBUG] Arguments: %q", strings.Join(args, " "))
	}

	workingDirectory, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("Error retrieving current working directory: %+v", err)
	}

	cmd := exec.Command("go", args...)
	cmd.Env = os.Environ()
	cmd.Dir = workingDirectory

	out, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("Error obtaining stdout: %+v", err)
	}

	errOut, err := cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("Error obtaining stderr: %+v", err)
	}

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("Error starting: %+v", err)
	}

	outScanner := bufio.NewScanner(out)
	go e.readFromScanner(outScanner, input.Debug)

	errScanner := bufio.NewScanner(errOut)
	go e.readFromScanner(errScanner, input.Debug)

	cmd.Wait()

	if !loggedAtLeastOneResult {
		return fmt.Errorf("No Tests were found/logged!")
	}

	return nil
}

func (e *Executor) logTestResult(result models.TestResult) {
	loggedAtLeastOneResult = true
	fmt.Printf(e.logger.Log(result))
}

func (e *Executor) readFromScanner(scanner *bufio.Scanner, debug bool) {
	for scanner.Scan() {
		text := scanner.Text()

		if debug {
			log.Printf("[DEBUG] Found: %q", text)
		}

		e.parser.ParseLine(text)
	}
}
