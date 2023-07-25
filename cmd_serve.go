package main

import (
	"fmt"
	"github.com/juxuny/webhook/config"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
	"log"
	"net/http"
	"os"
)

var (
	configFile string
)

// validate config and set default value
func validateConfig(inputConfig *config.Config) error {
	nameMapper := make(map[string]bool)
	for i, deployment := range inputConfig.Deployments {
		if deployment.BashInterpreter == "" {
			inputConfig.Deployments[i].BashInterpreter = "/bin/bash"
		}
		// check if deployment name is empty
		if deployment.Name == "" {
			log.Fatalf("Empty name found in deployments[%d]\n", i)
		}

		// check if the deployment name is already used
		if nameMapper[deployment.Name] {
			log.Fatal("Duplicated name: " + deployment.Name)
		}
		nameMapper[deployment.Name] = true

		// deployment scripts is not allow empty
		if len(deployment.Scripts) == 0 {
			log.Fatal("scripts is empty")
		}

		// validate working dir
		if _, err := os.Stat(deployment.WorkDir); os.IsNotExist(err) {
			log.Fatalf("not found directory: %s\n", deployment.WorkDir)
		}
	}

	if inputConfig.Port <= 0 {
		log.Fatal("invalid listen port: ", inputConfig.Port)
	}
	return nil
}

func initConfig() {
	fileContent, err := os.ReadFile(configFile)
	if err != nil {
		log.Fatal(err)
	}
	readingConfig := &config.Config{}
	err = yaml.Unmarshal(fileContent, readingConfig)
	if err != nil {
		log.Fatal(err)
	}
	if err = validateConfig(readingConfig); err != nil {
		log.Fatal(err)
	}
	config.Set(readingConfig)
}

var serveCommand = &cobra.Command{
	Use: "serve",
	Run: func(cmd *cobra.Command, args []string) {
		initConfig()
		c := config.Get()
		address := fmt.Sprintf(":%d", c.Port)
		log.Println("listen " + address)
		err := http.ListenAndServe(address, NewHandler())
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	serveCommand.PersistentFlags().StringVarP(&configFile, "config", "c", "config.yaml", "YAML config file")
	rootCommand.AddCommand(serveCommand)
}
