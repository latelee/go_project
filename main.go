package main

import (
	_ "fmt"
	"os"
	rootCmd "webdemo/cmd"
)

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
