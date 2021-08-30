package client

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
	"net_transfer/define"
	"net_transfer/utils"
	"os"
	"time"
)

func HandleError(err error, when string) {
	if err != nil {
		fmt.Println(err, when)
		os.Exit(1)
	}
}

func StartClient(ip, port string) {
	addr := fmt.Sprintf("%s:%s", ip, port)
	fmt.Println("connect ", addr)
	conn, err := net.Dial("tcp", addr)
	HandleError(err, "client conn error")
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Println("请输入文件/文件夹路径(退出请输入quit)")
		lineByte, _, _ := reader.ReadLine()
		line := string(lineByte)
		if line == "q" || line == "quit" {
			return
		}
		fileList := utils.PathFileListInfo(line)
		for _, meta := range fileList {
			meta.Md5 = utils.FileMd5(meta.LocalPath)
			SendFile(conn, meta)
		}
	}

}

func SendFile(conn net.Conn, fileMata define.FileMeta) {
	fileLocalPath := fileMata.LocalPath
	file, err := os.Open(fileLocalPath)
	if err != nil {
		log.Fatalln("打开文件出错")
		return
	}
	metaByteData, err := utils.MakeSendFileMata(fileMata)
	if err != nil {
		fmt.Println("构造文件信息错误")
		return
	}
	//fmt.Println(fileMata)
	// 写入meta信息
	conn.Write(metaByteData)
	fileBuffer := make([]byte, define.BLOCKSIZE)
	size := 0

	//进度条
	var bar utils.Bar
	bar.NewOption(0, fileMata.Size)
	//END
	fmt.Printf("\n开始发送文件%s 文件大小%s\n", fileMata.Name, utils.HumanSize(fileMata.Size))
	start := time.Now().UnixNano() / 1e6
	for {
		n, err := file.Read(fileBuffer)

		if n > 0 {
			state := SendData(conn, fileBuffer[:n], "SEND_FILE")
			if bytes.Equal(state, define.DATA_SEND_OK) {
				size += n
				bar.Play(int64(size))
				//fmt.Println("发送成功")
			} else if bytes.Equal(state, define.DATA_SEND_FAIL) {
				fmt.Println("发送失败 重试10次")
				for i := 0; i < 10; i++ {
					//fmt.Printf("重试第%d中 \n", i)
					state := SendData(conn, fileBuffer[:n], "SEND_FILE")
					if bytes.Equal(state, define.DATA_SEND_OK) {
						break
					} else if bytes.Equal(state, define.DATA_SEND_OK) {
						size += n
						bar.Play(int64(size))
					}
				}
				fmt.Println("网络状况不好 传输失败")
				return
			}
		} else {
			//fmt.Println("文件读完")
			tmp := make([]byte, 4)
			SendData(conn, tmp, "SEND_FILE_END")
			break
		}

		if err == io.EOF {
			tmp := make([]byte, 4)
			SendData(conn, tmp, "SEND_FILE_END")
			break
		}
	}
	end := time.Now().UnixNano() / 1e6
	cost := end - start
	if cost == 0 {
		cost = 100
	}
	bar.Finish()
	fmt.Printf("本次文件传输完成 耗时%.2fs 速度%s/s \n", float64(cost)/float64(1000), utils.HumanSize((1000*fileMata.Size)/cost))
}

func SendData(con net.Conn, data []byte, dataType string) []byte {
	buffer := make([]byte, define.BLOCKSIZE*2)
	md5 := utils.BytesMd5(data)
	var Head []byte
	if dataType == "SEND_FILE" {
		Head = define.DATA_FILE_BODY
	} else if dataType == "SEND_FILE_END" {
		Head = define.DATA_FILE_END
	}
	headSize := len(Head)
	copy(buffer[:headSize], Head)
	md5Size := len([]byte(md5))
	copy(buffer[headSize:headSize+md5Size], []byte(md5))

	dataSize := len(data)
	binary.BigEndian.PutUint32(buffer[headSize+md5Size:headSize+md5Size+4], uint32(dataSize))

	copy(buffer[headSize+md5Size+4:headSize+md5Size+4+dataSize], data)

	con.Write(buffer[:headSize+md5Size+4+dataSize])
	if dataType == "SEND_FILE" {
		stateBuffer := make([]byte, 4)
		io.ReadFull(con, stateBuffer)
		return stateBuffer
	} else {
		return define.DATA_SEND_OK
	}

}
