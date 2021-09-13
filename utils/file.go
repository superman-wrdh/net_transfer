package utils

import (
	"crypto/md5"
	"fmt"
	"io"
	"log"
	"net_transfer/define"
	"os"
	"path/filepath"
	"strings"
)

func FileMd5(FilePath string) (string, error) {
	file, err := os.Open(FilePath)
	if err != nil {
		log.Panicln("读取文件出错")
		return "", err
	}
	defer file.Close()
	md5Hash := md5.New()
	if _, err := io.Copy(md5Hash, file); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", md5Hash.Sum(nil)), nil
}

func BytesMd5(bytes []byte) string {
	hash := md5.Sum(bytes)
	return fmt.Sprintf("%x", hash)
}

func IsDir(Path string) bool {
	info, err := os.Stat(Path)
	if err != nil {
		return false
	}
	return info.IsDir()
}

func Mkdir(dirs []string) error {
	fullPath := filepath.Join(dirs...)
	if IsDir(fullPath) {
		return nil
	}
	err := os.MkdirAll(fullPath, 0777)
	if err != nil {
		return err
	}
	return nil
}

func FileList(path string) []string {
	var result []string
	err := filepath.Walk(path, func(path string, f os.FileInfo, err error) error {
		if f == nil {
			return err
		}
		if f.IsDir() {
			return nil
		}

		result = append(result, path)
		return nil
	})
	if err != nil {
		return nil
	}
	return result
}

func PathFileListInfo(localPath string) []define.FileMeta {
	var fileList []define.FileMeta
	if !IsDir(localPath) {
		meta, err := MakeFileMeta(localPath, false)
		if err == nil {
			fileList = append(fileList, meta)
		}
	} else {
		if !strings.HasSuffix(localPath, string(os.PathSeparator)) {
			localPath = fmt.Sprintf("%s/", localPath)
		}
		fileFullPathList := FileList(localPath)

		localPath = localPath[1 : len(localPath)-1]
		arr := strings.Split(localPath, string(os.PathSeparator))
		arr = arr[:len(arr)-1]
		subPath := fmt.Sprintf("%s%s%s", string(os.PathSeparator), strings.Join(arr, string(os.PathSeparator)), string(os.PathSeparator))

		for _, p := range fileFullPathList {
			meta, err := MakeFileMeta(p, false)
			if err == nil {
				dir, _ := filepath.Split(p)
				var folders []string
				if dir == localPath {
					folders = []string{}
				} else {
					p = strings.Replace(p, subPath, "", -1)
					reDir, _ := filepath.Split(p)
					if strings.HasSuffix(reDir, string(os.PathSeparator)) {
						reDir = reDir[:len(reDir)-1]
					}
					folders = strings.Split(reDir, string(os.PathSeparator))
				}
				strings.Split(p, string(os.PathSeparator))
				meta.Folders = folders
				fileList = append(fileList, meta)
			}
		}

	}
	return fileList
}

func MakeFileMeta(filePath string, md5Cal bool) (define.FileMeta, error) {
	info, err := os.Stat(filePath)
	if err != nil {
		return define.FileMeta{}, err
	}
	_, name := filepath.Split(filePath)
	fileMd5 := ""
	if md5Cal {
		fileMd5, err = FileMd5(filePath)
		if err != nil {
			return define.FileMeta{}, err
		}
	}
	meta := define.FileMeta{
		Name:      name,
		LocalPath: filePath,
		Size:      info.Size(),
		Md5:       fileMd5,
		Folders:   nil,
	}
	return meta, nil
}
