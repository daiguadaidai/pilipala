package server

import (
	"net/http"
	"github.com/cihub/seelog"
	"sync"
	"github.com/daiguadaidai/pilipala/pili/config"
	"github.com/gin-gonic/gin"
	"github.com/daiguadaidai/pilipala/pili/server/handler"
	"time"
	"strings"
	"fmt"
)

func StartHttpServer(_wg *sync.WaitGroup) {
	defer _wg.Done()

	// 注册路由
	router := gin.Default()
	router.Use(Cors())
	handler.Register(router)

	// 获取pala启动配置信息
	piliStartConfig := config.GetPiliStartConfig()
	s := &http.Server{
		Addr:           piliStartConfig.PiliAddress(),
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	seelog.Infof("Pili监听地址为: %v", piliStartConfig.PiliAddress())
	err := s.ListenAndServe()
	if err != nil {
		seelog.Errorf("pili启动服务出错: %v", err)
	}
}

// 跨域
func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method  //请求方法
		origin := c.Request.Header.Get("Origin")  //请求头部
		var headerKeys []string  // 声明请求头keys
		for k, _ := range c.Request.Header {
			headerKeys = append(headerKeys, k)
		}

		headerStr := strings.Join(headerKeys, ", ")
		if headerStr != "" {
			headerStr = fmt.Sprintf("access-control-allow-origin, access-control-allow-headers, %s", headerStr)
		} else {
			headerStr = "access-control-allow-origin, access-control-allow-headers"
		}

		if origin != "" {
			c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
			// 这是允许访问所有域
			c.Header("Access-Control-Allow-Origin", "*")
			//服务器支持的所有跨域请求的方法,为了避免浏览次请求的多次'预检'请求
			c.Header("Access-Control-Allow-Methods", "*")
			//  header的类型
			c.Header("Access-Control-Allow-Headers", "*")
			//  允许跨域设置,  // 跨域关键设置 让浏览器可以解析
			c.Header("Access-Control-Expose-Headers", "*")
			// 缓存请求信息 单位为秒
			c.Header("Access-Control-Max-Age", "172800")
			//  跨域请求是否需要带cookie信息 默认设置为true
			c.Header("Access-Control-Allow-Credentials", "false")
			// 设置返回格式是json
			c.Set("content-type", "application/json")
		}

		// 放行所有OPTIONS方法
		if method == "OPTIONS" {
			c.JSON(http.StatusOK, "Options Request!")
		}

		c.Next() // 处理请求
	}
}

