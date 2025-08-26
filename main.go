package main

import (
	"flag"
	"fmt"
	"os"
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
		Arch:    args.Arch,
	}

	fmt.Println(scanTarget)

	results := ScanELF(args.Path, scanTarget, args.Worker)

	if args.Out == "" {
		printFuncCalls(results)
	} else {
		csvOutput(args.Out, results)
	}

}
