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

	scanTarget := make([]ScanTarget, len(args.Funcs))
	for i := range args.Funcs {
		scanTarget[i] = ScanTarget{
			Funcs:   args.Funcs[i],
			ArgNums: args.ArgNums[i],
		}
    }

	fmt.Println(scanTarget)

}
