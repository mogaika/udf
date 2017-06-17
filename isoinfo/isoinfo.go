package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/mogaika/udf"
)

func printDir(spaces string, files []udf.File) {
	for _, f := range files {
		fmt.Printf("%s %-10d %s %-20s %v\n", f.Mode().String(), f.Size(), spaces, f.Name(), f.ModTime())
		if f.IsDir() {
			printDir(spaces+"   ", f.ReadDir())
		}
	}
}

func main() {
	flag.Parse()
	rdr, err := os.Open(flag.Arg(0))
	if err != nil {
		panic(err)
	}

	u := udf.NewUdfFromReader(rdr)

	printDir("", u.ReadDir(nil))
}
