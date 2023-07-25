package main

import (
	"github.com/spf13/cobra"
	"log"
)

var rootCommand = &cobra.Command{}

func main() {
	err := rootCommand.Execute()
	if err != nil {
		log.Fatal(err)
	}
}
