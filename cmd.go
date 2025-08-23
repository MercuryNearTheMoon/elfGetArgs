package main

import (
	"errors"
	"flag"
	"fmt"
)

type Args struct {
    Path    string
    Arch    string
    Funcs   []string
    ArgNums []int
    Out     string
    Worker  int
}

// ParseArgs parses CLI arguments
func ParseArgs() (Args, error) {
    var a Args
    funcs := multiFlag{}
    args := multiIntFlag{}

    flag.StringVar(&a.Path, "path", "", "Target directory or file to scan (required)")
    flag.StringVar(&a.Arch, "arch", "", "Target architecture (amd64, arm64) (required)")
    flag.Var(&funcs, "func", "Function name to search (can be repeated)")
    flag.Var(&args, "arg", "Argument index for corresponding function (can be repeated)")
    flag.StringVar(&a.Out, "out", "", "Output CSV file (optional)")
    flag.IntVar(&a.Worker, "worker", 1, "Number of workers (optional, default 1)")
    flag.Parse()

    a.Funcs = funcs
    a.ArgNums = args

    if a.Path == "" || a.Arch == "" {
        return a, errors.New("path and arch are required")
    }
    if len(a.Funcs) == 0 || len(a.ArgNums) == 0 {
        return a, errors.New("at least one func and arg are required")
    }
    if len(a.Funcs) != len(a.ArgNums) {
        return a, errors.New("number of func and arg must match")
    }

    return a, nil
}

type multiFlag []string

func (m *multiFlag) String() string {
    return fmt.Sprintf("%v", *m)
}

func (m *multiFlag) Set(value string) error {
    *m = append(*m, value)
    return nil
}

type multiIntFlag []int

func (m *multiIntFlag) String() string {
    return fmt.Sprintf("%v", *m)
}

func (m *multiIntFlag) Set(value string) error {
    var v int
    _, err := fmt.Sscan(value, &v)
    if err != nil {
        return err
    }
    *m = append(*m, v)
    return nil
}
