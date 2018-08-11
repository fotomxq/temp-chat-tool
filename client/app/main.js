//全局模块
const Electron = require('electron');
const App = Electron.app;
const BrowserWindow = Electron.BrowserWindow;
const Menu = Electron.Menu;
const Tray = Electron.Tray;
const Session = Electron.session;
const IPCMain = Electron.ipcMain;
const FS = require('fs');
const Path = require('path');

/////////////////////////////////////////////////////////////////////////////////////////////////////////////
// 全局设置
/////////////////////////////////////////////////////////////////////////////////////////////////////////////

let GlobConfigs = {
    "Mark": "temp-chat-tool",
    "Token": "",
    "Username": "",
    "Password": "",
    "Server-URL": "./index.html",
    "Debug-Server-URL": "./index.html",
    "Version": "0.0.1"
};

let Debug = true;

//主窗口
let MainWin;

/////////////////////////////////////////////////////////////////////////////////////////////////////////////
// 控制器主体设置
/////////////////////////////////////////////////////////////////////////////////////////////////////////////

//控制器主体
let MainApp = new Object;

MainApp.CreateWindow = function() {
    //注意，正式版应该启动缓冲功能
    if (Debug === true) {
        App.commandLine.appendSwitch('--disable-http-cache');
    }
    // 查询所有 cookies.
    const ses = Session.fromPartition('persist:name');
    ses.clearStorageData([], function(data) {});
    Session.defaultSession.cookies.get({}, function(error, cookies) {
        for (key in cookies) {
            var val = cookies[key];
            Session.defaultSession.cookies.remove('http://' + val['domain'], val['name'], function() {});
            Session.defaultSession.cookies.remove('https://' + val['domain'], val['name'], function() {});
        }
    });

    //构建托盘
    iconPath = 'logo.ico';
    appIcon = new Tray(iconPath);
    const contextMenu = Menu.buildFromTemplate([{
        label: '显示程序',
        click: function() {
            MainWin.show();
        }
    }, {
        label: '退出程序',
        click: function() {
            App.Quit();
        }
    }]);
    appIcon.setToolTip('临时加密通讯工具');
    appIcon.setContextMenu(contextMenu);

    //获取屏幕尺寸
    let electronScreen = Electron.screen;
    let size = electronScreen.getPrimaryDisplay().workAreaSize;

    //构建主程序框架
    let params = {
        width: size.width - 100,
        height: size.height - 50,
        minWidth: 960,
        minHeight: 600,
        nodeIntegration: 'iframe',
        webPreferences: {
            webSecurity: true,
            allowDisplayingInsecureContent: true
        }
    };
    //正式发布时，请使该配置生效
    //frame: false,
    if (Debug === false) {
        params['frame'] = false;
    }
    MainWin = new BrowserWindow(params);

    //禁止多开
    const shouldQuit = App.makeSingleInstance((commandLine, workingDirectory) => {
        // Someone tried to run a second instance, we should focus our window.
        if (MainWin) {
            if (MainWin.isMinimized()) {
                MainWin.restore();
            }
            MainWin.focus();
        }
    });
    if (shouldQuit) {
        App.quit();
        return false;
    }

    /**
     * MainWin 其他配置信息
     *  这是最初使用的配置信息，但发现配置可能无效，为了确保延续性，保留旧的配置信息。
     *  "node-integration": "iframe",
     *  "web-preferences": {
     *  "web-security": false
     *  }
     */
    let gotoURL = GlobConfigs['Server-URL'];
    if (Debug) {
        gotoURL = GlobConfigs['Debug-Server-URL'];
    }

    MainWin.loadFile(gotoURL);

    //是否显示调试工具
    if (Debug === true) {
        MainWin.webContents.openDevTools();
    }

    //是否显示菜单栏
    // 正式发布时，请使该配置生效
    if (Debug === false) {
        //考虑到稳定性，去掉该选项
        Menu.setApplicationMenu(null);
    }

    MainWin.on('closed', function() {
        MainWin = null;
    });

    //启动日志存储
    Log.Auto();
};

//关闭程序
MainApp.Quit = function() {
    MainWin = null;
    App.quit();
};

/////////////////////////////////////////////////////////////////////////////////////////////////////////////
// 加载应用程序
/////////////////////////////////////////////////////////////////////////////////////////////////////////////

//程序加载完成后，开始创建主体框架
App.on('ready', MainApp.CreateWindow);

//当窗口全部关闭，自动退出程序
App.on('window-all-closed', function() {
    if (process.platform !== 'darwin') {
        MainApp.Quit();
    }
});

//当程序加载完成后，启动框架主体
App.on('activate', function() {
    if (MainWin === null) {
        MainApp.CreateWindow();
    }
});

