package main

import (
	"fmt"
	"os"

	"github.com/akroshchenko/cron-parser/cron"
)

func usage() {
	fmt.Println(`
	Usage: cron-parser '<cron expression>'

	Example: cron-parser "*/15 0 1,15 * 1-5 /usr/bin/find"`)
}

func main() {
	args := os.Args[1:]
	if len(args) == 0 || len(args) > 1 {
		usage()
		os.Exit(1)
	}
	cronStr := args[0]

	expr, err := cron.Parse(cronStr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "# %s: %s\n", os.Args[0], err)
		os.Exit(1)
	}

	fmt.Print(expr)
}
