package test

import (
	"fmt"
	"testing"
)

func TestShow(t *testing.T) {
	var size int64 = 1024 * 8
	fmt.Println(size)
}

//func HumanSize(size int64) string {
//	var kb int64 = 1024
//	Mb := 1024 * kb
//	GB := 1024 * Mb
//	TB := 1024 * GB
//	r := ""
//	if size <= 1024 {
//		r = fmt.Sprintf("%dB", size)
//	} else if size < Mb {
//		r = fmt.Sprintf("%.2fB", (float64)size/kb)
//	} else if size < GB {
//		r = fmt.Sprintf("%.2fMB", size/Mb)
//	} else if size < TB {
//		r = fmt.Sprintf("%.2fGB", size/GB)
//	}
//	return r
//}
