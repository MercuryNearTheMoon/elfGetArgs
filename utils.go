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

func printFuncCalls(fs []funcCall) {
	for _, f := range fs {
		fmt.Printf("Caller: %s\tCallee: %s\tArgument:%s\tFilename:%s\tOffset:%s\n", f.caller, f.callee, f.argument, f.filename, f.offset)
	}
}

func csvOutput(filename string, result []funcCall) {

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
	var writeErr error
	no := 1
	for _, v := range result {
		record := []string{
			strconv.Itoa(no),
			v.filename,
			"0x" + v.offset,
			v.caller,
			v.callee,
			v.argument}
		writeErr = writer.Write(record)
		no++
		if writeErr != nil {
			log.Fatal(writeErr)
		}
	}
}
