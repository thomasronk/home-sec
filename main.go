package main

import (
	"fmt"
	"home-sec/detect"
	"os"
)

func main() {
	if len(os.Args) > 1 {
		detect.Run(os.Args[1])
	} else {
		fmt.Println("Need to enter an argument")
		os.Exit(1)
	}
}
