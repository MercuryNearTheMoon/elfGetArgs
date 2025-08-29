package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

var X86_64Registers = []string{"rdi", "rsi", "rdx", "rcx", "r8", "r9"}
var ARM64Registers = []string{"x0", "x1", "x2", "x3", "x4", "x5"}

func printArgs(args Args) {
	type field struct {
		key string
		val interface{}
	}

	fields := []field{
		{"Path", args.Path},
		{"Arch", args.Arch},
		{"Funcs", args.Funcs},
		{"Output", args.Out},
		{"Worker", args.Worker},
	}

	if argNames, err := parseRegisters(args.ArgNums, ArchType(args.Arch)); err == nil {
		fields = append(fields[:3], append([]field{{"Regs", argNames}}, fields[3:]...)...)
	}

	maxKeyLen := 0
	for _, f := range fields {
		if len(f.key) > maxKeyLen {
			maxKeyLen = len(f.key)
		}
	}

	frameWidth := 50
	fmt.Println("+" + strings.Repeat("-", frameWidth) + "+")

	for _, f := range fields {
		valStr := fmt.Sprintf("%v", f.val)
		lines := splitString(valStr, frameWidth-maxKeyLen-4)
		for i, line := range lines {
			if i == 0 {
				fmt.Printf("| %-*s: %-*s |\n", maxKeyLen, f.key, frameWidth-maxKeyLen-4, line)
			} else {
				fmt.Printf("| %-*s  %-*s |\n", maxKeyLen, "", frameWidth-maxKeyLen-4, line)
			}
		}
	}

	fmt.Println("+" + strings.Repeat("-", frameWidth) + "+")
}

func splitString(s string, width int) []string {
	words := strings.Fields(s)
	var lines []string
	var current string

	for _, w := range words {
		if len(current)+len(w)+1 > width {
			lines = append(lines, current)
			current = w
		} else {
			if current != "" {
				current += " "
			}
			current += w
		}
	}
	if current != "" {
		lines = append(lines, current)
	}
	return lines
}

func parseRegisters(regIdx []int, arch ArchType) ([]string, error) {
	var (
		targetRegs []string
		results    []string
	)

	switch arch {
	case AMD64:
		targetRegs = X86_64Registers
	case ARM64:
		targetRegs = ARM64Registers
	default:
		return nil, fmt.Errorf("Unsupported Arch Type: %s", arch)
	}

	for _, idx := range regIdx {
		results = append(results, targetRegs[idx])
	}
	return results, nil
}

func isELF(filePath string) bool {
	f, err := os.Open(filePath)
	if err != nil {
		return false
	}
	defer f.Close()

	buf := make([]byte, 4)
	n, err := f.Read(buf)
	if err != nil || n < 4 {
		return false
	}
	return buf[0] == 0x7f && buf[1] == 'E' && buf[2] == 'L' && buf[3] == 'F'
}

func hasTextSection(filePath string) bool {
	cmd := exec.Command("readelf", "-S", filePath)
	out, err := cmd.Output()
	if err != nil {
		return false
	}
	return strings.Contains(string(out), ".text")
}

func printFuncCalls(fc chan []funcCall) {
	for fs := range fc {
		for _, f := range fs {
			fmt.Printf("Caller: %s\tCallee: %s\tArgument:%s\tFilename:%s\tOffset:%s\n", f.caller, f.callee, f.argument, f.filename, f.offset)
		}
	}
}

func csvOutput(filename string, resultsCh <-chan []funcCall) {
	file, err := os.Create(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	if err := writer.Write([]string{
		"No.",
		"File",
		"Address",
		"Caller",
		"Callee",
		"Argument"}); err != nil {
		log.Fatal(err)
	}

	no := 1
	for results := range resultsCh {
		for _, v := range results {
			record := []string{
				strconv.Itoa(no),
				v.filename,
				"0x" + v.offset,
				v.caller,
				v.callee,
				v.argument,
			}
			if err := writer.Write(record); err != nil {
				log.Println("csv write error:", err)
			}
			no++

		}
		writer.Flush()
	}
}
