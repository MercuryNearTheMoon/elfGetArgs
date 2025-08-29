package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

func main() {
	args, err := ParseArgs()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		flag.Usage()
		os.Exit(1)
	}

	printArgs(args)

	scanTarget := ScanTarget{
		Funcs:   args.Funcs,
		ArgNums: args.ArgNums,
		Arch:    ArchType(strings.ToLower(args.Arch)),
	}

	resultsCh := make(chan []funcCall, 100)

	go func() {
		if err := ScanELF(args.Path, scanTarget, args.Worker, resultsCh); err != nil {
			fmt.Println(err)
			close(resultsCh)
			os.Exit(1)
		}
	}()

	if args.Out != "" {
		csvOutput(args.Out, resultsCh)
	} else {
		printFuncCalls(resultsCh)
	}
}
