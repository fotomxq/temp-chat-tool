/**
 * 引入库
 */
import Vue from 'vue';
import 'whatwg-fetch';
//路由库
import VueRouter from 'vue-router';
//element-ui
import ElementUI from 'element-ui';
import 'element-ui/lib/theme-chalk/index.css';
//全局主体
import './assets/vendor/material-design-icons/css/material-design-iconic-font.css';
import './assets/vendor/beagle/css/app.css';
import './assets/css/glob.css';
//全局配置
import GlobConfig, {GlobConfigData, GlobLogged, GlobMethods} from './base/glob-config';
//引入路由器
import routes from './routers/routers.js';

Vue.use(VueRouter);

Vue.use(ElementUI);

Vue.component('glob-config', GlobConfig);
Vue.prototype.$globConfig = GlobConfigData;
Vue.prototype.$globMethods = GlobMethods;
Vue.prototype.$globLogged = GlobLogged;

//定义路由
const router = new VueRouter({
  //不能使用该模式，否则将无法在electron下正常使用
  //如果希望使用，美化URL，必须利用nginx配合完善
  //mode: 'history',
  routes
});

/**
 * 初始化应用
 */
const app = new Vue({
  router,
}).$mount('#app');

//router.push('/index');
