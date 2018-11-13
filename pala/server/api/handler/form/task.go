package form

import (
	"fmt"
	"github.com/Unknwon/com"
	"github.com/gin-gonic/gin"
)

type TailFileForm struct {
	LogPath string `form:"log_path" json:"log_path" binding:"required"`
	Size    int64  `form:"pid" json:"pid"`
}

const (
	DEFAULT_TAIL_SIZE_STR = "20480"
	DEFAULT_TAIL_SIZE     = 20480
)

func NewTailFileForm(_c *gin.Context) (*TailFileForm, error) {
	logPath := _c.DefaultQuery("log_path", "")
	if logPath == "" {
		return nil, fmt.Errorf("必须输查看日志路径")
	}
	sizeStr := _c.DefaultQuery("size", DEFAULT_TAIL_SIZE_STR)
	size, err := com.StrTo(sizeStr).Int64()
	if err != nil {
		size = DEFAULT_TAIL_SIZE
	}
	size = -size

	return &TailFileForm{
		LogPath: logPath,
		Size:    size,
	}, nil
}
