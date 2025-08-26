package main

import (
	"fmt"
	"io/fs"
	"log"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"
)

type (
	X64Extractor   struct{}
	Arm64Extractor struct{}
	AddrExtractor  interface {
		Extract(lines []string, tarRegister string) (addrs []string, targetAddrs []int64)
	}
)

func (x X64Extractor) Extract(lines []string, tarRegister string) ([]string, []int64) {
	re := regexp.MustCompile(`#\s*([0-9a-fA-F]+)`)
	var addrs []string
	var targetAddrs []int64

	for _, line := range lines {
		if strings.Contains(line, "lea") && strings.Contains(line, tarRegister) {
			match := re.FindStringSubmatch(line)
			if len(match) < 2 {
				continue
			}
			addrHex := match[1]
			refAddr, err := strconv.ParseInt(addrHex, 16, 64)
			if err != nil {
				continue
			}
			addr := strings.TrimSpace(strings.Split(line, ":")[0])
			addrs = append(addrs, addr)
			targetAddrs = append(targetAddrs, refAddr)
		}
	}
	return addrs, targetAddrs
}

func (a Arm64Extractor) Extract(lines []string, tarRegister string) ([]string, []int64) {
	reBase := regexp.MustCompile(`adrp\s+\w+,\s+([0-9a-fA-F]+)`)
	reOffset := regexp.MustCompile(`#([0-9a-fA-Fx]+)`)

	var addrs []string
	var targetAddrs []int64
	var baseAddr int64
	var addr string
	foundBase := false

	for _, line := range lines {
		if !foundBase && strings.Contains(line, "adrp") && strings.Contains(line, tarRegister) {
			match := reBase.FindStringSubmatch(line)
			if len(match) < 2 {
				continue
			}
			baseAddr, _ = strconv.ParseInt(match[1], 16, 64)
			addr = strings.TrimSpace(strings.Split(line, ":")[0])
			foundBase = true
		} else if foundBase && strings.Contains(line, "add") && strings.Contains(line, tarRegister) {
			match := reOffset.FindStringSubmatch(line)
			if len(match) < 2 {
				continue
			}
			offsetAddr, _ := strconv.ParseInt(match[1], 0, 64)
			targetAddrs = append(targetAddrs, baseAddr+offsetAddr)
			addrs = append(addrs, addr)
			foundBase = false
		}
	}
	return addrs, targetAddrs
}

func scanFile(filePath string, scanTarget ScanTarget) []funcCall {
	functionNames := scanTarget.Funcs
	registers, err := parseRegisters(scanTarget.ArgNums, scanTarget.Arch)
	if err != nil || len(functionNames) != len(registers) {
		fmt.Println("Invalid input:", err)
		return nil
	}
	var extractor AddrExtractor
	switch scanTarget.Arch {
	case "arm64":
		extractor = Arm64Extractor{}
	case "amd64":
		extractor = X64Extractor{}
	default:
		fmt.Printf("Unsupported Arch Type: %s", scanTarget.Arch)
		return nil
	}

	objFinder := NewObjdumpFinder("objdump", functionNames)
	fnLinesMap, err := objFinder.FindFunctions(filePath)
	if err != nil {
		fmt.Println("Execute Error:", err)
		return nil
	}

	var results []funcCall

	for idx, funcName := range functionNames {
		tarRegister := registers[idx]
		lines := fnLinesMap[funcName]

		addrs, targetAddrs := extractor.Extract(lines, tarRegister)
		if len(addrs) == 0 {
			continue
		}

		strChan := make(chan []byte)
		sf := NewStringFinder(targetAddrs)
		go func(c chan []byte) {
			if err := sf.FindStrings(filePath, c); err != nil {
				fmt.Println("StringFinder Error:", err)
			}
		}(strChan)

		var strSlice [][]byte
		for str := range strChan {
			strSlice = append(strSlice, str)
		}

		callerNames, _ := objFinder.FindFunctionsByAddrs(filePath, addrs)
		for i := 0; i < len(strSlice); i++ {
			caller := ""
			if i < len(callerNames) {
				caller = callerNames[i]
			}
			results = append(results, funcCall{
				caller:   caller,
				callee:   funcName,
				argument: string(strSlice[i]),
				filename: filePath,
				offset:   addrs[i],
			})
		}
	}

	return results
}

func ScanELF(targetPath string, scanTargets ScanTarget, workerNum int) []funcCall {
	var (
		allResults []funcCall
		wg         sync.WaitGroup
	)

	filesCh := make(chan string, 100)
	resultsCh := make(chan []funcCall, 100)

	// Collector
	collectorDone := make(chan struct{})
	go func() {
		for results := range resultsCh {
			allResults = append(allResults, results...)
		}
		close(collectorDone)
	}()

	// Workers
	for i := 0; i < workerNum; i++ {
		go func() {
			for path := range filesCh {
				resultsCh <- scanFile(path, scanTargets)
				fmt.Println("Scanned ELF:", path)
				wg.Done()
			}
		}()
	}

	err := filepath.WalkDir(targetPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if !isELF(path) {
			return nil
		}
		if !hasTextSection(path) {
			fmt.Printf("Skipped %s, .text section not found\n", path)
			return nil
		}

		wg.Add(1)
		filesCh <- path
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	wg.Wait()
	close(filesCh)
	close(resultsCh)

	// wait for collector end
	<-collectorDone

	return allResults
}
