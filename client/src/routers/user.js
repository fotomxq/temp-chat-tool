let PageUser = r => require.ensure([], () => r(require('../page/user.vue')), 'page-user');

export default {
  path: '/user',
  name: 'user',
  components: {
    'login-before': PageUser
  },
  children: [{
      //管理中心
      path: "center",
      name: "user",
      component: r => require.ensure([], () => r(require("../page/center.vue")), "page-center")
    }
  ]
}
