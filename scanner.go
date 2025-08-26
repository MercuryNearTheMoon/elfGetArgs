package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

func scanFile(filePath string, scanTarget ScanTarget) []funcCall {
	functionNames := scanTarget.Funcs
	registers, err := parseRegisters(scanTarget.ArgNums, scanTarget.Arch)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	if len(functionNames) != len(registers) {
		fmt.Println("functionNames and registers must be the same length")
		return nil
	}

	var results []funcCall
	re := regexp.MustCompile(`#\s*([0-9a-fA-F]+)`)

	objFinder := NewObjdumpFinder("objdump", functionNames)
	fnLinesMap, err := objFinder.FindFunctions(filePath)
	if err != nil {
		fmt.Println("Execute Error:", err)
		return nil
	}

	for idx, funcName := range functionNames {
		tarRegister := registers[idx]
		lines := fnLinesMap[funcName]

		var targetAddrs []int64
		var addrs []string
		found := false

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
				found = true
			}
		}

		if !found {
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
		tmpResults := make([]funcCall, 0, len(strSlice))
		for i := 0; i < len(strSlice); i++ {
			caller := ""
			if i < len(callerNames) {
				caller = callerNames[i]
			}
			tmpResults = append(tmpResults, funcCall{
				caller:   caller,
				callee:   funcName,
				argument: string(strSlice[i]),
				filename: filePath,
				offset:   addrs[i],
			})
		}

		results = append(results, tmpResults...)
	}

	return results
}
