package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

type objdumpFinder struct{
	objdumpPath string
	functionNames []string
}

type OdbFinder interface{
	FindFunctions(filePath string) (map[string][]string, error)
	FindFunctionsByAddrs(filePath string, addrs []string) ([]string, error)
}

func NewObjdumpFinder(objdumpPath string, functionNames []string) OdbFinder {
	return &objdumpFinder{
		functionNames: functionNames,
		objdumpPath:   objdumpPath,
	}
}

func (of *objdumpFinder) FindFunctions(filePath string) (map[string][]string, error) {
	results := make(map[string][]string)

	for _, fn := range of.functionNames {
		cmdStr := fmt.Sprintf("%s -D '%s' -j .text | grep -B 10 '%s'", of.objdumpPath, filePath, fn)
		cmd := exec.Command("bash", "-c", cmdStr)

		var stderr bytes.Buffer
		cmd.Stderr = &stderr

		outPipe, err := cmd.StdoutPipe()
		if err != nil {
			return nil, fmt.Errorf("stdout pipe error: %v", err)
		}

		if err := cmd.Start(); err != nil {
			return nil, fmt.Errorf("cmd start error: %v, stderr: %s", err, stderr.String())
		}

		var lines []string
		scanner := bufio.NewScanner(outPipe)
		for scanner.Scan() {
			lineStr := scanner.Text()
			lines = append(lines, strings.TrimSpace(lineStr))
		}
		if scanErr := scanner.Err(); scanErr != nil {
			return nil, fmt.Errorf("scanner read error: %v", scanErr)
		}

		err = cmd.Wait()
		if err != nil {
			if exitErr, ok := err.(*exec.ExitError); ok && exitErr.ExitCode() == 1 && len(lines) == 0 {
				continue
			}
			return nil, fmt.Errorf("cmd wait error: %v, stderr: %s", err, stderr.String())
		}

		if len(lines) > 0 {
			results[fn] = lines
		}
	}

	return results, nil
}

func (of *objdumpFinder) FindFunctionsByAddrs(filePath string, addrs []string) ([]string, error) {
	cmd := exec.Command("bash", "-c", "readelf -sW '"+filePath+"'")
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	type funcSym struct {
		start uint64
		end   uint64
		name  string
	}

	var funcs []funcSym
	scanner := bufio.NewScanner(bytes.NewReader(out))
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) < 8 {
			continue
		}
		value := fields[1]
		size := fields[2]
		typ := fields[3]
		name := fields[len(fields)-1]

		if typ != "FUNC" {
			continue
		}

		start, err1 := strconv.ParseUint(value, 16, 64)
		sz, err2 := strconv.ParseUint(size, 10, 64)
		if err1 != nil || err2 != nil {
			continue
		}
		end := start + sz
		funcs = append(funcs, funcSym{start, end, name})
	}

	results := make([]string, len(addrs))
	for i, addr := range addrs {
		target, err := strconv.ParseUint(strings.TrimPrefix(addr, "0x"), 16, 64)
		if err != nil {
			results[i] = ""
			continue
		}

		found := ""
		for _, f := range funcs {
			if target >= f.start && target < f.end {
				found = f.name
				break
			}
		}
		results[i] = found
	}

	return results, nil
}