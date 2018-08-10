package main

import (
	"fmt"
	"github.com/fotomxq/gobase/app"
	"github.com/fotomxq/gobase/file"
	"github.com/fotomxq/gobase/ipaddr"
	"github.com/fotomxq/gobase/log"
	"github.com/fotomxq/gobase/pedometer"
	"io"
	"os"
	"time"

	"github.com/fotomxq/temp-chat-tool/app/serverbase"
	"github.com/gin-gonic/gin"
	"net/http"
)

//启动主体框架
func main() {
	err := Run()
	if err != nil {
		log.SendError("0.0.0.0", "", app.AppMark+".main()", err)
		fmt.Println(err.Error())
		return
	}
}

//初始化主体进程
//return error
func Run() error {
	//设置基本参数
	var err error

	//初始化底层应用结构
	// 加载配置文件
	err = app.RunLoadConfig()
	if err != nil {
		return err
	}

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
	routers := router.Use(nil)
	// 注册核心路由部分
	CoreURL(routers)
	// 注册app url
	AppURL(routers)
	//启动服务器
	router.Run(app.SystemConfig.RouterHost)

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
		log.SendError("0.0.0.0", "", app.AppMark+".main()", err)
		fmt.Println(err.Error())
	}
	//创建自动程序
	for {
		//更新gin配置
		gin.DisableConsoleColor()
		logSrc := logDir + file.Sep + "server-" + time.Now().Format("2006010215.log")
		f, err := os.Create(logSrc)
		if err != nil {
			log.SendError("0.0.0.0", "", app.AppMark+".main()", err)
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

}
