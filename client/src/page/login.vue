<template>
  <div class="max-size">
    <el-main>
      <div class="splash-container">
        <div class="card card-border-color card-border-color-primary">
          <div class="card-header">
            <img src="../assets/imgs/logo.png" alt="LOGO" width="102" height="27" class="logo-img">
            <span class="splash-description">登陆账户</span>
          </div>
          <div class="card-body button-widest">
            <el-form ref="form" :model="form">
              <el-form-item>
                <el-input name="host" v-model="form.host" placeholder="服务器"></el-input>
              </el-form-item>
              <el-form-item>
                <el-input name="nice_name" v-model="form.nice_name" placeholder="用户昵称"></el-input>
              </el-form-item>
              <el-form-item class="ch-size-100">
                <el-row>
                  <el-col :span="24">
                    <el-button type="primary" v-on:click="Login">登陆</el-button>
                  </el-col>
                </el-row>
              </el-form-item>
            </el-form>
          </div>
        </div>
        <div class="splash-footer">
          <p>
            <glob-config name="copyright"></glob-config>
          </p>
        </div>
      </div>
    </el-main>
  </div>
</template>
<script>
  export default {
    components: {},
    data: function () {
      return {
        //基础数据
        form: {
          // 服务器
          host: this.$globConfig.debug ? this.$globConfig.apiURL : "",
          // 用户名
          nice_name: this.$globConfig.debug ? this.$globConfig.debugUsername : "",
        }
      };
    },
    methods: {
      //登陆
      Login: function (event) {
        if(! this.form.nice_name){
          this.$message.error(
            "请输入正确的用户昵称..."
          );
        }
        this.$globConfig.apiURL = this.form.host;
        let parent = this;
        let url = "/login";
        let body = {
          nice_name: this.form.nice_name,
        };
        this.$globConfig.ajaxLoadTip = '正在登陆...';
        this.$globMethods.ajaxPostJSON(url, body).then(function (data) {
          if (data["status"] === true) {
            parent.$message.success("登陆成功，正在进入...");
            parent.$globConfig.cookie = data['data'];
            window.setTimeout(function () {
              parent.$router.push('/user');
            }, 2000);
            return;
          }
          parent.$message.error(
            "登陆失败，请检查您的用户名、密码、验证码是否正确？"
          );
        });
      },
      //显示注册
      goReg: function (event) {
        this.$router.push("/reg");
      },
      //显示忘记密码
      goForget: function (event) {
        this.$router.push("/forget");
      },
      //当验证码改变
      OnVcodeChange: function (val) {
        this.form.vcode = val;
      }
    }
  };
</script>
