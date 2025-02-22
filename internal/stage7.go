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

func testType2(stageHarness *test_case_harness.TestCaseHarness) error {
	// Add the random directory to PATH (where the my_exe file is created)

	randomDir, err := getRandomDirectory()
	if err != nil {
		return err
	}

	path := os.Getenv("PATH")
	logger := stageHarness.Logger
	shell := shell_executable.NewShellExecutable(stageHarness)
	shell.Setenv("PATH", fmt.Sprintf("%s:%s", randomDir, path))

	customExecutablePath := filepath.Join(randomDir, "my_exe")
	err = custom_executable.CreateExecutable(getRandomString(), customExecutablePath)
	if err != nil {
		return err
	}

	if err := shell.Start(); err != nil {
		return err
	}

	availableExecutables := []string{"cat", "cp", "mkdir", "my_exe"}

	for _, executable := range availableExecutables {
		command := fmt.Sprintf("type %s", executable)

		var expectedPath string
		if executable == "my_exe" {
			expectedPath = customExecutablePath
		} else {
			path, err := exec.LookPath(executable)
			if err != nil {
				return fmt.Errorf("CodeCrafters internal error. Error finding %s in PATH", executable)
			}

			expectedPath = path
		}

		testCase := test_cases.SingleLineExactMatchTestCase{
			Command:                    command,
			FallbackPatterns:           []*regexp.Regexp{regexp.MustCompile(fmt.Sprintf(`^(%s is )?%s$`, executable, expectedPath))},
			ExpectedPatternExplanation: fmt.Sprintf(`%s is %s`, executable, expectedPath),
			SuccessMessage:             "Received expected response",
		}
		if err := testCase.Run(shell, logger); err != nil {
			return err
		}
	}

	nonAvailableExecutables := []string{"nonexistent"}

	for _, executable := range nonAvailableExecutables {
		command := fmt.Sprintf("type %s", executable)
		testCase := test_cases.SingleLineExactMatchTestCase{
			Command:                    command,
			FallbackPatterns:           []*regexp.Regexp{regexp.MustCompile(fmt.Sprintf(`^(bash: type: )?%s: not found$`, executable))},
			ExpectedPatternExplanation: fmt.Sprintf(`%s: not found`, executable),
			SuccessMessage:             "Received expected response",
		}
		if err := testCase.Run(shell, logger); err != nil {
			return err
		}
	}

	return assertShellIsRunning(shell, logger)
}
