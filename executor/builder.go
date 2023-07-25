package executor

import (
	"fmt"
	"io"
)

const DefaultBashFilePath = "/bin/bash"

type BashExecutorBuilder struct {
	Scripts []string
	Workdir string
	Env     []string
	Logger  Logger
	Output  io.Writer
}

func NewBashBuilder() *BashExecutorBuilder {
	return &BashExecutorBuilder{}
}

func (t *BashExecutorBuilder) SetWorkdir(workdir string) *BashExecutorBuilder {
	if t == nil {
		return t
	}
	t.Workdir = workdir
	return t
}

func (t *BashExecutorBuilder) SetScripts(scripts []string) *BashExecutorBuilder {
	if t == nil {
		return t
	}
	t.Scripts = scripts
	return t
}

func (t *BashExecutorBuilder) AddScriptFile(scriptFile string) *BashExecutorBuilder {
	if t == nil {
		return t
	}
	t.Scripts = append(t.Scripts, scriptFile)
	return t
}

func (t *BashExecutorBuilder) AddVariables(variables map[string]string) *BashExecutorBuilder {
	if t == nil {
		return t
	}
	if variables == nil {
		return t
	}
	for k, v := range variables {
		t.Env = append(t.Env, fmt.Sprintf("%v=%v", k, v))
	}
	return t
}

func (t *BashExecutorBuilder) SetOutput(output io.Writer) *BashExecutorBuilder {
	if t == nil {
		return t
	}
	t.Output = output
	return t
}

func (t *BashExecutorBuilder) SetLogger(l Logger) *BashExecutorBuilder {
	if t == nil {
		return t
	}
	t.Logger = l
	return t
}
func (t *BashExecutorBuilder) Build() IExecutor {
	ret := &bashExecutor{}
	ret.bash = DefaultBashFilePath
	ret.logger = t.Logger
	ret.scripts = t.Scripts
	ret.Workdir = t.Workdir
	ret.Env = t.Env
	ret.output = t.Output
	ret.errorOutput = t.Output
	return ret
}
