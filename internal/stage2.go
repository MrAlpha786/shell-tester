package internal

import (
	"regexp"

	"github.com/codecrafters-io/shell-tester/internal/shell_executable"
	"github.com/codecrafters-io/shell-tester/internal/test_cases"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

func testMissingCommand(stageHarness *test_case_harness.TestCaseHarness) error {
	logger := stageHarness.Logger
	shell := shell_executable.NewShellExecutable(stageHarness)

	if err := shell.Start(); err != nil {
		return err
	}

	testCase := test_cases.SingleLineExactMatchTestCase{
		Command:                    "nonexistent",
		FallbackPatterns:           []*regexp.Regexp{regexp.MustCompile(`^(bash: )?nonexistent: (command )?not found$`)},
		ExpectedPatternExplanation: "nonexistent: command not found",
		SuccessMessage:             "Received command not found message",
	}

	if err := testCase.Run(shell, logger); err != nil {
		return err
	}

	// At this stage the user might or might not have implemented a REPL to print the prompt again, so we won't test further

	return nil
}
