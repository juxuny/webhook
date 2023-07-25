package executor

import (
	"bytes"
	"github.com/pkg/errors"
	"io"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"
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
		if strings.Index(script, "/") != 0 {
			script = path.Join(t.Workdir, script)
		}
		log.Println("load script: ", script)
		content, err := os.ReadFile(script)
		if err != nil {
			t.logger.Println(err)
			return errors.Wrap(err, "load script failed")
		}
		buffer.Write(content)
	}
	if buffer.Len() == 0 {
		return errors.New("script file is empty")
	}
	command.Stdin = buffer
	command.Env = t.Env
	return command.Run()
}
