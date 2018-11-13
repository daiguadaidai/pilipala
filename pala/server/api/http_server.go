package api

import (
	"fmt"
	"github.com/cihub/seelog"
	"github.com/daiguadaidai/pilipala/pala/config"
	"github.com/daiguadaidai/pilipala/pala/server/api/handler"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"sync"
	"time"
)

func StartHttpServer(_wg *sync.WaitGroup) {
	defer _wg.Done()

	// 注册路由
	router := gin.Default()
	router.Use(Cors())
	handler.Register(router)

	// 获取pala启动配置信息
	palaStartConfig := config.GetPalaStartConfig()
	s := &http.Server{
		Addr:           palaStartConfig.PalaAddress(),
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	seelog.Infof("Pala监听地址为: %v", palaStartConfig.PalaAddress())
	err := s.ListenAndServe()
	if err != nil {
		seelog.Errorf("pala启动服务出错: %v", err)
	}
}

// 跨域
func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method               //请求方法
		origin := c.Request.Header.Get("Origin") //请求头部
		var headerKeys []string                  // 声明请求头keys
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
