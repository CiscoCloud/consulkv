package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/mitchellh/cli"
)

func main() {
	log.SetOutput(ioutil.Discard)

	cli := &cli.CLI{
		Args:		os.Args[1:],
		Commands:	Commands,
		HelpFunc:	cli.BasicHelpFunc("consulkv"),
	}

	exitCode, err := cli.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error executing CLI: %s\n", err.Error())
		os.Exit(1)
	}

	os.Exit(exitCode)
}
