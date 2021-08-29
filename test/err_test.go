package test

import (
	"fmt"
	"os"
	"testing"
)

func TestErr(t *testing.T) {
	var fs *os.File
	var err = fmt.Errorf("创建文件出错")
	fs, err = os.Create("a.txt")
	if err != nil {
		fmt.Println("出错了")
	}
	fs.Write([]byte("1231"))
}
