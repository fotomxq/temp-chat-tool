import Index from './index.js';
import Login from './login.js';
import Error from './error.js';
import Version from './version.js';
import User from './user.js';

export default [
  {
    path: '/',
    name: '',
    redirect: '/login',
  },
  //主页
  Index,
  //登陆页面
  Login,
  //错误页面
  Error,
  //版本页面
  Version,
  //登陆后页面
  User,
];
