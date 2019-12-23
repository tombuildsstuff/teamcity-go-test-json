package runner

import (
	"fmt"
	"strings"
)

func (input ExecuteInput) toArgs() []string {
	args := []string{
		"test",
		"-v",
	}

	if input.Scope != "" {
		scope := input.Scope
		if !strings.HasSuffix(scope, "/...") {
			scope = fmt.Sprintf("%s/...", scope)
		}
		args = append(args, scope)
	} else {
		args = append(args, "./...")
	}

	if input.Prefix != "" {
		args = append(args, fmt.Sprintf("-run=%s", input.Prefix))
	}

	if input.Parallelism != 0 {
		args = append(args, fmt.Sprintf("-test.parallel=%d", input.Parallelism))
	} else {
		args = append(args, "-test.parallel=1")
	}

	if input.Count != 0 {
		args = append(args, fmt.Sprintf("-count=%d", input.Count))
	} else {
		args = append(args, "-count=1")
	}

	if input.TimeoutInHours != 0 {
		timeoutInMinutes := input.TimeoutInHours * 60
		timeoutInSeconds := timeoutInMinutes * 60
		args = append(args, fmt.Sprintf("-timeout=%ds", timeoutInSeconds))
	} else {
		args = append(args, "-timeout=3600s")
	}

	args = append(args, "-json")

	return args
}
