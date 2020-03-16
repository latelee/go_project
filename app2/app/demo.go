package app

import (
    "fmt"
    "github.com/latelee/go_project/pkg/com"
)

func doit() {
    for {
        fmt.Println(".")
        com.Sleep(1000)
    }
}

func Demo(debug bool) {
    fmt.Println("Demo....")
    
    if debug == true {
        fmt.Println("in debug mode")
    } else {
        fmt.Println("normal...")
        go doit()
    }
}