/////////////////////////////////////////////////////////////////////////////////////////////////////////////
// 异步通讯技术
/////////////////////////////////////////////////////////////////////////////////////////////////////////////
//关闭程序申请
IPCMain.on('close-app', (event, name) => {
    MainApp.Quit();
});

//隐藏到托盘
IPCMain.on('hide-app', (event, name) => {
    MainWin.hide();
});

//最大化
IPCMain.on('window-max-app', (event, name) => {
    MainWin.maximize();
});

//还原
IPCMain.on('window-return-app', (event, name) => {
    MainWin.unmaximize();
});

//最小化
IPCMain.on('window-hide-app', (event, name) => {
    MainWin.minimize();
});

/////////////////////////////////////////////////////////////////////////////////////////////////////////////
// 日志模块
/////////////////////////////////////////////////////////////////////////////////////////////////////////////
//日志
let Log = new Object;

//日志存储文件夹
Log.Dir = './log';
//定时器
Log.Timer = '';
//日志列队数据
Log.Data = [];

//支持类型
Log.ERROR = 'error';
Log.MESSAGE = 'message';
Log.SYSTEM = 'system';

/**
 * 接收一个日志
 */
Electron.ipcMain.on('send-log', (event, logType, content) => {
    Log.Send(logType, content);
});

/**
 * 发送日志
 * @param logType string 日志类型
 * @param content string 日志内容
 * @constructor
 */
Log.Send = function(logType, content) {
    let date = new Date();
    let nowTime = AppCore.GetDateTimeUTCFormat(date, 'yyyy-MM-dd_hh:mm:ss:SS');
    let logContent = nowTime + ' ' + '[' + logType + '] ' + content + '\r\n';
    Log.Data[Log.Data.length] = logContent;
    console.log(logContent);
};

//自动维护程序
// 定时存储日志数据
Log.Auto = function() {
    //初始化
    let date = new Date();
    let logFileName = AppCore.GetDateTimeUTCFormat(date, 'yyyy-MM-dd') + '.log';
    //确保日志目录存在
    let isDir = AppCore.CreateDir(Log.Dir);
    //如果创建失败，则稍后再试
    if (!isDir) {
        //1秒后重试
        Log.Timer = setTimeout(function() {
            Log.Auto();
        }, 1000);
        //不继续后续内容
        return;
    }
    //如果存在数据，则遍历
    let logSrc = Log.Dir + '/' + logFileName;
    if (Log.Data) {
        let content = '';
        for (let key in Log.Data) {
            let val = Log.Data[key];
            if (val) {
                content += val;
            }
        }
        Log.Data = [];
        if (content) {
            FS.appendFileSync(Path.join(logSrc), content, function(err) {
                if (err) {
                    console.log(err);
                }
            });
        }
    }
    //将日志发送到服务器
    // 500毫秒后继续运行
    Log.Timer = setTimeout(function() {
        Log.Auto();
    }, 500);
    //不继续后续内容
    return;
};

/////////////////////////////////////////////////////////////////////////////////////////////////////////////
// 基础模块
/////////////////////////////////////////////////////////////////////////////////////////////////////////////
var AppCore = new Object;

/**
 * 将日期时间转为Unix时间戳
 * @param date Date 时间句柄
 * @param format string 时间结构 qq yyyy-MM-dd hh:mm:ss:SS
 * @return string 日期时间 eg : 2018-05-03 10:02
 **/
AppCore.GetDateTimeUTCFormat = function(date, format) {
    var o = {
        "M+": date.getMonth() + 1, //月份
        "d+": date.getDate(), //日
        "h+": date.getHours(), //小时
        "m+": date.getMinutes(), //分
        "s+": date.getSeconds(), //秒
        "q+": Math.floor((date.getMonth() + 3) / 3), //季度
        "S": date.getMilliseconds() //毫秒
    };
    if (/(y+)/.test(format))
        format = format.replace(RegExp.$1, (date.getFullYear() + "").substr(4 - RegExp.$1.length));
    for (var k in o)
        if (new RegExp("(" + k + ")").test(format))
            format = format.replace(RegExp.$1, (RegExp.$1.length == 1) ? (o[k]) : (("00" + o[k]).substr(("" + o[k]).length)));
    return format;
};

/**
 * 确保目录被创建
 * @param src string 文件夹路径
 * @returns {*}
 * @constructor
 */
AppCore.CreateDir = async function(src) {
    return await FS.stat(src, async function(err, stats) {
        if (err) {
            return await FS.mkdir(Path.join(src), await
                function(err) {
                    if (err) {
                        console.log(err);
                        return false;
                    }
                    return true;
                });
        }
        if (!await stats.isDirectory()) {
            return await FS.mkdir(Path.join(src), await
                function(err) {
                    if (err) {
                        console.log(err);
                        return false;
                    }
                    return true;
                });
        }
        return false;
    });
};