package main

import "fmt"

var X86_64Registers = []string{"rdi", "rsi", "rdx", "rcx", "r8", "r9"}
var ARM64Registers = []string{"x0", "x1", "x2", "x3", "x4", "x5"}

func printArgs(args Args) {
	fmt.Println("Path   :", args.Path)
	fmt.Println("Arch   :", args.Arch)
	fmt.Println("Funcs  :", args.Funcs)
	fmt.Println("Args   :", args.ArgNums)
	fmt.Println("Output :", args.Out)
	fmt.Println("Worker :", args.Worker)
}

func parseRegisters(regIdx []int, arch string) ([]string, error) {
	var (
		targetRegs []string
		results    []string
	)

	switch arch {
	case "amd64":
		targetRegs = X86_64Registers
	case "arm64":
		targetRegs = ARM64Registers
	default:
		return nil, fmt.Errorf("Unsupported Arch Type: %s", arch)
	}

	for _, idx := range regIdx {
		results = append(results, targetRegs[idx])
	}
	return results, nil
}
