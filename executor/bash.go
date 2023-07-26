package executor

import (
	"bytes"
	"github.com/pkg/errors"
	"io"
	"os/exec"
)

type bashExecutor struct {
	scripts []string
	bash    string
	Workdir string

	output      io.Writer
	errorOutput io.Writer
	logger      Logger
	Env         []string
}

func (t *bashExecutor) Exec() error {
	command := exec.Command(t.bash)
	command.Dir = t.Workdir
	command.Stdout = t.output
	command.Stderr = t.errorOutput
	buffer := bytes.NewBuffer(nil)
	for _, script := range t.scripts {
		buffer.WriteString(script + "\n")
	}
	if buffer.Len() == 0 {
		return errors.New("script file is empty")
	}
	command.Stdin = buffer
	command.Env = t.Env
	return command.Run()
}
