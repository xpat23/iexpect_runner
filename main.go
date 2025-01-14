package main

import (
	i "iexpect_go/internal"
	"os"
)

func main() {
	i.NewTestingOutput(
		i.NewTestingResults(
			i.NewTestMethods(os.Args[1]),
		),
	).Print()
}
