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

}
