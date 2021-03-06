# temp-chat-tool

## 临时加密聊天工具

自动加密通讯聊天工具，服务端不存储任何聊天数据，只做转接处理。用户间聊天信息采用密钥加密，避免中间人获取数据后破解。

注意，消息仅在服务端保留10分钟。

## 使用方法

1、下载最新的releases，client、server两个压缩包。

2、server中默认发布linux64和win64两个版本，将配置文件和_linux64或_win64二进制文件，复制到对应系统下，即可运行使用。

3、client为客户端，由于采用html5开发，可放到任何环境下运行。你也可以使用electron打包为客户端使用。

## 服务描述

标识码 / 端口 / 描述

server-api          :9001   API服务

## 文件结构

这里仅介绍核心的重要文件夹结构，部分子文件夹、文件请根据文件名或文件内容注释判断。

路径 / 描述

/app 应用逻辑库

/script 脚本

    ./build     编译脚本

/server-api API服务

/client 客户端

## 编译和部署方案

1、下载本项目后编译即可，不依赖任何数据库。

2、本项目依赖:

    github.com/fotomxq/gobase

    github.com/gin-gonic/gin

## 服务端维护方法

1、关闭本程序killall server-api_linux64

2、cd /var/temp-chat-tool/

3、chmod 777 -R /var/temp-chat-tool/

2、重新启动程序 nohup ./server-api_linux64 &

## 打包本地应用

1、运行client/build.bat进行打包。

2、复制出client下的index.html、dist目录即可完成打包工作。