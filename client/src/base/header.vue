<template>
    <el-header class="header">
        <el-row class="header-container">
        <el-col :span="4">
            <router-link to="/user/center" class="logo"><img src="../assets/imgs/logo.png" alt="LOGO" width="102" height="27" class="logo-img"></router-link>
        </el-col>
        <el-col :span="20" class="header-right">
            <el-row>
            <el-col :span="12" class="header-title">
                <span>{{title}}</span>
            </el-col>
            <el-col :span="12">
                <div class="float-right">
                <router-link to="/login"><span class="icon mdi mdi-power"></span></router-link>
                </div>
            </el-col>
            </el-row>
        </el-col>
        </el-row>
    </el-header>
</template>

<script>
export default {
  props: {
      title: String,
  },
}
</script>
