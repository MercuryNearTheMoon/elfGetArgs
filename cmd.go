package main

import (
	"errors"
	"fmt"

	"github.com/spf13/pflag"
)

// ParseArgs parses CLI arguments
func ParseArgs() (Args, error) {
    var a Args

    pflag.StringVarP(&a.Path, "path", "p", "", "Target directory or file to scan (required)")
    pflag.StringVarP(&a.Arch, "arch", "A", "", "Target architecture (amd64, arm64) (required)")
    pflag.StringArrayVarP(&a.Funcs, "func", "f", []string{}, "Function name to search (can repeat)")
    pflag.IntSliceVarP(&a.ArgNums, "arg", "a", []int{}, "Argument index for corresponding function (can repeat)")
    pflag.StringVarP(&a.Out, "out", "o", "", "Output CSV file (optional)")
    pflag.IntVarP(&a.Worker, "worker", "w", 4, "Number of workers (optional, default 4)")

    pflag.Parse()

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
