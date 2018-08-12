<template>
  <el-main>
    <el-row>
      <el-col :span="24">
        <el-breadcrumb separator="/">
          <el-breadcrumb-item :to="{ path: '/user/center' }">管理中心</el-breadcrumb-item>
          <el-breadcrumb-item>标题</el-breadcrumb-item>
        </el-breadcrumb>
      </el-col>
    </el-row>
    <el-row>
      <el-col :span="24">
        <el-card shadow="always">
          <h3>标题</h3>
          <hr>
        </el-card>
      </el-col>
    </el-row>
  </el-main>
</template>

<script>
  export default {
    components: {},
    data: function () {
      return {}
    },
    props: {},
    methods: {},
    mounted: function () {
    }
  }
</script>
