package main

import (
	"flag"
	"fmt"
	"os"
)

const (
	CatCommand    = "cat"
	FormatCommand = "format"
)

func run() error {
	formatCommands := flag.NewFlagSet(FormatCommand, flag.ExitOnError)
	catCommands := flag.NewFlagSet(CatCommand, flag.ExitOnError)

	defFormatArgs(formatCommands)
	defCatArgs(catCommands)

	if len(os.Args) < 2 {
		return fmt.Errorf("must specify '%s' or '%s' as first arg", FormatCommand, CatCommand)
	}
	switch os.Args[1] {
	case FormatCommand:
		if err := formatCommands.Parse(os.Args[2:]); err != nil {
			return err
		}
		return runFormat()
	case CatCommand:
		if err := catCommands.Parse(os.Args[2:]); err != nil {
			return err
		}
		return runCat()
	default:
		return fmt.Errorf("unknown command %s", os.Args[1])
	}
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		os.Exit(1)
	}
}
