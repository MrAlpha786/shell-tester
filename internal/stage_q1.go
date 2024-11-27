package internal

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/codecrafters-io/shell-tester/internal/shell_executable"
	"github.com/codecrafters-io/shell-tester/internal/test_cases"
	"github.com/codecrafters-io/tester-utils/logger"
	"github.com/codecrafters-io/tester-utils/random"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

func testQ1(stageHarness *test_case_harness.TestCaseHarness) error {
	logger := stageHarness.Logger
	shell := shell_executable.NewShellExecutable(stageHarness)

	randomDir, err := GetShortRandomDirectory()
	if err != nil {
		return err
	}
	defer os.RemoveAll(randomDir)

	filePaths := []string{
		path.Join(randomDir, fmt.Sprintf("f   %d", random.RandomInt(1, 100))),
		path.Join(randomDir, fmt.Sprintf("f   %d", random.RandomInt(1, 100))),
		path.Join(randomDir, fmt.Sprintf("f   %d", random.RandomInt(1, 100))),
	}

	if err := shell.Start(); err != nil {
		return err
	}

	_, L := getRandomWordsSmallAndLarge(5, 5)
	inputs := []string{
		fmt.Sprintf(`echo '%s %s'`, L[0], L[1]),
		fmt.Sprintf(`echo %s     %s`, L[1], L[4]),
		fmt.Sprintf(`echo '%s     %s'`, L[2], L[3]),
		fmt.Sprintf(`cat '%s' '%s' '%s'`, filePaths[0], filePaths[1], filePaths[2]),
	}
	expectedOutputs := []string{
		fmt.Sprintf("%s %s", L[0], L[1]),
		fmt.Sprintf("%s %s", L[1], L[4]),
		fmt.Sprintf("%s     %s", L[2], L[3]),
		// fmt.Sprintf("%s %s %s", L[0], L[1], L[2]),
		`new line` + `new line` + `new     line`,
	}
	testCaseContents := newTestCaseContents(inputs, expectedOutputs)

	for _, testCaseContent := range testCaseContents {
		testCase := test_cases.SingleLineExactMatchTestCase{
			Command:        testCaseContent.Input,
			ExpectedOutput: testCaseContent.ExpectedOutput,
			SuccessMessage: "Received expected response",
		}
		if err := testCase.Run(shell, logger); err != nil {
			return err
		}
	}

	return assertShellIsRunning(shell, logger)
}

type testCaseContent struct {
	Input          string
	ExpectedOutput string
}

func newTestCaseContent(input string, expectedOutput string) testCaseContent {
	return testCaseContent{
		Input:          input,
		ExpectedOutput: expectedOutput,
	}
}

func newTestCaseContents(inputs []string, expectedOutputs []string) []testCaseContent {
	testCases := []testCaseContent{}
	for i, input := range inputs {
		testCases = append(testCases, newTestCaseContent(input, expectedOutputs[i]))
	}
	return testCases
}

func writeFile(path string, content string) error {
	return os.WriteFile(path, []byte(content), 0644)
}

func writeFiles(paths []string, contents []string, logger *logger.Logger) error {
	for i, content := range contents {
		logger.Infof("Writing file %s with content \"%s\"", paths[i], strings.TrimRight(content, "\n"))
		if err := writeFile(paths[i], content); err != nil {
			logger.Errorf("Error writing file %s: %v", paths[i], err)
			return err
		}
	}
	return nil
}

func writeFilesToDirectory(directory string, contents []string, logger *logger.Logger) ([]string, error) {
	writtenFiles := []string{}

	for _, content := range contents {
		path := path.Join(directory, fmt.Sprintf("f%d", random.RandomInt(1, 100)))
		writtenFiles = append(writtenFiles, path)
		logger.Infof("Writing file %q with content \"%s\"", path, strings.TrimRight(content, "\n"))
		if err := writeFile(path, content); err != nil {
			logger.Errorf("Error writing file %s: %v", path, err)
			return nil, err
		}
	}
	return writtenFiles, nil
}

var SMALL_WORDS = []string{"foo", "bar", "baz", "qux", "quz"}
var LARGE_WORDS = []string{"hello", "world", "test", "example", "shell", "script"}

func getRandomWordsSmallAndLarge(smallCount int, largeCount int) ([]string, []string) {
	smallWords := random.RandomElementsFromArray(SMALL_WORDS, smallCount)
	largeWords := random.RandomElementsFromArray(LARGE_WORDS, largeCount)
	return smallWords, largeWords
}

var ADJECTIVES = []string{"cute", "soft", "furry", "tiny", "cozy", "sweet", "warm", "calm"}
