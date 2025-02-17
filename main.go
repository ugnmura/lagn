package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/SushiWaUmai/lagn/core"
)

type Lagn struct {
	hadError bool
}

func createLagn() Lagn {
	return Lagn{}
}

func (client Lagn) report(line int, where string, message string) {
	fmt.Printf("[line %d] %s: %s", line, where, message)
	client.hadError = true
}

func (client Lagn) error(line int, message string) {
	client.report(line, "", message)
}

func (client Lagn) run(line string) {
	scanner := core.CreateScanner(line)
	scanner.ScanTokens()
	parser := core.CreateParser(scanner.Tokens)
	expr := parser.Parse()
	fmt.Println(expr)
}

func (client Lagn) runFile() {
	filePath := os.Args[1]
	content, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Println(err)
		return
	}

	client.run(string(content))
	if client.hadError {
		os.Exit(65)
	}
}

func (client Lagn) runPrompt() {
	bufScanner := bufio.NewScanner(os.Stdin)

	fmt.Print("> ")
	for bufScanner.Scan() {
		line := bufScanner.Text()
		client.run(line)
		client.hadError = false
		fmt.Print("> ")
	}
}

func main() {
	client := createLagn()

	if len(os.Args) > 2 {
		fmt.Println("Usage: lagn [script]")
		os.Exit(64)
	} else if len(os.Args) == 2 {
		client.runFile()
	} else {
		client.runPrompt()
	}
}
