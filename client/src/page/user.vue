<template>
  <el-container class="be-wrapper be-fixed-sidebar">
    <base-header :title="headerTitle"></base-header>
    <el-container>
      <el-aside width="230px" class="menu">
        <base-menu></base-menu>
      </el-aside>
      <el-main>
        <router-view></router-view>
      </el-main>
    </el-container>
  </el-container>
</template>

<script>
  /**
   * 全局初始化
   */
  export default {
    components: {
      "base-menu": r => require.ensure([], () => r(require("../base/menu.vue")), "base-menu"),
      "base-header": r => require.ensure([], () => r(require("../base/header.vue")), "base-header"),
      "base-breadcrumb": r => require.ensure([], () => r(require("../base/breadcrumb.vue")), "base-breadcrumb"),
    },
    data: function () {
      return {
        //修正loading为真
        loading: true,
        //顶部标题
        headerTitle: "聊天中心",
      };
    },
    props: {},
    methods: {
      //注入新的URL地址
      CreateRouter: function () {
        let parent = this;
        //获取用户基础数据
        this.$globLogged.GetUserBaseInfo()
          .then(function (res) {
            //失败退出
            if (!res) {
              //进入登陆界面
              window.setTimeout(function () {
                parent.$router.push("/login");
              }, 2000);
              return false;
            }
            //遍历所有权限，构建添加子路由器
            //增加新URL处理
            parent.$router.addRoutes([
              {
                //退出登陆
                path: "/logout",
                name: "logout",
                component: function () {
                  let url = "/logout";
                  parent.$globConfig.ajaxLoadTip = "正在登陆...";
                  parent.$globMethods.ajaxPostJSON(url, {}).then(function (data) {
                    parent.$message.success("退出成功...");
                    window.setTimeout(function () {
                      parent.$router.push("/login");
                    }, 2000);
                  });
                }
              },
              //用户及子路由控制器
              //待增加...
            ]);
            //反馈
            return true;
          })
          .then(function (res) {
            if (!res) {
              return false;
            }
            //反馈
            return true;
          });
      }
    },
    //初始化
    mounted: function () {
      //初始化子路由
      this.CreateRouter();
      //启动进入center模块
      this.$router.push('/user/center');
    }
  };
</script>
