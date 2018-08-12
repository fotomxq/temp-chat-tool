<template>
  <el-menu :router="isRouter" :unique-opened="isUniqueOpened" :default-active="menuActive">
    <el-menu-item index="/user/center">
      <i class="mdi mdi-home"></i>
      <span slot="title">聊天中心</span>
    </el-menu-item>
  </el-menu>
</template>

<script>
  export default {
    data: function () {
      return {
        //基础设定
        // 是否为路由结构
        isRouter: true,
        // 是否仅启动一个子菜单
        isUniqueOpened: true,
        //权限组
        authority: {},
      };
    },
    props: {
      //当前激活菜单
      menuActive: {
        type: String,
        default: "/center"
      }
    },
    methods: {
      /**
       * 判断是否具备对应权限序列
       * @param authoritys {String}
       * @returns {boolean}
       * @constructor
       */
      HasAuthority: function (authoritys) {
        let authorityList = authoritys.split(",");
        //只要满足一部分，则返回成功
        for (let key in this.authority) {
          let val = this.authority[key];
          for (let thisKey in authorityList) {
            let thisVal = authorityList[thisKey];
            if (val === thisVal) {
              return true;
            }
          }
        }
        //没找到返回失败
        return false;
      },
      //更新权限列队数据
      AutoUpdate: function () {
        let parent = this;
        if (this.$globLogged.lastGetBaseTime > 0) {
          this.authority = this.$globLogged.userAuthority;
        } else {
          window.setTimeout(function () {
            parent.AutoUpdate();
          }, 500);
        }
      }
    },
    mounted: function () {
      this.AutoUpdate();
    }
  };
</script>
