package serverbase

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

//通用反馈头
type ReportAction struct {
	//是否登陆
	Login bool `json:"login"`
	//状态值
	Status bool `json:"status"`
	//消息 如果发生错误，则为错误消息
	Message string `json:"message"`
	//是否为缓冲数据
	Cache bool `json:"cache"`
	//数据摘要
	SHA1 string `json:"sha1"`
	//数据个数
	Count int `json:"count"`
	//数据集合
	Data interface{} `json:"data"`
}

//登陆前错误处理模块
func ReportError(message string, c *gin.Context) {
	res := ReportAction{}
	res.Message = message
	c.JSON(http.StatusOK, res)
	c.Abort()
}

//登陆后错误处理模块
func ReportLoginError(message string, c *gin.Context) {
	res := ReportAction{}
	res.Login = true
	res.Message = message
	c.JSON(http.StatusOK, res)
	c.Abort()
}

//登陆后错误处理并记录错误信息
func ReportLoginErrorAndLog(err error, message string, c *gin.Context) {
	c.Set("err", err)
	ModLogError(c)
	res := ReportAction{}
	res.Login = true
	res.Message = message
	c.JSON(http.StatusOK, res)
	c.Abort()
}
