package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/SushiWaUmai/lagn/core"
)

func run(line string, environment core.Environment) error {
	scanner := core.CreateScanner(line)
	scanner.ScanTokens()
	parser := core.CreateParser(scanner.Tokens)
	program, err := parser.Parse()
	if err != nil {
		return err
	}

	var output any
	for _, expr := range program {
		output = expr.Interpret(environment)
	}
	fmt.Println(output)

	return nil
}

func runFile() {
	filePath := os.Args[1]
	content, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Println(err)
		return
	}

	environment := core.DefaultEnvironment()
	err = run(string(content), environment)
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
		err := run(line, environment)
		if err != nil {
			fmt.Println(err)
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
