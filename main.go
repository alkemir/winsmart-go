package main

import (
	"fmt"

	"github.com/alkemir/winsmart-go/smart"
)

func main() {
	fmt.Println("Detecting drives...")

	drives := smart.GetLogicalDrives()

	fmt.Print("Found: ")
	for _, d := range drives {
		fmt.Print(string([]byte{65 + d}), ": ")
	}
	fmt.Println()

	for _, d := range drives {
		fmt.Printf("Getting S.M.A.R.T for %s:", string([]byte{65 + d}))
		if err := smart.Read(d - 2); err != nil {
			panic(err)
		}
	}

}
