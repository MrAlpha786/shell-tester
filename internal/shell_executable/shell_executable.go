package shell_executable

import (
	"strings"

	"github.com/codecrafters-io/tester-utils/executable"
	"github.com/codecrafters-io/tester-utils/logger"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

type ShellExecutable struct {
	executable *executable.Executable
	logger     *logger.Logger
	args       []string
}

func NewShellExecutable(stageHarness *test_case_harness.TestCaseHarness) *ShellExecutable {
	b := &ShellExecutable{
		executable: stageHarness.NewExecutable(),
		logger:     stageHarness.Logger,
	}

	stageHarness.RegisterTeardownFunc(func() { b.Kill() })

	return b
}

func (b *ShellExecutable) Run(args ...string) error {
	b.args = args
	if b.args == nil || len(b.args) == 0 {
		b.logger.Infof("$ ./spawn_shell.sh")
	} else {
		var log string
		log += "$ ./spawn_shell.sh"
		for _, arg := range b.args {
			if strings.Contains(arg, " ") {
				log += " \"" + arg + "\""
			} else {
				log += " " + arg
			}
		}
		b.logger.Infof(log)
	}

	if err := b.executable.Start(b.args...); err != nil {
		return err
	}

	return nil
}

func (b *ShellExecutable) HasExited() bool {
	// Call Result() before HasExited(), Result will wait until everything is synced
	// So we get the correct result from HasExited()
	return b.executable.HasExited()
}

func (b *ShellExecutable) Kill() error {
	b.logger.Debugf("Terminating program")
	if err := b.executable.Kill(); err != nil {
		b.logger.Debugf("Error terminating program: '%v'", err)
		return err
	}

	b.logger.Debugf("Program terminated successfully")
	return nil // When does this happen?
}

func (b *ShellExecutable) feedStdin(command []byte) error {
	n, err := b.executable.StdinPipe.Write(command)
	b.logger.Debugf("Wrote %d bytes to stdin", n)
	if err != nil {
		return err
	}
	return nil
}

func (b *ShellExecutable) FeedStdin(command []byte) error {
	commandWithEnter := append(command, []byte("\n")...)
	return b.feedStdin(commandWithEnter)
}

func (b *ShellExecutable) Result() (executable.ExecutableResult, error) {
	return b.executable.Result()
}

func (b *ShellExecutable) ResultWithWait() (executable.ExecutableResult, error) {
	return b.executable.ResultWithWait()
}
