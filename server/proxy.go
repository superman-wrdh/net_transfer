package server

import (
	"fmt"
	"net"
	"net_transfer/define"
	"net_transfer/utils"
)

var ConnPoolMap map[string]define.UserConn
var CloseConnChan chan net.Conn

func init() {
	if ConnPoolMap == nil {
		ConnPoolMap = map[string]define.UserConn{}
	}
	CloseConnChan = make(chan net.Conn, 100)
}

func StartProxy(ip, port string) {
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
			success <- 1
			alive := utils.CheckConnIsAlive(conn)
			if alive {
				conn.Write(define.DATA_SEND_OK)
				// 连接放入连接池
				ConnPoolMap[authInfo.UserName] = define.UserConn{
					UserInfo: authInfo,
					Conn:     conn,
				}
				// 连接池中找到了连接
				if ToConn, ok := ConnPoolMap[authInfo.To]; ok {
					Proxy(ToConn.Conn, conn)
				} else {
					fmt.Println("等待连接")
				}

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

func Proxy(curr, remote net.Conn) {
	go Send(curr, remote)
	go Send(remote, curr)
}

func Send(from, to net.Conn) {
	buffer := make([]byte, 4096)
	for {
		n, err := from.Read(buffer)
		if err != nil {
			fmt.Println(err.Error())
		}
		n2, err := to.Write(buffer[:n])
		if err != nil {
			fmt.Println(err.Error())
		}
		if n != n2 {
			fmt.Println("send not completion")
		}
	}

}
