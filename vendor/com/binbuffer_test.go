
package com

import (
	"fmt"
	"testing"
)

func TestBinbuffer(t *testing.T) {
    
    fmt.Println("test...")
    
    strData := "4C456400640001002009110768656C6C6F776F726C64992019"

    //var reader BufferReader; //{[]byte(ToHexByte(strData)), 0
    reader := NewBufferReader([]byte(ToHexByte(strData)))

    reader.Init([]byte(ToHexByte(strData)))
    
    fmt.Printf("0x%x\n", reader.ReadUint8())
    reader.Dump()
    
    //var writer BufferWriter
    //writer.Init(100)
    
    writer := NewBufferWriter(100)
    writer.WriteUint8(0x4c)
    writer.Dump()
}

