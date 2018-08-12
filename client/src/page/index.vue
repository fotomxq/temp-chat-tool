<template>
  <div class="max-size" v-loading="loading">
  </div>
</template>

<script>
  export default {
    data: function () {
      return {
        loading: true,
      }
    },
    created: function () {
      this.$router.push('/login');
    }
  }
</script>
