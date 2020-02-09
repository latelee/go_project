
package main

import (
	"fmt"
    "os"
    rootCmd "github.com/latelee/myproject/app3/cmd"
)

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
