package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: zenx <command> [options]")
		return
	}

	switch os.Args[1] {
	case "new":
		RunNew(os.Args[2:])
	default:
		fmt.Println("Unknown command:", os.Args[1])
	}
}
