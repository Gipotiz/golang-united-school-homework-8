package main

import (
	"flag"
	"fmt"
	"io"
	"os"
)

type Arguments map[string]string

type FlagsCLI struct {
	id        string
	operation string
	item      string
	fileName  string
}

func parseArgs() Arguments {
	flags := &FlagsCLI{}

	flag.StringVar(&flags.id, "id", "", "id")
	flag.StringVar(&flags.operation, "operation", "", "operation")
	flag.StringVar(&flags.item, "item", "", "item")
	flag.StringVar(&flags.fileName, "fileName", "", "fileName")
	flag.Parse()

	return Arguments{
		"id":        flags.id,
		"operation": flags.operation,
		"item":      flags.item,
		"fileName":  flags.fileName,
	}
}

func Perform(args Arguments, writer io.Writer) error {
	if args["fileName"] == "" {
		return fmt.Errorf("-fileName flag has to be specified")
	}
	if args["operation"] == "" {
		return fmt.Errorf("-operation flag has to be specified")
	}

	file, err := os.OpenFile(args["fileName"], os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	defer file.Close()

	switch args["operation"] {
	case "add":
		return Add(args, file, writer)
	case "list":
		return List(file, writer)
	case "findById":
		return FindById(args, file, writer)
	case "remove":
		return Remove(args, file, writer)
	default:
		return fmt.Errorf("Operation %v not allowed!", args["operation"])
	}

	return nil
}

func main() {
	err := Perform(parseArgs(), os.Stdout)
	if err != nil {
		panic(err)
	}
}
