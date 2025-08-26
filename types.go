package main

type ArchType string

const (
	AMD64 ArchType = "amd64"
	ARM64 ArchType = "arm64"
)

type Args struct {
	Path    string
	Arch    string
	Funcs   []string
	ArgNums []int
	Out     string
	Worker  int
}

type ScanTarget struct {
	Funcs   []string
	ArgNums []int
	Arch    ArchType
}

type funcCall struct {
	caller, callee, argument, filename, offset string
}
