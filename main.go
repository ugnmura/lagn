package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/SushiWaUmai/lagn/core"
)

func run(line string, environment core.Environment) (any, error) {
	scanner := core.CreateScanner(line)
	scanner.ScanTokens()
	parser := core.CreateParser(scanner.Tokens)
	program, err := parser.Parse()
	if err != nil {
		return nil, err
	}

	var output any
	for _, expr := range program {
		output, err = expr.Interpret(environment)
    if err != nil {
      return nil, err
    }
	}

	return output, nil
}

func runFile() {
	filePath := os.Args[1]
	content, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Println(err)
		return
	}

	environment := core.DefaultEnvironment()
	_, err = run(string(content), environment)
	if err != nil {
		fmt.Println(err)
	}
}

func runPrompt() {
	bufScanner := bufio.NewScanner(os.Stdin)

	environment := core.DefaultEnvironment()
	fmt.Print("> ")
	for bufScanner.Scan() {
		line := bufScanner.Text()
		output, err := run(line, environment)
		if err != nil {
			fmt.Println(err)
		} else {
      fmt.Println(output)
    }

		fmt.Print("> ")
	}
}

func main() {
	if len(os.Args) > 2 {
		fmt.Println("Usage: lagn [script]")
		os.Exit(64)
	} else if len(os.Args) == 2 {
		runFile()
	} else {
		runPrompt()
	}
}
