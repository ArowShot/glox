package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/arowshot/glox"
)

func main() {
	if len(os.Args) == 2 {
		runFile(os.Args[1])
	} else if len(os.Args) == 1 {
		runPrompt()
	} else {
		fmt.Println("Usage: glox [script]")
	}
}

func runFile(path string) error {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	//fmt.Println(string(b))
	run(string(b))
	return nil
}

func runPrompt() {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("> ")
	for scanner.Scan() {
		run(scanner.Text())
		fmt.Print("> ")
	}
}

func run(script string) {
	scanner := glox.Scanner{
		Source: script,
	}

	scanner.ScanTokens()

	// for _, token := range scanner.Tokens {
	// 	fmt.Println(token)
	// }

	parser := glox.Parser{
		Tokens: scanner.Tokens,
	}
	stmts := parser.Parse()

	interpreter := &glox.Interpreter{
		Env: &glox.Environment{
			Values: make(map[string]interface{}),
		},
	}
	//astprinter := AstPrinter{}

	//fmt.Println(astprinter.print(expression))
	interpreter.Interpret(stmts)
	//fmt.Println()
}
