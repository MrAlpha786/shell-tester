package shell_executable

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/codecrafters-io/tester-utils/executable"
	"github.com/codecrafters-io/tester-utils/logger"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
	ptylib "github.com/creack/pty"
)

type ShellExecutable struct {
	executable    *executable.Executable
	logger        *logger.Logger
	programLogger *logger.Logger

	// Set after starting
	args      []string
	pty       *os.File
	ptyBuffer FileBuffer
}

func NewShellExecutable(stageHarness *test_case_harness.TestCaseHarness) *ShellExecutable {
	b := &ShellExecutable{
		executable:    stageHarness.NewExecutable(),
		logger:        stageHarness.Logger,
		programLogger: logger.GetLogger(stageHarness.Logger.IsDebug, "[your_program] "),
	}

	// TODO: Kill pty process?
	// stageHarness.RegisterTeardownFunc(func() { b.Kill() })

	return b
}

func (b *ShellExecutable) Start(args ...string) error {
	b.args = args
	b.logger.Infof(b.getInitialLogLine())

	cmd := exec.Command(b.executable.Path)

	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, "PS1=$ ")
	cmd.Env = append(cmd.Env, "TERM=dumb")

	pty, err := ptylib.Start(cmd)
	if err != nil {
		panic(err)
	}

	b.pty = pty
	b.ptyBuffer = NewFileBuffer(b.pty)

	return nil
}

func (b *ShellExecutable) AssertOutputMatchesRegex(regexp *regexp.Regexp) error {
	shouldStopReadingBuffer := func(buf []byte) error {
		if regexp.Match(StripANSI(buf)) {
			return nil
		} else {
			return fmt.Errorf("Expected output to match regex")
		}
	}

	if actualValue, err := b.ptyBuffer.ReadBuffer(shouldStopReadingBuffer); err != nil {
		// b.logger.Debugf("Read bytes: %q", actualValue)
		// TODO: Add regex to log message here
		return fmt.Errorf("Expected output to match regex, but got %q", string(actualValue))
	} else {
		// b.logger.Debugf("Read bytes: %q", actualValue)
		b.programLogger.Plainf("%s", string(StripANSI(actualValue)))
	}

	return nil
}

// TODO: Convert this to "AssertOutput", and "AssertNoMoreOutput"?
func (b *ShellExecutable) AssertPrompt(prompt string) error {
	shouldStopReadingBuffer := func(buf []byte) error {
		if string(StripANSI(buf)) == prompt {
			return nil
		} else {
			return fmt.Errorf("Prompt not found")
		}
	}

	if actualValue, err := b.ptyBuffer.ReadBuffer(shouldStopReadingBuffer); err != nil {
		// b.logger.Debugf("Read bytes: %q", actualValue)

		// If the user sent any output, let's print it.
		if len(actualValue) > 0 {
			b.programLogger.Plainf("%s", string(StripANSI(actualValue)))
		}

		return fmt.Errorf("Expected %q, but got %q", prompt, string(actualValue))
	} else {
		// b.logger.Debugf("Read bytes: %q", actualValue)
		b.programLogger.Plainf("%s", string(StripANSI(actualValue)))
	}

	return nil
}

func (b *ShellExecutable) SendCommand(command string) error {
	b.logger.Infof("> %s", command)

	if err := b.writeAndReadReflection(command); err != nil {
		return err
	}

	return nil
}

func (b *ShellExecutable) writeAndReadReflection(command string) error {
	b.pty.Write([]byte(command + "\n"))

	expectedReflection := command + "\r\n"
	readBytes := []byte{}

	shouldStopReadingBuffer := func(buf []byte) error {
		if string(buf) == expectedReflection {
			return nil
		} else {
			return fmt.Errorf("Expected %q, but got %q", expectedReflection, string(buf))
		}
	}

	readBytes, err := b.ptyBuffer.ReadBuffer(shouldStopReadingBuffer)
	if err != nil {
		return fmt.Errorf("Expected %q, but got %q", expectedReflection, string(readBytes))
	}

	return nil
}

func (b *ShellExecutable) getInitialLogLine() string {
	var log string

	if b.args == nil || len(b.args) == 0 {
		log = ("Running ./spawn_shell.sh")
	} else {
		log += "Running ./spawn_shell.sh"
		for _, arg := range b.args {
			if strings.Contains(arg, " ") {
				log += " \"" + arg + "\""
			} else {
				log += " " + arg
			}
		}
	}

	return log
}

// func (b *ShellExecutable) HasExited() bool {
// 	return b.executable.HasExited()
// }

// func (b *ShellExecutable) Kill() error {
// 	b.logger.Debugf("Terminating program")
// 	if err := b.executable.Kill(); err != nil {
// 		b.logger.Debugf("Error terminating program: '%v'", err)
// 		return err
// 	}

// 	b.logger.Debugf("Program terminated successfully")
// 	return nil // When does this happen?
// }

// func (b *ShellExecutable) feedStdin(command []byte) error {
// 	n, err := b.executable.StdinPipe.Write(command)
// 	b.logger.Debugf("Wrote %d bytes to stdin", n)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

// func (b *ShellExecutable) GetStdErrBuffer() *bytes.Buffer {
// 	return b.executable.StderrBuffer
// }

// func (b *ShellExecutable) GetStdOutBuffer() *bytes.Buffer {
// 	return b.executable.StdoutBuffer
// }

// func (b *ShellExecutable) Wait() (executable.ExecutableResult, error) {
// 	return b.executable.Wait()
// }
