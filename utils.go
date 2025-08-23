package main

import "fmt"

func printArgs(args Args) {
	fmt.Println("Path   :", args.Path)
	fmt.Println("Arch   :", args.Arch)
	fmt.Println("Funcs  :", args.Funcs)
	fmt.Println("Args   :", args.ArgNums)
	fmt.Println("Output :", args.Out)
	fmt.Println("Worker :", args.Worker)
}
