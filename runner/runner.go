package runner

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/tombuildsstuff/teamcity-go-test-json/logger"
	"github.com/tombuildsstuff/teamcity-go-test-json/parser"
)

func (input ExecuteInput) Execute() error {
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
	go input.readFromScanner(outScanner)

	errScanner := bufio.NewScanner(errOut)
	go input.readFromScanner(errScanner)

	cmd.Wait()

	return nil
}

func (input ExecuteInput) readFromScanner(scanner *bufio.Scanner) {
	for scanner.Scan() {
		text := scanner.Text()

		if input.Debug {
			log.Printf("[DEBUG] Found: %q", text)
		}

		parseLine(text, input.Logger)
	}
}

func parseLine(line string, testLogger logger.TeamCityTestLogger) {
	parsed, err := parser.ParseLine(line)
	if err != nil {
		log.Printf("[ERROR] %+v", err)
		return
	}

	// e.g. go mod
	if parsed == nil {
		return
	}

	if err := parsed.Log(testLogger); err != nil {
		log.Printf("[ERROR] %+v", err)
		return
	}
}
