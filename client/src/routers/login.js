let PageLogin = r => require.ensure([], () => r(require('../page/login.vue')), 'page-login');

export default {
  path: '/login',
  name: 'login',
  components: {
    'login-before': PageLogin
  },
};
