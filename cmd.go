package main

import (
	"fmt"
	"os"

	"github.com/spf13/pflag"
)

// ParseArgs parses CLI arguments
func ParseArgs() (Args, error) {
	var a Args

	pflag.StringVarP(&a.Path, "path", "p", "", "Target directory or file to scan (required)")
	pflag.StringVarP(&a.Arch, "arch", "A", "", "Target architecture (amd64, arm64) (required)")
	pflag.StringArrayVarP(&a.Funcs, "func", "f", []string{}, "Function name to search (can repeat)")
	pflag.IntSliceVarP(&a.ArgNums, "arg", "a", []int{}, "Argument index for corresponding function (can repeat, start from 0)")
	pflag.StringVarP(&a.Out, "out", "o", "", "Output CSV file (optional)")
	pflag.IntVarP(&a.Worker, "worker", "w", 4, "Number of workers (optional, default 4)")

	pflag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n\n", os.Args[0])
		fmt.Fprintln(os.Stderr, "Required flags:")
		fmt.Fprintln(os.Stderr, "  -p, --path string       Target directory or file to scan")
		fmt.Fprintln(os.Stderr, "  -A, --arch string       Target architecture (amd64, arm64)")
		fmt.Fprintln(os.Stderr, "  -f, --func string       Function name to search (can repeat, at least one required)")
		fmt.Fprintln(os.Stderr, "  -a, --arg int           Argument index for corresponding function (can repeat, start from 0, at least one required)")

		fmt.Fprintln(os.Stderr, "\nOptional flags:")
		fmt.Fprintln(os.Stderr, "  -o, --out string        Output CSV file (optional)")
		fmt.Fprintln(os.Stderr, "  -w, --worker int        Number of workers (optional, default 4)")

		fmt.Fprintln(os.Stderr, "\nExample:")
		fmt.Fprintf(os.Stderr, "  %s -p ./binaries -A amd64 -f open -a 0 -f strcmp -a 1 -o output.csv -w 8\n", os.Args[0])
	}

	pflag.Parse()
	if a.Path == "" || a.Arch == "" {
		fmt.Fprintln(os.Stderr, "Error: path and arch are required")
		pflag.Usage()
		os.Exit(1)
	} else if len(a.Funcs) == 0 || len(a.ArgNums) == 0 {
		fmt.Fprintln(os.Stderr, "Error: at least one func and arg are required")
		pflag.Usage()
		os.Exit(1)
	} else if len(a.Funcs) != len(a.ArgNums) {
		fmt.Fprintln(os.Stderr, "Error: number of func and arg must match")
		pflag.Usage()
		os.Exit(1)
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
