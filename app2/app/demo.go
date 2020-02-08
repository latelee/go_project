package app

import (
    "fmt"
    
)

func Demo(debug bool) {
    fmt.Println("Demo....")
    
    if debug == true {
        fmt.Println("in debug mode")
    } else {
        fmt.Println("normal...")
    }
}