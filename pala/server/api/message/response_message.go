package message

type ResponseMessage struct {
	Code    int         `json:"code"`    // 返回码
	Message string      `json:"message"` // 返回信息
	Data    interface{} `json:"data"`
}

func NewResponseMessage() *ResponseMessage {
	responseMessage := new(ResponseMessage)
	responseMessage.Code = 20000

	return responseMessage
}
