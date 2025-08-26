package main

import (
	"bufio"
	"bytes"
	"os/exec"
	"strconv"
	"strings"
)

type stringFinder struct {
	targetAddrs map[int64]struct{}
}

type StringFinder interface {
	FindStrings(filePath string, outChan chan []byte) error
}

func NewStringFinder(addrs []int64) StringFinder {
	addrMap := make(map[int64]struct{})
	for _, a := range addrs {
		addrMap[a] = struct{}{}
	}
	return &stringFinder{
		targetAddrs: addrMap,
	}
}

func (sf *stringFinder) FindStrings(filePath string, outChan chan []byte) error {
	defer close(outChan)

	cmd := exec.Command("strings", "-a", "-td", filePath)
	outPipe, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	defer outPipe.Close()

	if err := cmd.Start(); err != nil {
		return err
	}

	reader := bufio.NewReader(outPipe)
	for {
		line, err := reader.ReadBytes('\n')
		if err != nil {
			break
		}
		line = bytes.TrimRight(line, "\r\n")
		parts := strings.SplitN(string(line), " ", 3)
		if len(parts) < 3 {
			continue
		}
		addr, err := strconv.ParseInt(parts[1], 10, 64)
		if err != nil {
			continue
		}

		if _, ok := sf.targetAddrs[addr]; ok {
			outChan <- []byte(parts[2])
		}
	}

	return cmd.Wait()
}
