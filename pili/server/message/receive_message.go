package message

type TailLogMessage struct {
	Code    int         `json:"code"`    // 返回码
	Message string      `json:"message"` // 返回信息
	Data    interface{} `json:"data"`
}

func NewTailLogMessage() *TailLogMessage {
	msg := new(TailLogMessage)
	msg.Code = 20000

	return msg
}
