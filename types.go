package main

type Args struct {
	Path    string
	Arch    string
	Funcs   []string
	ArgNums []int
	Out     string
	Worker  int
}

type ScanTarget struct {
	Funcs   string
	ArgNums int
}
