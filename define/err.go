package define

import "fmt"

type InfoError struct {
	Err error
}

func (e *InfoError) Error() string {
	return fmt.Sprintf("文件信息协议错误")
}

type DataBrokeError struct {
	Err error
}

func (e *DataBrokeError) Error() string {
	return fmt.Sprintf("信息不完整,检验失败")
}
