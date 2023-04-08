package utils

import (
	"fmt"
	"testing"
)

func TestOpenApi(t *testing.T) {

    a, b := GenerateKey("latelee")
	fmt.Println("key:", a, b)
}