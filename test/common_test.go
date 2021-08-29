package test

import (
	"fmt"
	"net_transfer/utils"
	"os"
	"strings"
	"testing"
	"time"
)

func TestTime(t *testing.T) {
	start := time.Now().UnixNano()
	time.Sleep(time.Second * 2)
	end := time.Now().UnixNano()
	fmt.Println(end - start)
}

func TestHumanSize(t *testing.T) {
	var size int64 = 1024 * 1024 * 1024 * 2
	fmt.Printf(utils.HumanSize(size))
}

func TestFileMd5(t *testing.T) {
	fmt.Println(utils.FileMd5("/Users/super/MyDocument/CentOS-7-x86_64-Minimal-2003.iso"))
}

func TestBytesMd5(t *testing.T) {
	// 202cb962ac59075b964b07152d234b70
	fmt.Printf(utils.BytesMd5([]byte("123")))
}

func TestPathFileListInfo(t *testing.T) {
	utils.PathFileListInfo("/Users/super/Desktop/222")
}

func TestSP(t *testing.T) {
	p := "/a/b/c/d/"
	p = p[1 : len(p)-1]
	arr := strings.Split(p, string(os.PathSeparator))
	arr = arr[:len(arr)-1]
	path := fmt.Sprintf("%s%s%s", string(os.PathSeparator), strings.Join(arr, string(os.PathSeparator)), string(os.PathSeparator))
	fmt.Println(path)
}

func TestMkdir(t *testing.T) {
	arr := []string{"1", "2", "3"}
	err := utils.Mkdir(arr)
	if err != nil {
		fmt.Println("创建失败")
	} else {
		fmt.Println("创建成功")
	}
}

func TestCreateFile(t *testing.T) {
	file, _ := os.Create("a/1.txt")
	file.Write([]byte("111111"))
	file.Close()
}
