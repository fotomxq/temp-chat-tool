<template>
  <div>
    <el-breadcrumb separator="/">
      <el-breadcrumb-item :to="{ path: '/user/center' }">管理中心</el-breadcrumb-item>
    </el-breadcrumb>
  </div>
</template>

<script>
  export default {
    props: {
      items: String,
    },
    watch: {
      items: function () {
      }
    }
  };
</script>
