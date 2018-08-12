package main

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/fotomxq/gobase/file"
	"github.com/fotomxq/gobase/ipaddr"
	"github.com/fotomxq/gobase/log"
	"github.com/fotomxq/gobase/pedometer"

	"net/http"

	"github.com/fotomxq/gobase/gin_session"
	"github.com/fotomxq/temp-chat-tool/app"
	"github.com/fotomxq/temp-chat-tool/app/serverbase"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

var (
	AppMark string
)

//启动主体框架
func main() {
	err := Run()
	if err != nil {
		log.SendError("0.0.0.0", "", AppMark+".main()", err)
		fmt.Println(err.Error())
		return
	}
}

//初始化主体进程
//return error
func Run() error {
	//启动各类子模块
	//启动ipaddr
	ipaddr.Run()
	// 启动计步器
	go pedometer.Run()

	//启动本项目维护组
	go RunAuto()

	//打开日志
	//注册路由
	router := gin.Default()
	// 注册通用中间件
	routers := router.Use(ModGlobMiddleWare())
	// 注册核心路由部分
	CoreURL(routers)
	// 注册app url
	AppURL(routers)
	//启动服务器
	router.Run(":9001")

	return nil
}

//自动维护服务
// 针对本服务的部分信息
// 更大规模的操作放到专用的服务框架内
func RunAuto() {
	//确保log目录存在
	logDir := "." + file.Sep + "log"
	err := file.CreateFolder(logDir)
	if err != nil {
		log.SendError("0.0.0.0", "", AppMark+".main()", err)
		fmt.Println(err.Error())
	}
	//创建自动程序
	for {
		//更新gin配置
		gin.DisableConsoleColor()
		logSrc := logDir + file.Sep + "server-" + time.Now().Format("2006010215.log")
		f, err := os.Create(logSrc)
		if err != nil {
			log.SendError("0.0.0.0", "", AppMark+".main()", err)
			fmt.Println(err.Error())
		}
		gin.DefaultWriter = io.MultiWriter(f)
		//1小时运行一次
		time.Sleep(time.Hour * 1)
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//基础路由处理
////////////////////////////////////////////////////////////////////////////////////////////////////////////////

//中间件
func ModGlobMiddleWare() gin.HandlerFunc {
	//为了确保线程独立
	return func(c *gin.Context) {
		//初始化反馈头
		res := serverbase.ReportAction{}
		c.Set("from", "ModGlobMiddleWare")
		//允许跨域
		origin := c.Request.Header.Get("Origin")
		c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		//构建cookie
		cookie, err := gin_session.GetMark(c)
		//如果不存在cookie
		if err != nil {
			c.Set("err", err)
			ModLogError(c)
			res.Message = "无法建立cookie。"
			c.JSON(http.StatusOK, res)
			c.Abort()
			return
		}
		if cookie == "" {
			c.Set("err", errors.New("cannot set cookie."))
			ModLogError(c)
			res.Message = "无法建立cookie。"
			c.JSON(http.StatusOK, res)
			c.Abort()
			return
		}
		//设置上下文cookie
		c.Set("cookie", cookie)
		//返回
		c.Next()
		return
	}
}

//核心处理路由部分
func CoreURL(router gin.IRoutes) {
	CoreLoginBeforeURL(router)
}

//登陆前页面
func CoreLoginBeforeURL(router gin.IRoutes) {
	//其他所有页面
	router.GET("/", func(c *gin.Context) {
		http.Redirect(c.Writer, c.Request, "/login", http.StatusFound)
	})

	//全局favicon.ico
	router.GET("/favicon.ico", CoreLogoFavicon)
}

//获取Favicon
func CoreLogoFavicon(c *gin.Context) {
	//获取文件数据
	src := "." + file.Sep + "assets" + file.Sep + "imgs" + file.Sep + "logo-mini.png"
	content, err := file.LoadFile(src)
	if err != nil {
		c.Set("err", err)
		ModLogError(c)
		c.Abort()
		return
	}
	c.Writer.Write(content)
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//核心方法集合覆盖
// 引用mod方法并覆盖，避免反复引用包的问题
////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//ModLogError
func ModLogError(c *gin.Context) {
	serverbase.ModLogError(c)
}

//ReportError
func ReportError(message string, c *gin.Context) {
	serverbase.ReportError(message, c)
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//APP处理器
////////////////////////////////////////////////////////////////////////////////////////////////////////////////

//app url
func AppURL(router gin.IRoutes) {
	//登陆用户
	router.POST("/login", RouterLogin)
	//检查用户登陆状态
	router.POST("/user/logged-on", func(c *gin.Context) {
		//初始化
		var b bool
		res := ReportCoreUserInfo{}
		//获取参数
		res.Cookie, b = AppCheckLogin(c)
		if b == false {
			return
		}
		//反馈数据
		res.Login = true
		res.Status = true
		c.JSON(http.StatusOK, res)
	})
	//获取用户列表
	router.POST("/user/list", func(c *gin.Context) {
		_, b := AppCheckLogin(c)
		if b == false {
			return
		}
		res := serverbase.ReportAction{}
		res.Login = true
		res.Status = true
		res.Data = App.GetUserList()
		c.JSON(http.StatusOK, res)
	})
	//获取消息列队
	router.POST("/message/get", func(c *gin.Context) {
		_, b := AppCheckLogin(c)
		if b == false {
			return
		}
		postToken := c.PostForm("post_token")
		res := serverbase.ReportAction{}
		res.Login = true
		res.Status = true
		res.Data = App.GetMessageList(c, postToken)
		c.JSON(http.StatusOK, res)
	})
	//发送一个消息
	router.POST("/message/send", func(c *gin.Context) {
		token, b := AppCheckLogin(c)
		if b == false {
			return
		}
		postToken := c.PostForm("post_token")
		message := c.PostForm("message")
		err := App.SendMessage(token, postToken, message)
		if err != nil {
			ReportError("发送失败.", c)
			return
		}
		res := serverbase.ReportAction{}
		res.Login = true
		res.Status = true
		c.JSON(http.StatusOK, res)
	})
}

//登陆
func RouterLogin(c *gin.Context) {
	niceName := c.PostForm("nice_name")
	token, err := App.Login(c, niceName)
	if err != nil {
		ReportError("登陆失败.", c)
		return
	}
	//返回token
	res := serverbase.ReportAction{}
	res.Login = true
	res.Status = true
	res.Data = token
	c.JSON(http.StatusOK, res)
}

//反馈信息组
type ReportCoreUserInfo struct {
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
	//cookie值
	Cookie string `json:"cookie"`
}

//检查登陆模块
func AppCheckLogin(c *gin.Context) (string, bool) {
	//检查token是否在线，否则返回失败
	token := c.MustGet("cookie").(string)
	err := App.CheckLogin(token)
	if err != nil {
		ReportError("没有登陆.", c)
		return "", false
	}
	return token, true
}
