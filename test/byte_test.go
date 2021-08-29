package test

import (
	"encoding/binary"
	"fmt"
	"net_transfer/define"
	"strconv"
	"testing"
)

func TestByte(t *testing.T) {
	bytes := []byte("123")
	fmt.Println(len(bytes))
}

func TestInt(t *testing.T) {
	i, _ := strconv.Atoi("004")
	fmt.Println(i)
}

func TestByte2(t *testing.T) {
	bytes := []byte("$\r\n")
	fmt.Println(len(bytes))
}

func TestRange(t *testing.T) {
	fileList := []string{
		"a", "b",
	}
	for _, filePath := range fileList {
		fmt.Println(filePath)
	}
}

func TestIntToStr(t *testing.T) {
	fmt.Println(string(1))
}

func TestBinary(t *testing.T) {
	buffer := make([]byte, 1024)
	binary.BigEndian.PutUint32(buffer[4:8], uint32(1023))

	numBytesUint32 := binary.BigEndian.Uint32(buffer[4:8])
	fmt.Println(numBytesUint32)
}

func TestByteEq(t *testing.T) {
	headBuffer := make([]byte, 4)
	bs := []byte("000111111")
	copy(headBuffer[:4], bs[:4])
	switch string(headBuffer) {
	case string(define.DATA_FILE_INFO):
		fmt.Println("传输文件信息")
	case string(define.DATA_FILE_BODY):
		fmt.Println("传输文件体")
	case string(define.DATA_FILE_END):
		fmt.Println("文件传输完成")
	default:
		fmt.Println("未知协议")
		return
	}
}
