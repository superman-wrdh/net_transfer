package define

import "net"

type FileMeta struct {
	Name      string   `json:"name"`
	LocalPath string   `json:"local_path"`
	Size      int64    `json:"size"`
	Md5       string   `json:"md5"`
	Folders   []string `json:"folders"`
}

type UserAuth struct {
	UserName string `json:"user_name"`
	Password string `json:"password"`
	To       string `json:"to"`
}

type UserConn struct {
	UserInfo UserAuth
	Conn     net.Conn
}
