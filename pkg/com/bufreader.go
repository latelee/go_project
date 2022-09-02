package com

/*
数据解析基础工具函数

默认为小端格式，大端函数名加BE

TODO：有符号数、浮点数应该要加，但当前没需求，后续加。
记录(待后续再测)：
使用json记录值，ReadUint32BE函数返回似乎只能输出十六进制，用int转，输出十进制
*/
import (
	"fmt"
	//"io"
	"strconv"
	//"bytes"
	"encoding/hex"
	//"encoding/binary"
)

type BufferReader struct {
	Buffer []byte
	offset int
	length int
}

func NewBufferReader(buf []byte) *BufferReader {
	return &BufferReader{
		Buffer: buf,
		length: len(buf),
		offset: 0,
	}
}

func NewBufferReaderWithLen(buf []byte, len int) *BufferReader {
	return &BufferReader{
		Buffer: buf,
		length: len,
		offset: 0,
	}
}

func (b *BufferReader) Init(buf []byte) {
	b.Buffer = buf
	b.offset = 0
}

func (b *BufferReader) Dump() {
	fmt.Printf("in BufferReader (len %d)\n", b.length)
	Dump(b.Buffer, b.length)
	//fmt.Printf("%v\n", hex.Dump(b.Buffer))
}

func (b *BufferReader) SkipBytes(n int) {
	b.offset += n

	return
}

func (b *BufferReader) backBytes(n int) {
	b.offset -= n

	return
}

// 返回当前的字符，但偏移量不递增
func (b *BufferReader) peekByte(o byte) {
	o = b.Buffer[b.offset]

	return
}

// 读取剩下所有字节
func (b *BufferReader) ReadBytesLeft() (o []byte) {
	o = b.Buffer[b.offset:]

	//b.offset += n;
	b.offset = b.length

	return
}

func (b *BufferReader) LeftLength() (o int) {
	o = b.length - b.offset
	return
}

func (b *BufferReader) ReadUint8() (o uint8) {
	//fmt.Printf("%x", b.Buffer);

	o = b.Buffer[b.offset]
	b.offset += 1
	return
}

func (b *BufferReader) ReadInt8() (o int8) {
	//fmt.Printf("%x", b.Buffer);

	o = int8(b.Buffer[b.offset])
	b.offset += 1
	return
}

func (b *BufferReader) ReadUint16() (o uint16) {
	var u1, u2 uint16
	u1 = uint16(b.Buffer[b.offset])
	u2 = uint16(b.Buffer[b.offset+1])
	o = u2<<8 | u1
	b.offset += 2

	return
}

func (b *BufferReader) ReadUint16BE() (o uint16) {
	var u1, u2 uint16
	u2 = uint16(b.Buffer[b.offset])
	u1 = uint16(b.Buffer[b.offset+1])
	o = u2<<8 | u1
	b.offset += 2

	return
}

func (b *BufferReader) ReadUint32() (o uint32) {
	var u1, u2, u3, u4 uint32
	u1 = uint32(b.Buffer[b.offset])
	u2 = uint32(b.Buffer[b.offset+1])
	u3 = uint32(b.Buffer[b.offset+2])
	u4 = uint32(b.Buffer[b.offset+3])
	o = (u4 << 24) | (u3 << 16) | (u2 << 8) | u1

	b.offset += 4

	return
}

func (b *BufferReader) ReadUint32BE() (o uint32) {
	var u1, u2, u3, u4 uint32
	u4 = uint32(b.Buffer[b.offset])
	u3 = uint32(b.Buffer[b.offset+1])
	u2 = uint32(b.Buffer[b.offset+2])
	u1 = uint32(b.Buffer[b.offset+3])
	o = (u4 << 24) | (u3 << 16) | (u2 << 8) | u1

	b.offset += 4

	return
}

func (b *BufferReader) ReadUint64() (o uint64) {
	var u1, u2, u3, u4, u5, u6, u7, u8 uint64
	u1 = uint64(b.Buffer[b.offset])
	u2 = uint64(b.Buffer[b.offset+1])
	u3 = uint64(b.Buffer[b.offset+2])
	u4 = uint64(b.Buffer[b.offset+3])
	u5 = uint64(b.Buffer[b.offset+4])
	u6 = uint64(b.Buffer[b.offset+5])
	u7 = uint64(b.Buffer[b.offset+6])
	u8 = uint64(b.Buffer[b.offset+7])
	o = (u8 << 56) | (u7 << 48) | (u6 << 40) | (u5 << 32) | (u4 << 24) | (u3 << 16) | (u2 << 8) | u1

	b.offset += 8

	return
}

func (b *BufferReader) ReadUint64BE() (o uint64) {
	var u1, u2, u3, u4, u5, u6, u7, u8 uint64
	u8 = uint64(b.Buffer[b.offset])
	u7 = uint64(b.Buffer[b.offset+1])
	u6 = uint64(b.Buffer[b.offset+2])
	u5 = uint64(b.Buffer[b.offset+3])
	u4 = uint64(b.Buffer[b.offset+4])
	u3 = uint64(b.Buffer[b.offset+5])
	u2 = uint64(b.Buffer[b.offset+6])
	u1 = uint64(b.Buffer[b.offset+7])
	o = (u8 << 56) | (u7 << 48) | (u6 << 40) | (u5 << 32) | (u4 << 24) | (u3 << 16) | (u2 << 8) | u1

	b.offset += 8

	return
}

func (b *BufferReader) ReadBytes(n int) (o []byte) {
	o = b.Buffer[b.offset : b.offset+n]

	b.offset += n

	return
}

// 转为十六进制格式的字符串
func (b *BufferReader) ReadHexString(n int) (o string) {
	buf := b.Buffer[b.offset : b.offset+n]
	o = hex.EncodeToString(buf)

	b.offset += n

	return
}

// 只能读取正常的字符串，即可打印的
func (b *BufferReader) ReadString(n int) (o string) {
	o = string(b.Buffer[b.offset : b.offset+n])

	// 去除末尾的\0 ???
	// o = bytes.Trim(o, "\x00")
	b.offset += n

	return
}

// 读取bcd码，如十六进制的0x20，转换十进制的20
func (b *BufferReader) ReadBCD() (o int) {
	b1 := hex.EncodeToString(b.Buffer[b.offset : b.offset+1])
	o1, _ := strconv.ParseInt(b1, 10, 8)
	o = int(o1)

	b.offset += 1

	return
}

// 与HexString相同，只是原始字符表现形式不同
func (b *BufferReader) ReadBCDString(n int) (o string) {
	o = hex.EncodeToString(b.Buffer[b.offset : b.offset+n])
	b.offset += n

	return
}

// n为1、2、4
func (b *BufferReader) ReadBCDNumber(n int) (o int) {
	b1 := hex.EncodeToString(b.Buffer[b.offset : b.offset+n])
	o1, _ := strconv.ParseInt(b1, 10, 8*n)

	o = int(o1)

	b.offset += n

	return
}

func None() {
	fmt.Println("none...")
}

/*
// 转入十六进制形式的字符，输出为byte类型
func ToHexByte(str string) (ob []byte) {
    ob, _ = hex.DecodeString(str);

    return;
}

// 转入十六进制形式的数组，输出为对应的字符
// 如 4c 77数组，将转换成4c77字符串，可保存到文件
func ToHexString(b []byte) (ostr string) {
    ostr = hex.EncodeToString(b);
    return;
}
*/
