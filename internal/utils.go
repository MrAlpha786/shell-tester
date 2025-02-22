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
)

var SMALL_WORDS = []string{"foo", "bar", "baz", "qux", "quz"}
var LARGE_WORDS = []string{"hello", "world", "test", "example", "shell", "script"}

func assertShellIsRunning(shell *shell_executable.ShellExecutable, logger *logger.Logger) error {
	testCase := test_cases.NewSilentPromptTestCase("$ ")

	if err := testCase.Run(shell, logger); err != nil {
		return fmt.Errorf("Expected shell to print prompt after last command, but it didn't: %v", err)
	}
	return nil
}

// getRandomDirectory creates a random directory in /tmp, creates the directories and returns the full path
// directory is of the form `/tmp/<random-word>/<random-word>/<random-word>`
func getRandomDirectory() (string, error) {
	randomDir := path.Join("/tmp", random.RandomWord(), random.RandomWord(), random.RandomWord())
	if err := os.MkdirAll(randomDir, 0755); err != nil {
		return "", fmt.Errorf("CodeCrafters internal error. Error creating directory %s: %v", randomDir, err)
	}
	return randomDir, nil
}

// getShortRandomDirectory creates a random directory in /tmp, creates the directories and returns the full path
// directory is of the form `/tmp/<random-word>`
func getShortRandomDirectory() (string, error) {
	randomDir := path.Join("/tmp", random.RandomElementFromArray(SMALL_WORDS))
	if err := os.MkdirAll(randomDir, 0755); err != nil {
		return "", fmt.Errorf("CodeCrafters internal error. Error creating directory %s: %v", randomDir, err)
	}
	return randomDir, nil
}

func getRandomString() string {
	// We will use a random numeric string of length = 6
	var result string
	for i := 0; i < 5; i++ {
		result += fmt.Sprintf("%d", random.RandomInt(10, 99))
	}

	return result
}

func getRandomName() string {
	names := []string{"Alice", "David", "Emily", "James", "Maria"}
	return names[random.RandomInt(0, len(names))]
}

// writeFile writes a file to the given path with the given content
func writeFile(path string, content string) error {
	return os.WriteFile(path, []byte(content), 0644)
}

// writeFiles writes a list of files to the given paths with the given contents
func writeFiles(paths []string, contents []string, logger *logger.Logger) error {
	for i, content := range contents {
		logger.Infof("Writing file \"%s\" with content \"%s\"", paths[i], strings.TrimRight(content, "\n"))
		if err := writeFile(paths[i], content); err != nil {
			logger.Errorf("Error writing file %s: %v", paths[i], err)
			return err
		}
	}
	return nil
}
