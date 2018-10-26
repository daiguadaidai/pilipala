package handler

import "github.com/gin-gonic/gin"

var handlers []Registrant

type Registrant interface {
	Register(_router *gin.Engine)
}

// 注册路由
func Register(_router *gin.Engine) {
	// 注册任务的路由
	for _, handler := range handlers {
		handler.Register(_router)
	}
}

// 添加 handler
func AddHandler(_handler Registrant) {
	if handlers == nil {
		handlers = make([]Registrant, 0, 1)
	}

	handlers = append(handlers, _handler)
}
