package define

var VERSION = "20210815"

var BLOCKSIZE = 1024 * 64

var (
	DATA_FILE_INFO = []byte("0001")
	DATA_FILE_BODY = []byte("0002")
	DATA_FILE_END  = []byte("0003")
)

var (
	DATA_SEND_OK   = []byte("1001")
	DATA_SEND_FAIL = []byte("1002")
)
