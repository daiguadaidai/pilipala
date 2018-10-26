package message

type ResponseMessage struct {
	Code    int                    `json:"code"`    // 返回码
	Message string                 `json:"message"` // 返回信息
	Data    map[string]interface{} `json:"data"`    // 存储的数据
}

func NewResponseMessage() *ResponseMessage {
	responseMessage := new(ResponseMessage)
	responseMessage.Data = make(map[string]interface{})
	responseMessage.Code = 20000

	return responseMessage
}
