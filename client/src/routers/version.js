let PageVersion = r => require.ensure([], () => r(require('../page/version.vue')), 'page-version');

export default {
  path: '/version',
  name: 'version',
  components: {
    'login-before': PageVersion
  },
};
