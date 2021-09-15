package utils

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"io"
	"log"
	"net"
	"net_transfer/define"
)

func MakeSendFileMata(fileMeta define.FileMeta) ([]byte, error) {
	buffer := make([]byte, define.BLOCKSIZE*2)
	meta, err := json.Marshal(fileMeta)
	if err != nil {
		log.Println("json序列化失败")
		return buffer, err
	}
	copy(buffer[:4], define.DATA_FILE_INFO)
	md5Byte := []byte(BytesMd5(meta))
	copy(buffer[4:36], md5Byte)
	infoSize := len(meta)
	binary.BigEndian.PutUint32(buffer[36:40], uint32(infoSize))
	copy(buffer[40:40+infoSize], meta)
	return buffer[:40+infoSize], nil
}

func GetReceiveFileMeta(conn net.Conn) (define.FileMeta, error) {
	infoByteData := make([]byte, 40)
	io.ReadFull(conn, infoByteData)
	meta := define.FileMeta{}
	if !bytes.Equal(infoByteData[:4], define.DATA_FILE_INFO) {
		err := define.InfoError{}
		return meta, &err
	}
	md5Byte := infoByteData[4:36]
	bodySize := binary.BigEndian.Uint32(infoByteData[36:40])
	bodyBuffer := make([]byte, bodySize)
	io.ReadFull(conn, bodyBuffer)
	if !bytes.Equal([]byte(BytesMd5(bodyBuffer)), md5Byte) {
		err := define.DataBrokeError{}
		return meta, &err
	}
	err := json.Unmarshal(bodyBuffer, &meta)
	if err != nil {
		return meta, err
	}
	return meta, nil
}

func ClientAuth(conn net.Conn, userName, password string) bool {
	buffer := make([]byte, define.BLOCKSIZE)
	copy(buffer[:4], define.DATA_USER_AUTH)
	auth := define.UserAuth{
		UserName: userName,
		Password: password,
	}
	userByte, err := json.Marshal(auth)
	if err != nil {
		return false
	}
	bodySize := len(userByte)
	binary.BigEndian.PutUint32(buffer[4:8], uint32(bodySize))
	copy(buffer[8:8+bodySize], userByte)
	conn.Write(buffer[:8+bodySize])
	resp := make([]byte, 4)
	io.ReadFull(conn, resp)
	if bytes.Equal(resp, define.DATA_SEND_OK) {
		return true
	}
	return false
}

func ServerVerify(conn net.Conn) bool {
	Head := make([]byte, 8)
	io.ReadFull(conn, Head)
	if bytes.Equal(Head[:4], define.DATA_USER_AUTH) {
		bodySize := binary.BigEndian.Uint32(Head[4:8])
		Body := make([]byte, bodySize)
		io.ReadFull(conn, Body)
		u := define.UserAuth{}
		err := json.Unmarshal(Body, &u)
		if err != nil {
			return false
		}
		if u.UserName == define.DEFAULT_USER && u.Password == define.DEFAULT_PASSWORD {
			return true
		}
	}
	return false
}

func CheckConnIsAlive(conn net.Conn) bool {
	_, err := conn.Read(make([]byte, 0))
	if err != io.EOF {
		return true
	}
	return false
}
