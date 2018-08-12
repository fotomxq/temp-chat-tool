let PageError = r => require.ensure([], () => r(require('../page/error.vue')), 'page-error');

export default {
  path: '*',
  name: 'error',
  components: {
    'login-before': PageError
  },
};
