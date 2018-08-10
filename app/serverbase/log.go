package serverbase

import (
	"github.com/fotomxq/gobase/log"

	"github.com/gin-gonic/gin"
)

//日志处理模块
func ModLogError(c *gin.Context) {
	username, b := c.Get("username")
	if b == false {
		username = ""
	}
	from := c.MustGet("from").(string)
	err := c.MustGet("err").(error)
	log.SendError(c.ClientIP(), username.(string), from, err)
}
