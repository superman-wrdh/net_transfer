package utils

import (
	"fmt"
)

func HumanSize(size int64) string {
	var KB int64 = 1024
	var MB = 1024 * KB
	var GB = 1024 * MB
	var TB = 1024 * GB
	res := ""
	if size < KB {
		res = fmt.Sprintf("%dB", size)
	} else if size < MB {
		res = fmt.Sprintf("%.2fKB", float64(size)/float64(KB))
	} else if size < GB {
		res = fmt.Sprintf("%.2fMB", float64(size)/float64(MB))
	} else if size < TB {
		res = fmt.Sprintf("%.2fGB", float64(size)/float64(GB))
	} else {
		res = fmt.Sprintf("%.2fTB", float64(size)/float64(GB*1024))
	}
	return res
}
