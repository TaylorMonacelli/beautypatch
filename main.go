package main

import (
	"fmt"
	"os"
)

func main() {
	filePaths := []string{
		"C:\\Windows\\System32\\control.exe",
		"C:\\Windows\\System32\\windowspowershell\\v1.0\\powershell.exe",
		"C:\\Windows\\explorer.exe",
	}

	for _, path := range filePaths {
		if _, err := os.Stat(path); os.IsNotExist(err) {
			fmt.Printf("File %s does not exist\n", path)
		} else {
			fmt.Printf("File %s exists\n", path)
		}
	}
}
