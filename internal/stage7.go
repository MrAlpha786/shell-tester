package internal

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"

	"github.com/codecrafters-io/shell-tester/internal/custom_executable"
	"github.com/codecrafters-io/shell-tester/internal/shell_executable"
	"github.com/codecrafters-io/shell-tester/internal/test_cases"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

func getPath(executable string, exePath string) string {
	if executable == "my_exe" {
		return exePath
	}

	path, err := exec.LookPath(executable)
	if err != nil {
		return fmt.Sprintf(`%s[:]? not found`, executable)
	}

	return path
}

func testType2(stageHarness *test_case_harness.TestCaseHarness) error {
	// Add the random directory to PATH (where the my_exe file is created)

	randomDir, err := GetRandomDirectory()
	if err != nil {
		return err
	}

	path := os.Getenv("PATH")
	logger := stageHarness.Logger
	shell := shell_executable.NewShellExecutable(stageHarness)
	shell.Setenv("PATH", fmt.Sprintf("%s:%s", randomDir, path))

	exePath := filepath.Join(randomDir, "my_exe")
	err = custom_executable.CreateExecutable(GetRandomString(), exePath)
	if err != nil {
		return err
	}

	if err := shell.Start(); err != nil {
		return err
	}

	availableExecutables := []string{"cat", "cp", "mkdir", "my_exe"}

	for _, executable := range availableExecutables {
		command := fmt.Sprintf("type %s", executable)
		actualPath := getPath(executable, exePath)
		expectedPattern := fmt.Sprintf(`^(%s is )?%s\r\n`, executable, actualPath)
		testCase := test_cases.RegexTestCase{
			Command:                    command,
			ExpectedPattern:            regexp.MustCompile(expectedPattern),
			ExpectedPatternExplanation: fmt.Sprintf("match %q", expectedPattern),
			SuccessMessage:             "Received expected response",
		}
		if err := testCase.Run(shell, logger); err != nil {
			return err
		}
	}

	nonAvailableExecutables := []string{"nonexistent"}

	for _, executable := range nonAvailableExecutables {
		command := fmt.Sprintf("type %s", executable)
		actualPath := getPath(executable, exePath)
		expectedPattern := fmt.Sprintf(`^(bash: type: )?%s\r\n`, actualPath)
		testCase := test_cases.RegexTestCase{
			Command:                    command,
			ExpectedPattern:            regexp.MustCompile(expectedPattern),
			ExpectedPatternExplanation: fmt.Sprintf("match %q", expectedPattern),
			SuccessMessage:             "Received expected response",
		}
		if err := testCase.Run(shell, logger); err != nil {
			return err
		}
	}

	return assertShellIsRunning(shell, logger)
}
