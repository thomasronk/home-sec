package main

import (
	"fmt"
	"home-sec/detect"
	"os"
)

func main() {
	if err := detect.Run(); err != nil {
		fmt.Print("uh oh")
		os.Exit(1)
	}
}
