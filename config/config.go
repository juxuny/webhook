package config

import "sync"

type Config struct {
	Port        int          `yaml:"port"`
	UrlPrefix   string       `yaml:"url-prefix"`
	Auth        []Auth       `yaml:"auth"`
	Deployments []Deployment `yaml:"deployments"`
}

type AuthType string

type Auth struct {
	Name  string `yaml:"name"`
	Token string `yaml:"token"`
}

type Deployment struct {
	Name            string     `yaml:"name"`
	WorkDir         string     `yaml:"work-dir"`
	BashInterpreter string     `yaml:"bash-interpreter"`
	Scripts         []string   `yaml:"scripts"`
	Variables       []Variable `yaml:"variables"`
}

type VariableType string

const (
	VariableTypeString  VariableType = "string"
	VariableTypeNumber  VariableType = "number"
	VariableTypeInteger VariableType = "integer"
	VariableTypeInt64   VariableType = "int64"
	VariableTypeBool    VariableType = "bool"
)

type Variable struct {
	Name     string       `yaml:"name"`
	Type     VariableType `yaml:"type"`
	Nullable bool         `yaml:"nullable"`
}

func (t Variable) IsRequired() bool {
	return !t.Nullable
}

var (
	globalConfig     *Config
	globalConfigLock = &sync.RWMutex{}
)

func Set(c *Config) {
	globalConfigLock.Lock()
	defer globalConfigLock.Unlock()
	globalConfig = c
}

func Get() *Config {
	globalConfigLock.RLock()
	defer globalConfigLock.RUnlock()
	return globalConfig
}

func (t *Config) GetAuthByToken(token string) (ret Auth, found bool) {
	for _, item := range t.Auth {
		if item.Token == token {
			return item, true
		}
	}
	return ret, false
}

func (t *Config) GetDeployment(name string) (ret Deployment, found bool) {
	if t == nil {
		return ret, false
	}
	for _, item := range t.Deployments {
		if item.Name == name {
			return item, true
		}
	}
	return
}

func (t *Config) GetHookPrefix() string {
	if t == nil {
		return ""
	}
	return t.UrlPrefix
}
