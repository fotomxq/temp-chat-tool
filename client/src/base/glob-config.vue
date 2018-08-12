<template>
  <div>{{value}}</div>
</template>

<script>
  //引入url组合库
  import queryString from "query-string";
  import {Loading} from 'element-ui';

  //常用变量组合
  export let GlobConfigData = {
    //应用基础
    appName: "temp-chat-tool",
    copyright: "© fotomxq 2018",

    //远端总服务器
    apiURL: "http://localhost:9001",
    //全局cookie，该值可修改
    cookie: "",

    //登陆页面URL
    LoginURL: "/login",

    //debug
    debug: true,
    debugUsername: "",
    debugPassword: "",

    //fetch ajax 配置组件
    ajaxLoadOn: true,
    ajaxLoadTip: "",
    ajaxGetConfig: {
      credentials: "include",
      method: "get"
    },
    ajaxPostConfig: {
      credentials: "include",
      method: "post",
      headers: {
        Accept: "application/json, application/xml, text/plain, text/html, *.*",
        "Content-Type": "application/x-www-form-urlencoded; charset=utf-8"
      }
    },
    ajaxGetURL: function (url) {
      return this.apiURL + url;
    }
  };

  //全局方法结构封装
  export let GlobMethods = {
    //fetch封装
    //由于使用yarn add whatwg-fetch 安装，所以需明确封装，否则未来调用包可能出现异常
    ajax: async function (url, params) {
      let loader;
      if (GlobConfigData.ajaxLoadOn) {
        loader = Loading.service({
          lock: true,
          text: GlobConfigData.ajaxLoadTip,
          spinner: "el-icon-loading",
          background: "rgba(0, 0, 0, 0.7)"
        });
      }
      url = GlobConfigData.apiURL + url;
      return fetch(url, params).then(resolve => {
        if (GlobConfigData.ajaxLoadOn) {
          loader.close();
        }
        return resolve;
      });
    },
    ajaxGet: async function (url) {
      return this.ajax(url, GlobConfigData.ajaxGetConfig).then(resolve => {
        return resolve;
      });
    },
    ajaxGetText: async function (url) {
      return this.ajaxGet(url).then(resp => {
        return resp.text();
      });
    },
    ajaxGetJSON: async function (url) {
      return this.ajaxGet(url).then(resp => {
        return resp.json();
      });
    },
    ajaxPost: async function (url, body) {
      let init = GlobConfigData.ajaxPostConfig;
      body['cookie'] = GlobConfigData.cookie;
      init.body = queryString.stringify(body);
      return this.ajax(url, init).then(resolve => {
        return resolve;
      });
    },
    ajaxPostText: async function (url, body) {
      return this.ajaxPost(url, body).then(resp => {
        return resp.text();
      });
    },
    ajaxPostJSON: async function (url, body) {
      return this.ajaxPost(url, body).then(resp => {
        return resp.json();
      });
    },
    //通过unix时间戳获取当前时间
  };

  //登陆方法封装
  export let GlobLogged = {
    //权限组
    userAuthority: [],
    //更新时间
    lastGetBaseTime: 0,
    //cookie
    cookie: "",
    //是否登陆？
    isLogin: false,
    //数据是否异常
    isError: false,
    //异常消息
    errorMessage: "",
    //获取用户基本数据
    GetUserBaseInfo: async function () {
      let parent = this;
      return await GlobMethods
        .ajaxPostJSON("/user/logged-on",{})
        .then(function (data) {
          //初始化所有数据
          parent.isError = false;
          parent.errorMessage = "";
          parent.isLogin = false;
          //是否获取成功？
          if (!data["status"]) {
            parent.errorMessage = data["message"];
            parent.$message.error("无法获取用户信息，请检查您的网络是否正常？");
            return false;
          }
          //是否已经登陆？
          if (!data["login"]) {
            parent.$message.error(
              "尚未登陆，可能您的登陆状态超时，请重新登陆。"
            );
            parent.isLogin = false;
            return false;
          }
          //存储基础数据
          parent.cookie = data["cookie"];
          //更新时间
          parent.lastGetBaseTime = new Date().getTime();
          //反馈
          return true;
        });
    },
  };

  /**
   * 全局配置
   */
  export default {
    data: function () {
      return {
        value: "",
        config: {
          appname: GlobConfigData.appname,
          copyright: GlobConfigData.copyright
        }
      };
    },
    props: {
      name: String
    },
    methods: {
      get: function (name) {
        return this.config[name];
      }
    },
    mounted: function () {
      this.value = this.get(this.name);
    }
  };
</script>
