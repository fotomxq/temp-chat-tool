<template>
  <el-main>
    <el-row>
      <el-col :span="24">
        <el-breadcrumb separator="/">
          <el-breadcrumb-item :to="{ path: '/user/center' }">聊天中心</el-breadcrumb-item>
        </el-breadcrumb>
      </el-col>
    </el-row>
    <el-row>
      <el-col :span="24">
      </el-col>
      <el-col :span="24">
        <el-row :gutter="20">
          <el-col :span="4">
            <el-card shadow="always">
              <el-table :data="UserList" style="width: 100%" @row-click="selectUser">
                <el-table-column prop="nice_name" label="在线用户（单击选择）"></el-table-column>
              </el-table>
            </el-card>
          </el-col>
          <el-col :span="20">
            <el-row :gutter="20">
              <el-col :span="24">
                <el-card shadow="always">
                  <el-row :gutter="20">
                    <el-col :span="20">
                      <el-input type="input" placeholder="加密和消息解密 密钥 必须是数字组成..." v-model="NowUserKey"></el-input>
                    </el-col>
                    <el-col :span="4">
                      <el-button @click="randKey">生成密钥</el-button>
                    </el-col>
                    <el-col :span="20">
                      <el-input type="textarea" :rows="2" placeholder="消息内容..." v-model="MessageContent"></el-input>
                    </el-col>
                    <el-col :span="4">
                      <el-button type="primary" @click="sendMessage">发送消息</el-button>
                    </el-col>
                  </el-row>
                </el-card>
              </el-col>
              <el-col :span="24">
                <el-card shadow="hover">
                  <h5>正在与{{NowUserInfo.nice_name}}对话...</h5>
                  <ul>
                    <li v-for="(item, index) in MessageList">
                      {{item.nice_name}} - {{item.create_time}} : {{item.content}}
                    </li>
                  </ul>
                </el-card>
              </el-col>
            </el-row>
          </el-col>
        </el-row>
      </el-col>
    </el-row>
  </el-main>
</template>

<script>
  //引入加密库
  import sjcl from '../assets/vendor/sjcl/sjcl.js';

  export default {
    data: function () {
      return {
        //刷新用户列表定时器
        timerRefUserList: null,
        //刷新消息列表定时器
        timerRefMessageList: null,
        //用户列表
        UserList: [],
        //用户对应的密钥
        // token => value
        UserListKey: {},
        //当前聊天对象用户信息
        NowUserInfo: {},
        //当前聊天用户key
        NowUserKey: '',
        //发送消息内容
        MessageContent: '',
        //消息列队
        MessageList: [],
      }
    },
    methods: {
      //刷新用户列表
      refUserList: function(){
        let parent = this;
        let url = "/user/list";
        let body = {};
        this.$globMethods.ajaxPostJSON(url, body).then(function (data) {
          if (data["status"] === true) {
            //parent.$message.success("刷新用户列表成功...");
            parent.UserList = data['data'];
            return;
          }
          parent.$message.error(
            "无法获取用户列表数据...请确保网络正常或是否已登陆？"
          );
        });
      },
      //选择用户
      selectUser: function(row){
        //找到用户信息
        this.NowUserInfo = {
          nice_name: row['nice_name'],
          token: row['token'],
        };
        let nowKey = '';
        for(let key in this.UserList){
          let val = this.UserList[key];
          for(let keyK in this.UserListKey){
            let valK = this.UserListKey[keyK];
            if(val['token'] === keyK){
              nowKey = valK;
              break;
            }
          }
          if(nowKey){
            this.NowUserKey = nowKey;
            break;
          }
        }
        this.MessageList = [];
        this.autoRefMessageList();
      },
      //发送消息
      sendMessage: function(){
        if(!this.NowUserInfo['token']){
          this.$message.error(
            "必须选择用户后才能给对方发送信息..."
          );
          return false;
        }
        let parent = this;
        let url = "/message/send";
        let body = {
          post_token: this.NowUserInfo['token'],
          message: sjcl.encrypt(this.NowUserKey,this.MessageContent),
        };
        this.$globMethods.ajaxPostJSON(url, body).then(function (data) {
          if (data["status"] === true) {
            parent.$message.success("发送成功...");
            return;
          }
          parent.$message.error(
            "发送失败..."
          );
        });
      },
      //随机生成key
      randKey: function(){
        this.NowUserKey = (Math.random() * 1000000000000000000).toString() + (Math.random() * 1000000000000000000).toString() + (Math.random() * 1000000000000000000).toString() + (Math.random() * 1000000000000000000).toString() + (Math.random() * 1000000000000000000).toString() + (Math.random() * 1000000000000000000).toString();
      },
      //刷新消息列表
      refMessageList: function(){
        let parent = this;
        //必须选择一个用户
        if(!this.NowUserInfo['token']){
          return false;
        }
        //发送请求获取消息
        let url = "/message/get";
        let body = {
          post_token: this.NowUserInfo['token'],
        };
        this.$globMethods.ajaxPostJSON(url, body).then(function (data) {
          if (data["status"] === true) {
            //parent.$message.success("获取消息成功...");
            //重构消息列队
            parent.MessageList = [];
            parent.MessageList = data['data'].reverse();
            for(let key in parent.MessageList){
              let val = parent.MessageList[key];
              parent.MessageList[key]['nice_name'] = parent.getUserNiceName(val['send_user_token']);
              parent.MessageList[key]['content'] = sjcl.decrypt(parent.NowUserKey,val['content']);
            }
            return;
          }
          parent.$message.error(
            "获取消息失败..."
          );
        });
      },
      //通过token获取用户名称
      getUserNiceName: function(token){
        for(let key in this.UserList){
          let val = this.UserList[key];
          if(val['token'] === token){
            return val['nice_name'];
          }
        }
        return '';
      },
      //自动刷新消息列表
      autoRefMessageList: function(){
        let parent = this;
        //定时刷新消息列表
        parent.refMessageList();
        clearInterval(this.timerRefMessageList);
        this.timerRefMessageList = setInterval(function(){
          parent.refMessageList();
        },1000);
      }
    },
    watch:{
      //key值变化，写入用户列队附属
      NowUserKey: function(){
        if(!this.NowUserInfo['token']){
          return;
        }
        this.UserListKey[this.NowUserInfo['token']] = this.NowUserKey;
      },
    },
    //初始化
    mounted: function () {
      let parent = this;
      //关闭遮罩层
      this.$globConfig.ajaxLoadOn = false;
      //定时刷新用户列表
      parent.refUserList();
      clearInterval(this.timerRefUserList);
      this.timerRefUserList = setInterval(function(){
        parent.refUserList();
      },5000);
      this.autoRefMessageList();
    },
    //离开激活
    destroyed: function(){
      clearInterval(this.timerRefUserList);
      clearInterval(this.timerRefMessageList);
    }
  };
</script>
