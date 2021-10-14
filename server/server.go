package server

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net_transfer/define"
	"net_transfer/utils"
	"os"
	"strings"
	"time"
)

func StartServer(ip, port string) {

	addr := fmt.Sprintf("%s:%s", ip, port)
	fmt.Println("start listen server ", addr)
	listener, err := net.Listen("tcp", addr)
	ServerHandleError(err, "net.listen")

	for {
		conn, e := listener.Accept()
		ServerHandleError(e, "listener.accept")
		fmt.Printf("%s connected \n", conn.RemoteAddr().String())

		success := make(chan int, 1)
		callback := make(chan int, 1)
		go TimeOutCheck(conn, success, callback)

		authInfo, verify := utils.ServerVerify(conn)
		if verify {
			fmt.Printf("用户认证成功%s\n", authInfo.UserName)
			success <- 1
			alive := utils.CheckConnIsAlive(conn)
			if alive {
				conn.Write(define.DATA_SEND_OK)
				go FileHandler(conn)
			} else {
				fmt.Printf("客户方在规定时间内没有认证 连接关闭")
			}
		} else {
			fmt.Printf("%s 验证失败 \n", conn.RemoteAddr().String())
			conn.Write(define.DATA_SEND_FAIL)
			conn.Close()
		}

	}
}

func ServerHandleError(err error, when string) {
	if err != nil && err != io.EOF {
		fmt.Println(err, when)
		os.Exit(1)
	}
}

func FileHandler(conn net.Conn) {
	var fs *os.File
	var fileSize = 0
	var bar utils.Bar
	var start int64 = 0
	var end int64 = 0
	meta := define.FileMeta{}
	tmpFileName := ""

	for {
		headBuffer := make([]byte, 40)
		io.ReadFull(conn, headBuffer)
		md5Bytes := headBuffer[4:36]

		bodySize := binary.BigEndian.Uint32(headBuffer[36:40])
		var body []byte
		if bodySize > 0 {
			body = make([]byte, bodySize)
			io.ReadFull(conn, body)
		}
		switch string(headBuffer[:4]) {

		case string(define.DATA_FILE_INFO):
			// 接收文件meta信息
			if bodySize == 0 {
				break
			}
			json.Unmarshal(body, &meta)
			var err = fmt.Errorf("创建文件出错")
			fmt.Printf("开始接收文件:%s 文件大小%s\n", meta.Name, utils.HumanSize(meta.Size))
			start = time.Now().UnixNano() / 1e6
			bar.NewOption(0, meta.Size)
			tmpFileName = fmt.Sprintf("%s.tmp", meta.Name)
			if len(meta.Folders) > 0 {
				// 创建文件夹
				err := utils.Mkdir(meta.Folders)
				if err != nil {
					return
				}
				// 文件夹 / 文件名
				tmpFileName = fmt.Sprintf("%s%s%s",
					strings.Join(meta.Folders,
						string(os.PathSeparator)),
					string(os.PathSeparator), tmpFileName)
			}
			fs, err = os.Create(tmpFileName)
			if err != nil {
				return
			}

		case string(define.DATA_FILE_BODY):
			// 接收文件体
			if bodySize == 0 {
				break
			}
			if bytes.Equal([]byte(utils.BytesMd5(body)), []byte(md5Bytes)) {
				fileSize += int(bodySize)
				fs.Write(body)
				conn.Write(define.DATA_SEND_OK)
				bar.Play(int64(fileSize))
				//fmt.Printf("文件已经写入%d \n", fileSize)
			} else {
				conn.Write(define.DATA_SEND_FAIL)
			}

		case string(define.DATA_FILE_END):
			// 文件传输完成
			end = time.Now().UnixNano() / 1e6
			cost := end - start
			bar.Finish()
			if cost == 0 {
				// 小文件传输速度很快 花费时间为0 就算100ms
				cost = 100
			}
			fileMd5, err := utils.FileMd5(tmpFileName)
			if err == nil {
				if fileMd5 == meta.Md5 {
					// 将临时文件重命名成正式文件
					err := os.Rename(tmpFileName, tmpFileName[:len(tmpFileName)-4])
					if err == nil {
						fmt.Printf("MD5校验通过,本次文件传输完成 耗时%.2fs 速度%s/s \n", float64(cost)/float64(1000), utils.HumanSize((1000*meta.Size)/cost))
					} else {
						fmt.Printf("文件重命名失败,%v", err.Error())
					}

				} else {
					fmt.Println("MD5校验不通过 文件损坏")
				}
			} else {
				fmt.Printf("计算文件md5失败,%v", err.Error())
			}

			start = 0
			end = 0

			bar.Finish()
			bar.Reset()

			fs.Close()
			fs = nil

			fileSize = 0

		default:
			//fmt.Println("未知协议")
			return
		}
	}
}

func TimeOutCheck(conn net.Conn, success, callback chan int) {
	ticker := time.NewTicker(20 * time.Second)
	defer ticker.Stop()
	select {
	case <-success:
		fmt.Printf("%s 完成验证 \n", conn.RemoteAddr().String())
		callback <- 1
		return
	case <-ticker.C:
		callback <- 0
		fmt.Printf("%s 客户端超时未完成验证 \n", conn.RemoteAddr().String())
		conn.Close()
		return
	}
}
